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

	return &releaseData{
		projects: createProjects(*configFile, *cacheDir),
	}
}

func createProjects(configFile string, cacheDir string) []*project {
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
			org:     info.org,
			version: version,

			repoDir:        fmt.Sprintf("%s/%s/%s", cacheDir, info.org, info.name),
			repoUrl:        fmt.Sprintf("https://github.com/%s/%s.git", info.org, info.name),
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
		mi, mj := metrics[i], metrics[j]
		if mi.Operator != mj.Operator {
			return mi.Operator == "kubevirt" || (mj.Operator != "kubevirt" && mi.Operator < mj.Operator)
		}
		if mi.Kind != mj.Kind {
			return mi.Kind == "Metric"
		}
		return mi.Name < mj.Name
	})
}

func getOperatorOrder(metrics []Metric) []string {
	seen := make(map[string]bool)
	var kubevirt, others []string
	for _, m := range metrics {
		if !seen[m.Operator] {
			if m.Operator == "kubevirt" {
				kubevirt = append(kubevirt, m.Operator)
			} else {
				others = append(others, m.Operator)
			}
			seen[m.Operator] = true
		}
	}
	sort.Strings(others)
	return append(kubevirt, others...)
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
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", p.org, p.name, p.version))

	if err == nil {
		defer resp.Body.Close()
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("[%s](https://github.com/%s/%s/tree/%s)", p.name, p.org, p.name, p.version)
	}
	return fmt.Sprintf("[%s - %s](https://github.com/%s/%s/releases/tag/%s)", p.name, p.version, p.org, p.name, p.version)
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

	for _, line := range content {
		if isEnd(line) {
			break
		}

		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "|------") {
			continue
		}

		if strings.HasPrefix(trimmed, "|") {
			if metric := p.parseTableRow(line, operatorName); metric != nil {
				metrics = append(metrics, *metric)
			}
		}
	}

	return metrics
}

func (p *project) parseTableRow(line string, operatorName string) *Metric {
	parts := strings.Split(line, "|")
	if len(parts) < 5 {
		return nil
	}

	name := strings.TrimSpace(parts[1])
	kind := strings.TrimSpace(parts[2])
	metricType := strings.TrimSpace(parts[3])
	description := strings.TrimSpace(parts[4])

	if name == "" || name == "Name" {
		return nil
	}

	return &Metric{
		Operator:    operatorName,
		Name:        name,
		Kind:        kind,
		Type:        metricType,
		Description: description,
	}
}

func isEnd(s string) bool {
	return strings.HasPrefix(s, "## ")
}
