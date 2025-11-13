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
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus/testutil/promlint"
	dto "github.com/prometheus/client_model/go"
)

func CustomMetricsValidation(problems []promlint.Problem, mf *dto.MetricFamily, operatorName, subOperatorName string) []promlint.Problem {
	// Validate metric prefix
	nameParts := strings.Split(*mf.Name, "_")

	if nameParts[0] != operatorName {
		problems = append(problems, promlint.Problem{
			Metric: *mf.Name,
			Text:   fmt.Sprintf(`name need to start with %s_`, operatorName),
		})
	} else if operatorName != subOperatorName && nameParts[1] != subOperatorName {
		problems = append(problems, promlint.Problem{
			Metric: *mf.Name,
			Text:   fmt.Sprintf(`name need to start with "%s_%s_"`, operatorName, subOperatorName),
		})
	}

	// If promlint fails on a "total" suffix, check also for "_timestamp_seconds" suffix. If it exists, do not fail
	var newProblems []promlint.Problem
	for _, problem := range problems {
		if strings.Contains(problem.Text, "counter metrics should have \"_total\" suffix") {
			if !strings.HasSuffix(problem.Metric, "_timestamp_seconds") {
				problem.Text = "counter metrics should have \"_total\" or \"_timestamp_seconds\" suffix"
				newProblems = append(newProblems, problem)
			}
		} else {
			newProblems = append(newProblems, problem)
		}
	}

	return newProblems
}

func CustomRecordingRuleValidation(rr recordingRule) []promlint.Problem {
	var problems []promlint.Problem

	name := rr.Record
	if name == "" {
		return problems
	}

	// Require name structure: level:metric:operations
	parts, nameOkProblems, ok := validateRecordingRuleNameStructure(name)
	if !ok {
		return append(problems, nameOkProblems...)
	}

	// Enforce metric segment prefix: it should start with operator_suboperator like metrics do
	problems = append(problems, validateRecordingRuleMetricPrefix(name, parts, operatorName, subOperatorName)...)

	// Single/Multiple operations: detect ops in expr and enforce suffix rules
	problems = append(problems, validateRecordingRuleOpsSuffix(rr, parts)...)

	// Forbid duplicated operations
	problems = append(problems, validateRecordingRuleNoDuplicateOps(name, parts)...)

	return problems
}

func validateRecordingRuleNameStructure(name string) ([]string, []promlint.Problem, bool) {
	parts := strings.Split(name, ":")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, []promlint.Problem{
			{Metric: name, Text: "recording rule name must be level:metric:operations. if there is no obvious operations, use 'sum'"},
		}, false
	}
	return parts, nil, true
}

func validateRecordingRuleMetricPrefix(name string, parts []string, operatorName, subOperatorName string) []promlint.Problem {
	var problems []promlint.Problem
	metricPart := parts[1]
	expectedPrefix := operatorName + "_"
	if operatorName != subOperatorName {
		expectedPrefix = operatorName + "_" + subOperatorName + "_"
	}
	if !strings.HasPrefix(metricPart, expectedPrefix) {
		problems = append(problems, promlint.Problem{
			Metric: name,
			Text:   fmt.Sprintf("metric (second segment in level:metric:operations) need to start with %q", expectedPrefix),
		})
	}
	return problems
}

func validateRecordingRuleOpsSuffix(rr recordingRule, parts []string) []promlint.Problem {
	name := rr.Record
	var problems []promlint.Problem

	aggOps := detectOps(rr.Expr, promAggRegex)
	timeOps := detectOps(rr.Expr, promTimeRegex)
	total := len(aggOps) + len(timeOps)
	opsPart := parts[2]

	if total == 1 {
		if len(aggOps) == 1 {
			expected := aggOps[0]
			if !strings.HasSuffix(opsPart, expected) {
				problems = append(problems, promlint.Problem{
					Metric: name,
					Text:   fmt.Sprintf("single aggregation detected in expr: operations (third segment in level:metric:operations) should end with %q", expected),
				})
			}
		} else if len(timeOps) == 1 {
			expected := timeOps[0]
			timeSuffixRe, err := regexp.Compile(`(?i)` + expected + `[0-9]+[smhdwy]+$`)
			if err != nil {
				fmt.Fprintf(os.Stderr, "time op suffix regex compile failed for %q: %v\n", expected, err)
			} else if !timeSuffixRe.MatchString(opsPart) {
				problems = append(problems, promlint.Problem{
					Metric: name,
					Text:   fmt.Sprintf("single time operation detected in expr: operations (third segment in level:metric:operations) should end with %q", expected+"<duration>"),
				})
			}
		}
	} else if total > 1 {
		present := false
		for _, op := range append(aggOps, timeOps...) {
			if strings.Contains(opsPart, op) {
				present = true
				break
			}
		}
		if !present {
			problems = append(problems, promlint.Problem{
				Metric: name,
				Text:   "multiple operations detected in expr: at least one operation should appear in operations (third segment in level:metric:operations)",
			})
		}
	}
	return problems
}

func validateRecordingRuleNoDuplicateOps(name string, parts []string) []promlint.Problem {
	var problems []promlint.Problem
	seen := map[string]struct{}{}
	tokens := strings.Split(parts[2], "_")
	for _, tok := range tokens {
		if tok == "" {
			continue
		}
		if _, ok := seen[tok]; ok {
			problems = append(problems, promlint.Problem{
				Metric: name,
				Text:   "operations (third segment in level:metric:operations) contains duplicate tokens; merge duplicates (e.g., min_min -> min)",
			})
			break
		}
		seen[tok] = struct{}{}
	}
	return problems
}

var (
	promAggOps = []string{
		// Aggregations
		"sum", "avg", "min", "max", "count", "quantile", "stddev", "stdvar",
		"group", "count_values", "limitk", "limit_ratio",
		// Selectors considered as aggregations
		"topk", "bottomk",
	}
	promTimeOps = []string{
		// Time/transform functions
		"rate", "irate", "increase", "delta", "idelta", "deriv",
		"avg_over_time", "min_over_time", "max_over_time", "sum_over_time",
		"count_over_time", "quantile_over_time", "stddev_over_time", "stdvar_over_time",
	}
	// Precompiled regexes for op detection
	promAggRegex  = regexp.MustCompile(`\b(` + strings.Join(promAggOps, "|") + `)\b`)
	promTimeRegex = regexp.MustCompile(`\b(` + strings.Join(promTimeOps, "|") + `)\b`)
)

func detectOps(expr string, re *regexp.Regexp) []string {
	uniq := map[string]struct{}{}
	for _, m := range re.FindAllStringSubmatch(expr, -1) {
		if len(m) > 1 {
			uniq[m[1]] = struct{}{}
		}
	}
	var opsList []string
	for k := range uniq {
		opsList = append(opsList, k)
	}
	return opsList
}
