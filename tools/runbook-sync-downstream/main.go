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
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/go-git/go-git/v5"
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

func (rbSync *runbookSync) updateRunbook(rb runbook) string {
	lastUpdateDate := rb.upstreamLastUpdated.Format("20060102150405")
	runbookName := strings.Replace(rb.name, ".md", "", -1)
	branchName := fmt.Sprintf("cnv-runbook-sync-%s/%s", lastUpdateDate, runbookName)

	prExists, pr, err := rbSync.prForBranchPreviouslyCreated(branchName)
	if err != nil {
		klog.Fatalf("failed to check if branch exists: %v", err)
	}

	if prExists {
		klog.Infof("PR for '%s' was previously created: %s", branchName, pr.GetHTMLURL())
		return branchName
	}

	worktree, err := newBranchFromMain(rbSync.downstreamRepo, branchName)
	if err != nil {
		klog.Fatalf("failed to create branch: %v", err)
	}

	err = copyRunbook(rb.name)
	if err != nil {
		klog.Fatalf("failed to copy file: %v", err)
	}

	commitMessage := fmt.Sprintf("Sync CNV runbook %s (Updated at %s)", rb.name, rb.upstreamLastUpdated)

	err = rbSync.commitAndPush(worktree, commitMessage)
	if err != nil {
		klog.Fatalf("failed to push changes: %v", err)
	}

	body := fmt.Sprintf(
		"This is an automated PR by 'tools/openshift-virtualization-operator/runbook-sync'.\n\n"+
			"CNV runbook '%s' was updated in upstream https://%s at %s.\n"+
			"This PR syncs the runbook in this repository to contain all new added changes.\n\n"+
			"/cc @machadovilaca",
		rb.name, upstreamRepositoryURL, rb.upstreamLastUpdated,
	)

	err = rbSync.createPR(branchName, commitMessage, body)
	if err != nil {
		klog.Fatalf("failed to create PR: %v", err)
	}

	return branchName
}

func (rbSync *runbookSync) deprecateRunbook(rb runbook) string {
	runbookName := strings.Replace(rb.name, ".md", "", -1)
	branchName := fmt.Sprintf("cnv-runbook-deprecate-%s", runbookName)

	prExists, pr, err := rbSync.prForBranchPreviouslyCreated(branchName)
	if err != nil {
		klog.Fatalf("failed to check if branch exists: %v", err)
	}

	if prExists {
		klog.Infof("PR for '%s' was previously created: %s", branchName, pr.GetHTMLURL())
		return branchName
	}

	worktree, err := newBranchFromMain(rbSync.downstreamRepo, branchName)
	if err != nil {
		klog.Fatalf("failed to create branch: %v", err)
	}

	klog.Infof("updating runbook with deprecation message")
	deprecatedRunbook(runbookName)

	commitMessage := fmt.Sprintf("Deprecate CNV runbook %s", runbookName)

	err = rbSync.commitAndPush(worktree, commitMessage)
	if err != nil {
		klog.Fatalf("failed to push changes: %v", err)
	}

	body := fmt.Sprintf(
		"This is an automated PR by 'tools/openshift-virtualization-operator/runbook-sync'.\n\n"+
			"CNV runbook '%s' was deprecated in upstream https://%s.\n"+
			"This PR moves the runbook to the 'deprecate' subdirectory.\n\n"+
			"/cc @machadovilaca",
		rb.name, upstreamRepositoryURL,
	)

	err = rbSync.createPR(branchName, commitMessage, body)
	if err != nil {
		klog.Fatalf("failed to create PR: %v", err)
	}

	return branchName
}

func (rbSync *runbookSync) commitAndPush(worktree *git.Worktree, msg string) error {
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

	if rbSync.dryRun {
		klog.Warning("[DRY RUN] skipping push")
		return nil
	}

	err = rbSync.downstreamRepo.Push(&git.PushOptions{
		RemoteName: forkRemoteName,
	})
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

	return true, prs[0], nil
}

func (rbSync *runbookSync) createPR(branchName string, title string, body string) error {
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
		return nil
	}

	pr, _, err := rbSync.ghClient.PullRequests.Create(context.Background(), downstreamRepositoryOwner, downstreamRepositoryName, prOpts)
	if err != nil {
		return err
	}

	klog.Infof("PR created: %s", pr.GetHTMLURL())

	return nil
}

func deprecatedRunbook(runbookName string) {
	tmpl, err := template.ParseFS(deprecatedRunbookTemplate, "templates/deprecated_runbook.tmpl")
	if err != nil {
		klog.Fatalf("failed to parse template: %v", err)
	}

	p := path.Join(downstreamCloneDir, downstreamRunbooksDir, runbookName+".md")
	f, err := os.OpenFile(p, os.O_RDWR, 0644)
	if err != nil {
		klog.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, struct{ RunbookName string }{RunbookName: runbookName})
	if err != nil {
		klog.Fatalf("failed to execute template: %v", err)
	}
}
