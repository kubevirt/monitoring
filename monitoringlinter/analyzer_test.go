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
