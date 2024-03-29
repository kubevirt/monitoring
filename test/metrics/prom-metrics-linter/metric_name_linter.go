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
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2023 Red Hat, Inc.
 *
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus/testutil/promlint"
	dto "github.com/prometheus/client_model/go"
)

var (
	metricFamilies  string
	operatorName    string
	subOperatorName string
)

func init() {
	// Define command-line flags
	flag.StringVar(&metricFamilies, "metric-families", "", "JSON representation of metric families")
	flag.StringVar(&operatorName, "operator-name", "", "Operator name")
	flag.StringVar(&subOperatorName, "sub-operator-name", "", "Sub-operator name")

	// Parse command-line flags
	flag.Parse()
}

func main() {
	if metricFamilies == "" || operatorName == "" || subOperatorName == "" {
		fmt.Println("metric-families, operator-name, and sub-operator-name must be provided")
		os.Exit(1)
	}

	// Parse metric families from the JSON representation
	families := parseMetricFamilies(metricFamilies)

	// Call customRules function to apply the custom rules to the linter
	linter := promlint.NewWithMetricFamilies(families)

	problems, err := linter.Lint()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run promlint; %v\n", err)
		os.Exit(1)
	}

	for _, family := range families {
		problems = CustomLinterRules(problems, family, operatorName, subOperatorName)
	}

	for _, problem := range problems {
		fmt.Printf("%s: %s\n", problem.Metric, problem.Text)
	}
}

func parseMetricFamilies(jsonStr string) []*dto.MetricFamily {
	var families []*dto.MetricFamily
	err := json.Unmarshal([]byte(jsonStr), &families)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing metric families: %v\n", err)
		os.Exit(1)
	}
	return families
}
