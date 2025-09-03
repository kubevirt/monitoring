package main

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

			updateRunbookContent = replaceContents(testRunbookContent)
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
