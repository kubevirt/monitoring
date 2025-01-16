package main

import (
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
})
