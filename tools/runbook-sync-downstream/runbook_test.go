package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kubevirt/monitoring/tools/runbook-sync-downstream/pkg/transform"
)

var _ = Describe("Runbook", func() {
	Context("Runbook content replacement", Ordered, func() {
		var updateRunbookContent string

		BeforeAll(func() {
			testRunbookContent :=
				"kubectl get <something> -n kubevirt -o json\n" +
					"kubectl get <other_something> --namespace kubevirt -o json\n" +
					"kubectl get <another_other_something> -n kubevirt-hyperconverged -o json\n" +
					"i'm a resource -> namespace: kubevirt-hyperconverged\n"

			updateRunbookContent = transform.ReplaceContents(testRunbookContent)
		})

		It("should replace namespace in '-n kubevirt' format", func() {
			Expect(updateRunbookContent).To(ContainSubstring("oc get <something> -n openshift-cnv -o json"))
		})

		It("should replace namespace in '--namespace kubevirt' format", func() {
			Expect(updateRunbookContent).To(ContainSubstring("oc get <other_something> --namespace openshift-cnv -o json"))
		})

		It("should replace namespace in '-n kubevirt-hyperconverged' format", func() {
			Expect(updateRunbookContent).To(ContainSubstring("oc get <another_other_something> -n openshift-cnv -o json"))
		})

		It("should replace namespace in 'namespace: kubevirt-hyperconverged' format", func() {
			Expect(updateRunbookContent).To(ContainSubstring("i'm a resource -> namespace: openshift-cnv"))
		})
	})

	Context("Generated file comparison", func() {
		var tempDir string

		BeforeEach(func() {
			var err error
			tempDir, err = os.MkdirTemp("", "runbook-diff-*")
			Expect(err).ToNot(HaveOccurred())

			runbooksDir := filepath.Join(tempDir, downstreamRunbooksDir)
			err = os.MkdirAll(runbooksDir, 0755)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDir)
			Expect(err).ToNot(HaveOccurred())
		})

		relPath := filepath.Join(downstreamRunbooksDir, "Foo.md")

		It("reports a match when the generated file is byte-identical", func() {
			content := []byte("# Foo\n\nidentical content\n")
			err := os.WriteFile(filepath.Join(tempDir, relPath), content, 0644)
			Expect(err).ToNot(HaveOccurred())

			matches, err := generatedFileMatches(tempDir, relPath, content)
			Expect(err).ToNot(HaveOccurred())
			Expect(matches).To(BeTrue())
		})

		It("reports a mismatch when the generated file differs", func() {
			err := os.WriteFile(filepath.Join(tempDir, relPath), []byte("# Foo\n\nfixed content\n"), 0644)
			Expect(err).ToNot(HaveOccurred())

			matches, err := generatedFileMatches(tempDir, relPath, []byte("# Foo\n\nbuggy content\n"))
			Expect(err).ToNot(HaveOccurred())
			Expect(matches).To(BeFalse())
		})

		It("returns an error when the generated file is missing", func() {
			_, err := generatedFileMatches(tempDir, relPath, []byte("anything"))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Empty commit handling", func() {
		var tempDir string
		var repo *git.Repository
		var worktree *git.Worktree

		var runbookPath string

		// The repo is seeded with a committed runbook so the index is populated,
		// which is what the real downstream clone looks like. go-git v5.11 only
		// reports git.ErrEmptyCommit for a completely empty index, so a populated
		// index is the only setup that exercises the real clean-worktree path.
		BeforeEach(func() {
			var err error
			tempDir, err = os.MkdirTemp("", "runbook-commit-*")
			Expect(err).ToNot(HaveOccurred())

			repo, err = git.PlainInit(tempDir, false)
			Expect(err).ToNot(HaveOccurred())

			worktree, err = repo.Worktree()
			Expect(err).ToNot(HaveOccurred())

			runbooksDir := filepath.Join(tempDir, downstreamRunbooksDir)
			Expect(os.MkdirAll(runbooksDir, 0755)).To(Succeed())

			runbookPath = filepath.Join(runbooksDir, "Foo.md")
			Expect(os.WriteFile(runbookPath, []byte("# Foo\n\noriginal content\n"), 0644)).To(Succeed())

			_, err = worktree.Add(downstreamRunbooksDir)
			Expect(err).ToNot(HaveOccurred())

			_, err = worktree.Commit("seed runbook", &git.CommitOptions{
				Author: &object.Signature{Name: githubUsername, Email: githubEmail, When: time.Now()},
			})
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDir)
			Expect(err).ToNot(HaveOccurred())
		})

		It("reports nothing to commit when the worktree is clean", func() {
			rbSync := &runbookSync{downstreamRepo: repo}
			err := rbSync.commit(worktree, "no changes")
			Expect(err).To(HaveOccurred())
		})

		It("does not create a commit when the worktree is clean", func() {
			rbSync := &runbookSync{downstreamRepo: repo}
			head, err := repo.Head()
			Expect(err).ToNot(HaveOccurred())

			_ = rbSync.commit(worktree, "no changes")

			newHead, err := repo.Head()
			Expect(err).ToNot(HaveOccurred())
			Expect(newHead.Hash()).To(Equal(head.Hash()))
		})

		It("skips creating a new PR branch when a clean worktree has no PR", func() {
			rbSync := &runbookSync{downstreamRepo: repo}
			err := rbSync.commit(worktree, "no changes")

			action, classifyErr := classifyCommitResult(err, false)
			Expect(classifyErr).ToNot(HaveOccurred())
			Expect(action).To(Equal(commitActionSkipNewPR))
		})

		It("continues to update an existing PR when a clean worktree has an open PR", func() {
			rbSync := &runbookSync{downstreamRepo: repo}
			err := rbSync.commit(worktree, "no changes")

			action, classifyErr := classifyCommitResult(err, true)
			Expect(classifyErr).ToNot(HaveOccurred())
			Expect(action).To(Equal(commitActionKeepExisting))
		})

		It("reports nothing to commit when only files outside the runbooks dir are dirty", func() {
			rbSync := &runbookSync{downstreamRepo: repo}
			Expect(os.WriteFile(filepath.Join(tempDir, "unrelated.txt"), []byte("noise\n"), 0644)).To(Succeed())

			err := rbSync.commit(worktree, "no runbook changes")

			Expect(err).To(MatchError(errNothingToCommit))
		})

		It("commits when the runbook content changed", func() {
			rbSync := &runbookSync{downstreamRepo: repo}
			Expect(os.WriteFile(runbookPath, []byte("# Foo\n\nnew content\n"), 0644)).To(Succeed())

			Expect(rbSync.commit(worktree, "real change")).To(Succeed())
		})

		It("proceeds normally when a commit was created", func() {
			action, classifyErr := classifyCommitResult(nil, false)
			Expect(classifyErr).ToNot(HaveOccurred())
			Expect(action).To(Equal(commitActionProceed))
		})

		It("returns other commit errors for fatal handling", func() {
			sentinel := errors.New("disk full")
			action, classifyErr := classifyCommitResult(sentinel, true)
			Expect(classifyErr).To(MatchError(sentinel))
			Expect(action).To(Equal(commitActionProceed))
		})
	})

	Context("Runbook deprecation", func() {
		var tempDir string
		var testRunbookPath string

		BeforeEach(func() {
			var err error
			tempDir, err = os.MkdirTemp("", "runbook-test-*")
			Expect(err).ToNot(HaveOccurred())

			runbooksDir := filepath.Join(tempDir, downstreamRunbooksDir)
			err = os.MkdirAll(runbooksDir, 0755)
			Expect(err).ToNot(HaveOccurred())

			testRunbookPath = filepath.Join(runbooksDir, "TestRunbook.md")
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDir)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should deprecate a runbook with original content preserved", func() {
			By("creating a test runbook")
			originalContent := `# TestRunbook

## Meaning

This is a test runbook with some content.

## Impact

This describes the impact of the alert.

## Diagnosis

How to diagnose the issue.

## Mitigation

How to fix the issue.`

			err := os.WriteFile(testRunbookPath, []byte(originalContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			By("calling the deprecation function")
			deprecatedRunbook("TestRunbook", tempDir)

			By("reading the updated content")
			updatedContent, err := os.ReadFile(testRunbookPath)
			Expect(err).ToNot(HaveOccurred())

			updatedStr := string(updatedContent)

			By("verifying the updated content")
			Expect(updatedStr).To(ContainSubstring("# TestRunbook [Deprecated]"))

			By("verifying the deprecation notice is added")
			Expect(updatedStr).To(ContainSubstring("This alert is deprecated. You can safely ignore or silence it."))

			By("verifying original content is preserved (without the original title)")
			Expect(updatedStr).To(ContainSubstring("## Meaning"))
			Expect(updatedStr).To(ContainSubstring("This is a test runbook with some content."))
			Expect(updatedStr).To(ContainSubstring("## Impact"))
			Expect(updatedStr).To(ContainSubstring("## Diagnosis"))
			Expect(updatedStr).To(ContainSubstring("## Mitigation"))
			Expect(updatedStr).To(ContainSubstring("How to fix the issue."))
		})

		It("should separate the deprecation notice from the body by a single blank line", func() {
			originalContent := "# TestRunbook\n\n## Meaning\n\nThis is a test runbook.\n"
			Expect(os.WriteFile(testRunbookPath, []byte(originalContent), 0644)).To(Succeed())

			deprecatedRunbook("TestRunbook", tempDir)

			updatedContent, err := os.ReadFile(testRunbookPath)
			Expect(err).ToNot(HaveOccurred())

			Expect(string(updatedContent)).To(Equal(
				"# TestRunbook [Deprecated]\n" +
					"\n" +
					"This alert is deprecated. You can safely ignore or silence it.\n" +
					"\n" +
					"## Meaning\n" +
					"\n" +
					"This is a test runbook.\n"))
		})

		It("should not re-deprecate an already deprecated runbook", func() {
			By("creating a runbook that's already deprecated")
			deprecatedContent := `# TestRunbook [Deprecated]

This alert is deprecated. You can safely ignore or silence it.

## Meaning

This is a test runbook with some content.`

			err := os.WriteFile(testRunbookPath, []byte(deprecatedContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			By("calling the deprecation function")
			deprecatedRunbook("TestRunbook", tempDir)

			By("reading the updated content")
			updatedContent, err := os.ReadFile(testRunbookPath)
			Expect(err).ToNot(HaveOccurred())

			updatedStr := string(updatedContent)

			By("verifying the content remains unchanged")
			Expect(updatedStr).To(Equal(deprecatedContent))

			By("verifying [Deprecated] appears only once")
			deprecatedCount := strings.Count(updatedStr, "[Deprecated]")
			Expect(deprecatedCount).To(Equal(1))
		})
	})

	Context("Filtering runbooks that are already deprecated", func() {
		var tempDir string
		var runbooksDir string

		BeforeEach(func() {
			var err error
			tempDir, err = os.MkdirTemp("", "runbook-filter-*")
			Expect(err).ToNot(HaveOccurred())

			runbooksDir = filepath.Join(tempDir, downstreamRunbooksDir)
			Expect(os.MkdirAll(runbooksDir, 0755)).To(Succeed())
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDir)
			Expect(err).ToNot(HaveOccurred())
		})

		writeRunbook := func(name, content string) {
			Expect(os.WriteFile(filepath.Join(runbooksDir, name), []byte(content), 0644)).To(Succeed())
		}

		It("drops a runbook whose downstream file already carries the deprecation notice", func() {
			writeRunbook("Done.md", "# Done [Deprecated]\n\nThis alert is deprecated.\n")

			kept := filterAlreadyDeprecated([]runbook{{name: "Done.md"}}, tempDir)

			Expect(kept).To(BeEmpty())
		})

		It("keeps a runbook that has not been deprecated yet", func() {
			writeRunbook("Pending.md", "# Pending\n\n## Meaning\n\nStill live.\n")

			kept := filterAlreadyDeprecated([]runbook{{name: "Pending.md"}}, tempDir)

			Expect(kept).To(HaveLen(1))
			Expect(kept[0].name).To(Equal("Pending.md"))
		})

		It("keeps a runbook whose file cannot be read so the error surfaces downstream", func() {
			kept := filterAlreadyDeprecated([]runbook{{name: "Missing.md"}}, tempDir)

			Expect(kept).To(HaveLen(1))
			Expect(kept[0].name).To(Equal("Missing.md"))
		})

		It("keeps only the runbooks still needing deprecation", func() {
			writeRunbook("Done.md", "# Done [Deprecated]\n\nGone.\n")
			writeRunbook("Pending.md", "# Pending\n\nStill live.\n")

			kept := filterAlreadyDeprecated([]runbook{{name: "Done.md"}, {name: "Pending.md"}}, tempDir)

			Expect(kept).To(HaveLen(1))
			Expect(kept[0].name).To(Equal("Pending.md"))
		})
	})

	Context("Listing runbooks that need a sync", func() {
		var upstreamDir, downstreamDir string
		var upstreamRepo, downstreamRepo *git.Repository

		BeforeEach(func() {
			var err error
			upstreamDir, err = os.MkdirTemp("", "runbook-upstream-*")
			Expect(err).ToNot(HaveOccurred())
			downstreamDir, err = os.MkdirTemp("", "runbook-downstream-*")
			Expect(err).ToNot(HaveOccurred())

			upstreamRepo = initRepoWithRunbooks(upstreamDir, upstreamRunbooksDir, map[string]string{
				"Live.md": "# Live\n\nStill documented upstream.\n",
			})
			downstreamRepo = initRepoWithRunbooks(downstreamDir, downstreamRunbooksDir, map[string]string{
				"Live.md":    "# Live\n\nStill documented upstream.\n",
				"Done.md":    "# Done [Deprecated]\n\nAlready deprecated downstream.\n",
				"Pending.md": "# Pending\n\nRemoved upstream, not yet deprecated.\n",
			})
		})

		AfterEach(func() {
			Expect(os.RemoveAll(upstreamDir)).To(Succeed())
			Expect(os.RemoveAll(downstreamDir)).To(Succeed())
		})

		It("excludes runbooks that are already deprecated downstream", func() {
			_, toDeprecate := listRunbooksThatNeedUpdate(downstreamRepo, upstreamRepo)

			var names []string
			for _, rb := range toDeprecate {
				names = append(names, rb.name)
			}
			Expect(names).To(ConsistOf("Pending.md"))
		})
	})
})

func initRepoWithRunbooks(dir, runbooksDir string, runbooks map[string]string) *git.Repository {
	GinkgoHelper()

	repo, err := git.PlainInit(dir, false)
	Expect(err).ToNot(HaveOccurred())

	Expect(os.MkdirAll(filepath.Join(dir, runbooksDir), 0755)).To(Succeed())
	for name, content := range runbooks {
		Expect(os.WriteFile(filepath.Join(dir, runbooksDir, name), []byte(content), 0644)).To(Succeed())
	}

	worktree, err := repo.Worktree()
	Expect(err).ToNot(HaveOccurred())
	_, err = worktree.Add(runbooksDir)
	Expect(err).ToNot(HaveOccurred())
	_, err = worktree.Commit("seed runbooks", &git.CommitOptions{
		Author: &object.Signature{Name: githubUsername, Email: githubEmail, When: time.Now()},
	})
	Expect(err).ToNot(HaveOccurred())

	return repo
}
