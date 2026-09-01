package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
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

		BeforeEach(func() {
			var err error
			tempDir, err = os.MkdirTemp("", "runbook-commit-*")
			Expect(err).ToNot(HaveOccurred())

			repo, err = git.PlainInit(tempDir, false)
			Expect(err).ToNot(HaveOccurred())

			worktree, err = repo.Worktree()
			Expect(err).ToNot(HaveOccurred())

			// An empty runbooks dir means commit() stages nothing, which is the
			// clean-worktree case that yields git.ErrEmptyCommit.
			runbooksDir := filepath.Join(tempDir, downstreamRunbooksDir)
			Expect(os.MkdirAll(runbooksDir, 0755)).To(Succeed())
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDir)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns git.ErrEmptyCommit when the worktree is clean", func() {
			rbSync := &runbookSync{downstreamRepo: repo}
			err := rbSync.commit(worktree, "no changes")
			Expect(errors.Is(err, git.ErrEmptyCommit)).To(BeTrue())
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
})
