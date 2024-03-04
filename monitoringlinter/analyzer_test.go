package monitoringlinter_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/kubevirt/monitoring/monitoringlinter"
)

func TestAllUseCases(t *testing.T) {
	for _, testcase := range []struct {
		name string
		data string
	}{
		{
			name: "Verify metrics registrations in pkg/controller, using operatorobservability, is reported.",
			data: "a/testrepo/pkg/controller/operatorobservability",
		},
		{
			name: "Verify metrics registrations in pkg/controller, using prometheus, is reported.",
			data: "a/testrepo/pkg/controller/prometheus",
		},
		{
			name: "Verify metrics registrations in pkg/monitoring, using operatorobservability, is not reported.",
			data: "a/testrepo/pkg/monitoring/operatorobservability",
		},
		{
			name: "Verify metrics registrations in pkg/monitoring, using prometheus, is reported.",
			data: "a/testrepo/pkg/monitoring/prometheus",
		},
	} {
		t.Run(testcase.name, func(tt *testing.T) {
			analysistest.Run(tt, analysistest.TestData(), monitoringlinter.NewAnalyzer(), testcase.data)
		})
	}
}
