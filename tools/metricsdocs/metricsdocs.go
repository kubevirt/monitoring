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

package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
)

//go:embed metrics.tmpl
var metricsTemplate string

var (
	metricsTmpl *template.Template
)

func init() {
	var err error
	metricsTmpl, err = template.New("metrics").Funcs(template.FuncMap{
		"escapePipe": func(s string) string {
			return strings.ReplaceAll(s, "|", "\\|")
		},
		"normalizeDescription": func(s string) string {
			s = strings.ReplaceAll(s, "|", "\\|")
			s = strings.ReplaceAll(s, "\n", " ")
			return strings.TrimSpace(s)
		},
		"codeSpan": func(s string) string {
			s = strings.ReplaceAll(s, "`", "\\`")
			s = strings.ReplaceAll(s, "|", "\\|")
			return "`" + s + "`"
		},
	}).Parse(metricsTemplate)
	if err != nil {
		log.Fatalf("ERROR: parsing metrics template: %s", err)
	}
}

func main() {
	r := parseArguments()
	r.checkoutProjects()
	r.createDoc()
}

func parseArguments() *releaseData {
	cacheDir := flag.String("cache-dir", "/tmp/release-tool", "The base directory used to cache git repos in")
	configFile := flag.String("config-file", "", "Name of file containing the components versions")

	flag.Parse()

	if *configFile == "" {
		log.Fatal("--config-file is a required argument")
	}

	org := "kubevirt"
	baseDir := fmt.Sprintf("%s/%s/", *cacheDir, org)

	return &releaseData{
		org:      org,
		projects: createProjects(*configFile, baseDir, org),
	}
}

func createProjects(configFile string, baseDir string, org string) []*project {
	config := getConfig(configFile)

	var projects []*project
	for _, info := range projectsInfo {
		version, ok := config[info.short+"_VERSION"]

		if !ok {
			log.Fatalf("ERROR: config doesn't contain '%s_VERSION' for %s", info.short, info.name)
		}

		projects = append(projects, &project{
			short:   info.short,
			name:    info.name,
			version: version,

			repoDir:        baseDir + info.name,
			repoUrl:        fmt.Sprintf("https://github.com/%s/%s.git", org, info.name),
			metricsDocPath: info.metricsDocPath,
		})
	}

	return projects
}

func getConfig(configFile string) map[string]string {
	config, err := godotenv.Read(configFile)
	if err != nil {
		log.Fatalf("ERROR: reading %s config file", configFile)
	}

	return config
}

func (r *releaseData) checkoutProjects() {
	for _, p := range r.projects {
		err := p.gitCheckoutUpstream()
		if err != nil {
			log.Fatalf("ERROR: checking out upstream: %s\n", err)
		}
		err = p.gitSwitchToBranch(p.version)
		if err != nil {
			log.Fatalf("ERROR: changing to version branch: %s\n", err)
		}
	}
}

func (r *releaseData) createDoc() {
	r.outFile = createFile()
	defer r.outFile.Close()

	metrics := r.parseMetrics()

	r.writeMetrics(metrics)
}

func createFile() *os.File {
	file, err := os.Create("docs/metrics.md")
	if err != nil {
		log.Fatalf("ERROR: creating output file: %s", err)
	}
	return file
}

func (r *releaseData) parseMetrics() []Metric {
	var metrics []Metric

	for _, p := range r.projects {
		content, err := readLines(path.Join(p.repoDir, "/", p.metricsDocPath))
		if err != nil {
			log.Printf("WARNING: %s project does not contain metrics documentation file in '%s'", p.name, p.version)
			continue
		}

		parsed := p.parseMetricsDoc(content, p.name)
		if len(parsed) != 0 {
			metrics = append(metrics, parsed...)
		} else {
			log.Printf("WARNING: %s project metrics documentation file is empty in '%s'", p.name, p.version)
		}
	}

	return metrics
}

func (r *releaseData) writeMetrics(metrics []Metric) {
	sortMetrics(metrics)

	templateData := struct {
		Operators []TemplateOperator
		Metrics   []Metric
	}{
		Operators: r.buildOperators(getOperatorOrder(metrics)),
		Metrics:   metrics,
	}

	if err := metricsTmpl.Execute(r.outFile, templateData); err != nil {
		log.Printf("ERROR: executing metrics template: %s", err)
	}
}

func sortMetrics(metrics []Metric) {
	sort.SliceStable(metrics, func(i, j int) bool {
		if metrics[i].Operator != metrics[j].Operator {
			if metrics[i].Operator == "kubevirt" {
				return true
			}
			if metrics[j].Operator == "kubevirt" {
				return false
			}
			return metrics[i].Operator < metrics[j].Operator
		}
		return metrics[i].Name < metrics[j].Name
	})
}

func getOperatorOrder(metrics []Metric) []string {
	seen := make(map[string]bool)
	var operatorOrder []string
	for _, m := range metrics {
		if !seen[m.Operator] {
			operatorOrder = append(operatorOrder, m.Operator)
			seen[m.Operator] = true
		}
	}
	return operatorOrder
}

func (r *releaseData) buildOperators(operatorOrder []string) []TemplateOperator {
	operatorLinkMap := make(map[string]string)
	seenOperators := make(map[string]bool)
	for _, opName := range operatorOrder {
		seenOperators[opName] = true
	}

	for _, p := range r.projects {
		if seenOperators[p.name] {
			operatorLinkMap[p.name] = p.writeComponentMetrics()
		}
	}

	var operators []TemplateOperator
	for _, opName := range operatorOrder {
		if link, ok := operatorLinkMap[opName]; ok {
			operators = append(operators, TemplateOperator{
				Name: opName,
				Link: link,
			})
		}
	}
	return operators
}

func (p *project) writeComponentMetrics() string {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/kubevirt/%s/releases/tags/%s", p.name, p.version))

	if err == nil {
		defer resp.Body.Close()
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("[%s](https://github.com/kubevirt/%s/tree/%s)", p.name, p.name, p.version)
	}
	return fmt.Sprintf("[%s - %s](https://github.com/kubevirt/%s/releases/tag/%s)", p.name, p.version, p.name, p.version)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (p *project) parseMetricsDoc(content []string, operatorName string) []Metric {
	var metrics []Metric
	var currentMetric *Metric
	inMetricsSection := false

	for i, line := range content {
		if isBeginning(line) {
			metrics, currentMetric = p.handleMetricBeginning(metrics, currentMetric, line, operatorName)
			inMetricsSection = true
		} else if inMetricsSection {
			if isEnd(line) {
				metrics = p.appendMetricIfValid(metrics, currentMetric)
				break
			}

			if currentMetric != nil {
				metrics, currentMetric = p.processMetricLine(metrics, currentMetric, line, content, i)
			}
		}
	}

	return p.appendMetricIfValid(metrics, currentMetric)
}

func (p *project) handleMetricBeginning(metrics []Metric, currentMetric *Metric, line, operatorName string) ([]Metric, *Metric) {
	metrics = p.appendMetricIfValid(metrics, currentMetric)
	metricName := strings.TrimSpace(strings.TrimPrefix(line, "### "))
	return metrics, &Metric{
		Operator:    operatorName,
		Name:        metricName,
		Type:        "",
		Description: "",
	}
}

func (p *project) processMetricLine(metrics []Metric, currentMetric *Metric, line string, content []string, lineIndex int) ([]Metric, *Metric) {
	trimmedLine := strings.TrimSpace(line)
	if trimmedLine == "" {
		if lineIndex+1 < len(content) && isBeginning(content[lineIndex+1]) {
			metrics = p.appendMetricIfValid(metrics, currentMetric)
			return metrics, nil
		}
		return metrics, currentMetric
	}

	if strings.Contains(trimmedLine, "Type:") {
		return p.parseTypeLine(metrics, currentMetric, trimmedLine)
	}

	p.appendToDescription(currentMetric, trimmedLine)
	return metrics, currentMetric
}

func (p *project) parseTypeLine(metrics []Metric, currentMetric *Metric, line string) ([]Metric, *Metric) {
	parts := strings.SplitN(line, "Type:", 2)
	if len(parts) != 2 {
		p.appendToDescription(currentMetric, line)
		return metrics, currentMetric
	}

	desc := strings.TrimSuffix(strings.TrimSpace(parts[0]), ".")
	p.appendToDescription(currentMetric, desc)

	typePart := strings.TrimSuffix(strings.TrimSpace(parts[1]), ".")
	currentMetric.Type = typePart

	metrics = p.appendMetricIfValid(metrics, currentMetric)
	return metrics, nil
}

func (p *project) appendToDescription(m *Metric, text string) {
	if m.Description != "" {
		m.Description += " " + text
	} else {
		m.Description = text
	}
}

func (p *project) appendMetricIfValid(metrics []Metric, m *Metric) []Metric {
	if m != nil && m.Name != "" {
		p.finalizeMetric(m)
		return append(metrics, *m)
	}
	return metrics
}

func (p *project) finalizeMetric(m *Metric) {
	m.Description = strings.TrimSpace(m.Description)
}

func isBeginning(s string) bool {
	return strings.HasPrefix(s, "### ")
}

func isEnd(s string) bool {
	return strings.HasPrefix(s, "## ")
}
