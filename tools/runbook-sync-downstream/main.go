/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2024 Red Hat, Inc.
 *
 */

package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v60/github"
	"k8s.io/klog/v2"
)

const (
	githubUsername = "hco-bot"
	githubEmail    = "71450783+hco-bot@users.noreply.github.com"

	upstreamCloneDir      = "/tmp/kubevirt-monitoring"
	upstreamRepositoryURL = "github.com/kubevirt/monitoring"
	upstreamRunbooksDir   = "docs/runbooks"

	downstreamMainBranch      = "master"
	downstreamCloneDir        = "/tmp/runbooks"
	downstreamRepositoryOwner = "openshift"
	downstreamRepositoryFork  = "hco-bot"
	downstreamRepositoryName  = "runbooks"
	downstreamRunbooksDir     = "alerts/openshift-virtualization-operator"

	originRemoteName = "origin"
	forkRemoteName   = "fork"
)

var (
	downstreamRepositoryURL = fmt.Sprintf("github.com/%s/%s", downstreamRepositoryOwner, downstreamRepositoryName)
	forkedRepositoryURL     = fmt.Sprintf("github.com/%s/%s", downstreamRepositoryFork, downstreamRepositoryName)

	prReviewersUsernames = []string{"machadovilaca", "sradco", "avlitman", "jherrman"}
	prReviewersFmt       = fmt.Sprintf("/cc @%s", strings.Join(prReviewersUsernames, " @"))

	//go:embed templates/deprecated_runbook.tmpl
	deprecatedRunbookTemplate embed.FS
)

type runbookSyncArgs struct {
	githubToken string
	dryRun      bool
}

type runbookSync struct {
	ghClient       *github.Client
	downstreamRepo *git.Repository
	dryRun         bool
}

func main() {
	rbSyncArgs := getRunbookSyncArgs()

	downstreamRepo, upstreamRepo := setup(rbSyncArgs.githubToken)
	runbooksToUpdate, runbooksToDeprecate := listRunbooksThatNeedUpdate(downstreamRepo, upstreamRepo)

	for _, r := range runbooksToUpdate {
		klog.Infof("runbook %s will be updated. Last update: %s, upstream last update: %s", r.name, r.lastLocalUpdate, r.upstreamLastUpdated)
	}

	for _, r := range runbooksToDeprecate {
		klog.Infof("runbook %s will be deprecated. Last update: %s", r.name, r.lastLocalUpdate)
	}

	rbSync := &runbookSync{
		ghClient:       github.NewClient(nil).WithAuthToken(rbSyncArgs.githubToken),
		downstreamRepo: downstreamRepo,
		dryRun:         rbSyncArgs.dryRun,
	}

	rbSync.createRunbooksBranches(runbooksToUpdate, runbooksToDeprecate)
}

func getRunbookSyncArgs() runbookSyncArgs {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		klog.Fatal("GITHUB_TOKEN environment variable is required")
	}

	dryRun := os.Getenv("DRY_RUN")
	if dryRun == "" {
		dryRun = "true"
	}
	if dryRun != "true" && dryRun != "false" {
		klog.Fatal("DRY_RUN environment variable must be 'true' or 'false'")
	}
	klog.Infof("dry run: %s", dryRun)

	return runbookSyncArgs{
		githubToken: githubToken,
		dryRun:      dryRun != "false",
	}
}

func (rbSync *runbookSync) createRunbooksBranches(runbooksToUpdate []runbook, runbooksToDeprecate []runbook) {
	if len(runbooksToUpdate) == 0 {
		klog.Info("no runbooks to update")
	}

	for _, rb := range runbooksToUpdate {
		klog.Infof("---")
		_ = rbSync.updateRunbook(rb)
	}

	if len(runbooksToDeprecate) == 0 {
		klog.Info("no runbooks to deprecate")
	}

	for _, rb := range runbooksToDeprecate {
		klog.Infof("---")
		_ = rbSync.deprecateRunbook(rb)
	}
}

// runbookPRSync captures everything the shared PR sync flow needs to decide
// whether to skip, create, or update a downstream runbook PR.
type runbookPRSync struct {
	branchName string
	filePath   string // path of the runbook file relative to the repo root
	commitMsg  string
	prTitle    string
	prBody     string
	generate   func() error                              // writes the runbook file to disk
	afterSync  func(currentPR *github.PullRequest) error // optional, runs after a PR is created or updated in place, receiving the PR to keep open
}

// syncRunbookPR skips closed PRs, updates open PRs in place when the generated
// content differs, and otherwise creates a fresh PR.
func (rbSync *runbookSync) syncRunbookPR(p runbookPRSync) string {
	prExists, pr, err := rbSync.prForBranchPreviouslyCreated(p.branchName)
	if err != nil {
		klog.Fatalf("failed to check if branch exists: %v", err)
	}

	if prExists && pr.GetState() == "closed" {
		klog.Infof("PR for '%s' is closed, skipping: %s", p.branchName, pr.GetHTMLURL())
		return p.branchName
	}

	worktree, err := newBranchFromMain(rbSync.downstreamRepo, p.branchName)
	if err != nil {
		klog.Fatalf("failed to create branch: %v", err)
	}

	if err := p.generate(); err != nil {
		klog.Fatalf("failed to generate runbook: %v", err)
	}

	action, err := classifyCommitResult(rbSync.commit(worktree, p.commitMsg), prExists)
	if err != nil {
		klog.Fatalf("failed to commit changes: %v", err)
	}

	if action == commitActionSkipNewPR {
		klog.Infof("no changes to sync for '%s', skipping new PR", p.branchName)
		return p.branchName
	}

	if prExists {
		rbSync.updateOpenPR(p, pr)
		if err := rbSync.runAfterSync(p, pr); err != nil {
			klog.Fatalf("failed to run post-sync step: %v", err)
		}
		return p.branchName
	}

	if err := rbSync.push(p.branchName, false); err != nil {
		klog.Fatalf("failed to push changes: %v", err)
	}

	newPR, err := rbSync.createPR(p.branchName, p.prTitle, p.prBody)
	if err != nil {
		klog.Fatalf("failed to create PR: %v", err)
	}

	if err := rbSync.runAfterSync(p, newPR); err != nil {
		klog.Fatalf("failed to run post-sync step: %v", err)
	}

	return p.branchName
}

// runAfterSync runs the optional post-sync hook, closing other outdated PRs for
// the same runbook on both the create and in-place update paths.
func (rbSync *runbookSync) runAfterSync(p runbookPRSync, currentPR *github.PullRequest) error {
	if p.afterSync == nil {
		return nil
	}
	return p.afterSync(currentPR)
}

// updateOpenPR force-pushes the regenerated content to an already-open PR's
// branch, but only when it differs from what is already there.
func (rbSync *runbookSync) updateOpenPR(p runbookPRSync, pr *github.PullRequest) {
	remoteContent, err := forkBranchFileContent(rbSync.downstreamRepo, p.branchName, p.filePath)
	if err != nil {
		klog.Fatalf("failed to read fork branch content: %v", err)
	}

	if remoteContent != nil {
		matches, err := generatedFileMatches(downstreamCloneDir, p.filePath, remoteContent)
		if err != nil {
			klog.Fatalf("failed to compare runbook content: %v", err)
		}
		if matches {
			klog.Infof("open PR '%s' is already up to date, skipping: %s", p.branchName, pr.GetHTMLURL())
			return
		}
	}

	klog.Infof("open PR '%s' is out of date, updating: %s", p.branchName, pr.GetHTMLURL())
	if err := rbSync.push(p.branchName, true); err != nil {
		klog.Fatalf("failed to force-push changes: %v", err)
	}
}

func (rbSync *runbookSync) updateRunbook(rb runbook) string {
	lastUpdateDate := rb.upstreamLastUpdated.Format("20060102150405")
	runbookName := strings.Replace(rb.name, ".md", "", -1)
	branchName := fmt.Sprintf("cnv-runbook-sync-%s/%s", lastUpdateDate, runbookName)

	commitMessage := fmt.Sprintf("Sync CNV runbook %s (Updated at %s)", rb.name, rb.upstreamLastUpdated)

	body := fmt.Sprintf(
		"This is an automated PR by 'tools/openshift-virtualization-operator/runbook-sync'.\n\n"+
			"CNV runbook '%s' was updated in upstream https://%s at %s.\n"+
			"This PR syncs the runbook in this repository to contain all new added changes.\n\n"+
			"%s",
		rb.name, upstreamRepositoryURL, rb.upstreamLastUpdated, prReviewersFmt,
	)

	return rbSync.syncRunbookPR(runbookPRSync{
		branchName: branchName,
		filePath:   path.Join(downstreamRunbooksDir, rb.name),
		commitMsg:  commitMessage,
		prTitle:    commitMessage,
		prBody:     body,
		generate:   func() error { return copyRunbook(rb.name) },
		afterSync: func(currentPR *github.PullRequest) error {
			return rbSync.closeOutdatedRunbookPRs(currentPR, runbookName)
		},
	})
}

func (rbSync *runbookSync) deprecateRunbook(rb runbook) string {
	runbookName := strings.Replace(rb.name, ".md", "", -1)
	branchName := fmt.Sprintf("cnv-runbook-deprecate-%s", runbookName)

	commitMessage := fmt.Sprintf("Deprecate CNV runbook %s", runbookName)

	body := fmt.Sprintf(
		"This is an automated PR by 'tools/openshift-virtualization-operator/runbook-sync'.\n\n"+
			"CNV runbook '%s' was deprecated in upstream https://%s.\n"+
			"This PR moves the runbook to the 'deprecate' subdirectory.\n\n"+
			"%s",
		rb.name, upstreamRepositoryURL, prReviewersFmt,
	)

	return rbSync.syncRunbookPR(runbookPRSync{
		branchName: branchName,
		filePath:   path.Join(downstreamRunbooksDir, rb.name),
		commitMsg:  commitMessage,
		prTitle:    commitMessage,
		prBody:     body,
		generate: func() error {
			klog.Infof("updating runbook with deprecation message")
			deprecatedRunbook(runbookName, downstreamCloneDir)
			return nil
		},
	})
}

// commitAction describes how syncRunbookPR should proceed after committing.
type commitAction int

const (
	commitActionProceed      commitAction = iota // a commit was created; continue the sync
	commitActionSkipNewPR                        // nothing to commit and no PR exists; skip creating an empty PR
	commitActionKeepExisting                     // nothing to commit but a PR exists; continue to update it in place
)

// classifyCommitResult decides how to proceed based on the commit result. A
// clean worktree (git.ErrEmptyCommit) means the generated content already
// matches the base branch: skip creating a brand-new empty PR, but still
// continue so an existing PR is updated in place. Any other commit error is
// returned to the caller for fatal handling.
func classifyCommitResult(err error, prExists bool) (commitAction, error) {
	if err == nil {
		return commitActionProceed, nil
	}

	if errors.Is(err, git.ErrEmptyCommit) {
		if prExists {
			return commitActionKeepExisting, nil
		}
		return commitActionSkipNewPR, nil
	}

	return commitActionProceed, err
}

func (rbSync *runbookSync) commit(worktree *git.Worktree, msg string) error {
	_, err := worktree.Add(downstreamRunbooksDir)
	if err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	_, err = worktree.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  githubUsername,
			Email: githubEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}
	klog.Infof("successfully committed: %s", msg)

	return nil
}

// push publishes branchName to the fork remote. When force is true it overwrites
// the remote branch, which is required to update an already-open PR whose branch
// has diverged from the freshly recreated local branch.
func (rbSync *runbookSync) push(branchName string, force bool) error {
	if rbSync.dryRun {
		if force {
			klog.Warningf("[DRY RUN] would force-push to update branch %s", branchName)
		} else {
			klog.Warning("[DRY RUN] skipping push")
		}
		return nil
	}

	pushOpts := &git.PushOptions{
		RemoteName: forkRemoteName,
	}
	if force {
		refSpec := config.RefSpec(fmt.Sprintf("+refs/heads/%s:refs/heads/%s", branchName, branchName))
		pushOpts.RefSpecs = []config.RefSpec{refSpec}
		pushOpts.Force = true
	}

	err := rbSync.downstreamRepo.Push(pushOpts)
	if err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}
	klog.Info("successfully pushed changes")

	return nil
}

func (rbSync *runbookSync) prForBranchPreviouslyCreated(branchName string) (bool, *github.PullRequest, error) {
	prs, _, err := rbSync.ghClient.PullRequests.List(context.Background(), downstreamRepositoryOwner, downstreamRepositoryName, &github.PullRequestListOptions{
		State: "all",
		Head:  fmt.Sprintf("%s:%s", downstreamRepositoryFork, branchName),
	})
	if err != nil {
		return false, nil, err
	}

	if len(prs) == 0 {
		return false, nil, nil
	}

	// Prefer an open PR when both open and closed PRs exist for the branch, so an
	// open PR is updated in place rather than being treated as closed.
	for _, pr := range prs {
		if pr.GetState() == "open" {
			return true, pr, nil
		}
	}

	return true, prs[0], nil
}

func (rbSync *runbookSync) createPR(branchName string, title string, body string) (*github.PullRequest, error) {
	headBranch := fmt.Sprintf("%s:%s", downstreamRepositoryFork, branchName)
	baseBranch := downstreamMainBranch

	prOpts := &github.NewPullRequest{
		Title: &title,
		Head:  &headBranch,
		Base:  &baseBranch,
		Body:  &body,
	}

	if rbSync.dryRun {
		klog.Warningf("[DRY RUN] skipping PR creation '%s', %s => %s/%s %s", *prOpts.Title, *prOpts.Head, downstreamRepositoryOwner, downstreamRepositoryName, *prOpts.Base)
		return nil, nil
	}

	pr, _, err := rbSync.ghClient.PullRequests.Create(context.Background(), downstreamRepositoryOwner, downstreamRepositoryName, prOpts)
	if err != nil {
		return nil, err
	}

	klog.Infof("PR created: %s", pr.GetHTMLURL())

	return pr, nil
}

func (rbSync *runbookSync) closeOutdatedRunbookPRs(keepPR *github.PullRequest, runbookName string) error {
	prs, _, err := rbSync.ghClient.PullRequests.List(context.Background(), downstreamRepositoryOwner, downstreamRepositoryName, &github.PullRequestListOptions{
		State: "open",
		Base:  downstreamMainBranch,
	})
	if err != nil {
		return err
	}

	for _, oldPR := range prs {
		if isAutomatedPRForSameRunbook(oldPR, runbookName) && oldPR.GetNumber() != keepPR.GetNumber() {
			if err := rbSync.closeOutdatedRunbookPR(oldPR, keepPR); err != nil {
				return err
			}
		}
	}

	return nil
}

func isAutomatedPRForSameRunbook(oldPR *github.PullRequest, runbookName string) bool {
	return strings.Contains(oldPR.GetTitle(), runbookName) && oldPR.GetUser().GetLogin() == githubUsername
}

func (rbSync *runbookSync) closeOutdatedRunbookPR(oldPR *github.PullRequest, keepPR *github.PullRequest) error {
	klog.Infof("closing outdated PR: %s", oldPR.GetHTMLURL())

	if rbSync.dryRun {
		klog.Warning("[DRY RUN] skipping PR close")
		return nil
	}

	body := *oldPR.Body + fmt.Sprintf("\n\nThis pull request has been closed in favor of another one. Please refer to the updated PR for the latest changes and discussion: %s.", keepPR.GetHTMLURL())
	_, _, err := rbSync.ghClient.PullRequests.Edit(context.Background(), downstreamRepositoryOwner, downstreamRepositoryName, oldPR.GetNumber(), &github.PullRequest{
		State: github.String("closed"),
		Body:  github.String(body),
	})
	if err != nil {
		return err
	}

	return nil
}

func deprecatedRunbook(runbookName string, cloneDir string) {
	p := path.Join(cloneDir, downstreamRunbooksDir, runbookName+".md")

	content, err := os.ReadFile(p)
	if err != nil {
		klog.Fatalf("failed to read file: %v", err)
	}

	originalContent := string(content)

	if strings.Contains(originalContent, "[Deprecated]") {
		klog.Infof("runbook %s is already deprecated", runbookName)
		return
	}

	lines := strings.Split(originalContent, "\n")
	var contentWithoutTitle []string
	titleRemoved := false

	for _, line := range lines {
		if !titleRemoved && strings.HasPrefix(line, "# ") {
			titleRemoved = true
			continue
		}
		contentWithoutTitle = append(contentWithoutTitle, line)
	}

	tmpl, err := template.ParseFS(deprecatedRunbookTemplate, "templates/deprecated_runbook.tmpl")
	if err != nil {
		klog.Fatalf("failed to parse template: %v", err)
	}

	f, err := os.Create(p)
	if err != nil {
		klog.Fatalf("failed to create file: %v", err)
	}
	defer f.Close()

	templateData := struct {
		RunbookName     string
		OriginalContent string
	}{
		RunbookName:     runbookName,
		OriginalContent: strings.Join(contentWithoutTitle, "\n"),
	}

	err = tmpl.Execute(f, templateData)
	if err != nil {
		klog.Fatalf("failed to execute template: %v", err)
	}
}
