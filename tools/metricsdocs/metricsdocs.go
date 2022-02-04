package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
)

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
	for _, n := range projectsInfo {
		version, ok := config[n.short+"_VERSION"]

		if !ok {
			log.Fatalf("ERROR: config doesn't contain '%s_VERSION' for %s", n.short, n.name)
		}

		projects = append(projects, &project{
			short:   n.short,
			name:    n.name,
			version: version,

			repoDir:        baseDir + n.name,
			repoUrl:        fmt.Sprintf("https://github.com/%s/%s.git", org, n.name),
			metricsDocPath: n.metricsDocPath,
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

	r.writeHeader()

	metrics := r.parseMetrics()

	r.writeTOC(metrics)
	r.writeMetrics(metrics)
}

func createFile() *os.File {
	file, err := os.Create("docs/metrics.md")
	if err != nil {
		log.Fatalf("ERROR: creating output file: %s", err)
	}
	return file
}

func (r *releaseData) writeHeader() {
	r.outFile.WriteString("# KubeVirt components metrics\n\n")
	r.outFile.WriteString("This document aims to help users that are not familiar with metrics exposed by all the KubeVirt components.\n")
	r.outFile.WriteString("All metrics documented here are auto-generated in each component repository and gathered here.\n")
	r.outFile.WriteString("They reflect and describe exactly what is being exposed.\n\n")
}

func (r *releaseData) parseMetrics() map[string]string {
	metrics := make(map[string]string)

	for _, p := range r.projects {
		content, err := readLines(path.Join(p.repoDir, "/", p.metricsDocPath))
		if err != nil {
			log.Printf("WARNING: %s project does not contain any metrics documentation in '%s'", p.name, p.version)
		}

		parsed := p.parseMetricsDoc(content)
		if len(parsed) != 0 {
			metrics[p.name] = parsed
		}
	}

	return metrics
}

func (r *releaseData) writeTOC(metrics map[string]string) {
	r.outFile.WriteString("## Table of Contents\n\n")

	for _, p := range r.projects {
		if _, ok := metrics[p.name]; ok {
			r.outFile.WriteString("- [" + p.name + "](#" + p.name + ")\n")
		}
	}

	r.outFile.WriteString("\n")
}

func (r *releaseData) writeMetrics(metrics map[string]string) {
	for _, p := range r.projects {
		if content, ok := metrics[p.name]; ok {
			p.writeComponentMetrics(r.outFile, content)
		}
	}
}

func (p *project) writeComponentMetrics(outFile *os.File, content string) {
	outFile.WriteString("<div id='" + p.name + "'></div>\n\n")

	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/kubevirt/%s/releases/tags/%s", p.name, p.version))

	if err != nil || resp.StatusCode != 200 {
		outFile.WriteString(fmt.Sprintf("## [%s](https://github.com/kubevirt/%s/tree/%s)\n\n", p.name, p.name, p.version))
	} else {
		outFile.WriteString(fmt.Sprintf("## [%s - %s](https://github.com/kubevirt/%s/releases/tag/%s)\n\n", p.name, p.version, p.name, p.version))
	}

	outFile.WriteString(content)
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

func (p *project) parseMetricsDoc(content []string) string {
	running := false
	var b strings.Builder

	for _, s := range content {
		if !running && isBeginning(s) {
			running = true
		} else if running {
			if isEnd(s) {
				break
			}
			fmt.Fprintln(&b, s)
		}
	}

	return b.String()
}

func isBeginning(s string) bool {
	return strings.HasPrefix(s, "## ") && strings.Contains(s, "Metrics List")
}

func isEnd(s string) bool {
	return strings.HasPrefix(s, "## ")
}
