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
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/prometheus/client_golang/prometheus/testutil/promlint"
	dto "github.com/prometheus/client_model/go"
)

var (
	metricFamilies  string
	operatorName    string
	subOperatorName string
)

//go:embed allowlist.json
var embeddedAllowlist []byte

type allowlistConfig struct {
	Operators map[string][]string `json:"operators"`
}

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

	// Parse input JSON containing metricFamilies and optional recordingRules.
	families, recordingRules, err := parseInput(metricFamilies)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing input: %v\n", err)
		os.Exit(1)
	}

	// Lint metrics with promlint, then apply custom validations
	linter := promlint.NewWithMetricFamilies(families)

	problems, err := linter.Lint()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run promlint; %v\n", err)
		os.Exit(1)
	}

	// Collect and print metric problems first (sorted by metric name)
	metricProblems := problems
	for _, family := range families {
		metricProblems = CustomMetricsValidation(metricProblems, family, operatorName, subOperatorName)
	}
	sort.Slice(metricProblems, func(i, j int) bool { return metricProblems[i].Metric < metricProblems[j].Metric })
	for _, p := range metricProblems {
		fmt.Printf("%s: %s\n", p.Metric, p.Text)
	}

	// Then collect and print recording rule problems (sorted by rule name), skipping allowlisted names
	allow := map[string]struct{}{}
	if len(embeddedAllowlist) > 0 {
		var cfg allowlistConfig
		if err := json.Unmarshal(embeddedAllowlist, &cfg); err != nil {
			panic("failed to read the allow list; " + err.Error())
		}
		for _, n := range cfg.Operators[subOperatorName] {
			allow[n] = struct{}{}
		}
	}
	ruleProblems := []promlint.Problem{}
	for _, rr := range recordingRules {
		if _, ok := allow[rr.Record]; !ok {
			ruleProblems = append(ruleProblems, CustomRecordingRuleValidation(rr)...)
		}
	}
	sort.Slice(ruleProblems, func(i, j int) bool { return ruleProblems[i].Metric < ruleProblems[j].Metric })
	for _, p := range ruleProblems {
		fmt.Printf("%s: %s\n", p.Metric, p.Text)
	}
}

type recordingRule struct {
	Record string `json:"record"`
	Expr   string `json:"expr"`
	Type   string `json:"type,omitempty"`
}

type inputJSON struct {
	MetricFamilies []*dto.MetricFamily `json:"metricFamilies"`
	RecordingRules []recordingRule     `json:"recordingRules"`
}

func parseInput(jsonStr string) ([]*dto.MetricFamily, []recordingRule, error) {
	// Expect a JSON with metricFamilies and optional recordingRules
	var env inputJSON
	if err := json.Unmarshal([]byte(jsonStr), &env); err != nil {
		return nil, nil, err
	}
	return env.MetricFamilies, env.RecordingRules, nil
}
