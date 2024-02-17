package monitoringlocationlinter_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/monitoring/test/monitoringlocationlinter"
)

func TestAllUseCases(t *testing.T) {
	for _, testcase := range []struct {
		name string
		data string
	}{
		{
			name: "Exclude files in pkg/monitoring.",
			data: "a/testrepo/pkg/monitoring",
		},
		{
			name: "Report metrics registrations in pkg/controller.",
			data: "a/testrepo/pkg/controller",
		},
	} {
		t.Run(testcase.name, func(tt *testing.T) {
			analysistest.Run(tt, analysistest.TestData(), monitoringlocationlinter.NewAnalyzer(), testcase.data)
		})
	}
}
