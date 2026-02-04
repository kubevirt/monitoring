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
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/grafana/regexp"
	"k8s.io/klog/v2"
)

var (
	runbookRegex = regexp.MustCompile(`.*\.md`)

	namespaceRegex                = regexp.MustCompile(`(namespace:|-n|--namespace) (kubevirt(?:-hyperconverged)?)`)
	downstreamCommentsRegex       = regexp.MustCompile(`(?s)<!--DS: (.*?)-->`)
	multipleNewLinesRegex         = regexp.MustCompile(`\n\n+`)
	asciiDocLinksRegex            = regexp.MustCompile(`link:(https://[^[\]]+)\[([^[\]]+)\]`)
	trailingSpaceAtEndOfLineRegex = regexp.MustCompile(`[ ]+\n`)
)

type runbook struct {
	name string

	lastLocalUpdate     time.Time
	upstreamLastUpdated time.Time
}

func listRunbooksThatNeedUpdate(downstreamRepo *git.Repository, upstreamRepo *git.Repository) ([]runbook, []runbook) {
	localRunbooks, err := findRunbooksLastUpdateDates(downstreamRepo, downstreamRunbooksDir)
	if err != nil {
		klog.Fatal(fmt.Errorf("failed to find runbooks last update dates: %w", err))
	}

	upstreamRunbooks, err := findRunbooksLastUpdateDates(upstreamRepo, upstreamRunbooksDir)
	if err != nil {
		klog.Fatal(fmt.Errorf("failed to find runbooks last update dates: %w", err))
	}

	return checkWhichRunbooksNeedUpdate(localRunbooks, upstreamRunbooks), checkWhichRunbooksNeedDeprecation(localRunbooks, upstreamRunbooks)
}

func findRunbooksLastUpdateDates(repo *git.Repository, dir string) (map[string]time.Time, error) {
	runbooksTree, err := getDirCurrentTree(repo, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get current runbooks tree: %w", err)
	}

	runbooks := make(map[string]time.Time)
	for _, entry := range runbooksTree.Entries {
		if !runbookRegex.MatchString(entry.Name) {
			continue
		}

		lastCommitDate, err := findLastCommitDate(repo, dir, entry.Name, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to find last commit date: %w", err)
		}
		runbooks[entry.Name] = lastCommitDate
	}

	return runbooks, nil
}

func checkWhichRunbooksNeedUpdate(localRunbooks, upstreamRunbooks map[string]time.Time) []runbook {
	var runbooksToUpdate []runbook

	for name, lastUpstreamUpdate := range upstreamRunbooks {
		lastLocalUpdate, ok := localRunbooks[name]

		if !ok {
			// minimal time to be considered as last update
			lastLocalUpdate = time.UnixMilli(0)
		}

		if lastLocalUpdate.Before(lastUpstreamUpdate) {
			runbooksToUpdate = append(runbooksToUpdate, runbook{
				name:                name,
				lastLocalUpdate:     lastLocalUpdate,
				upstreamLastUpdated: lastUpstreamUpdate,
			})
		}
	}

	return runbooksToUpdate
}

func checkWhichRunbooksNeedDeprecation(localRunbooks, upstreamRunbooks map[string]time.Time) []runbook {
	var runbooksToDeprecate []runbook

	for name, lastLocalUpdate := range localRunbooks {
		_, ok := upstreamRunbooks[name]
		if !ok {
			runbooksToDeprecate = append(runbooksToDeprecate, runbook{
				name:                name,
				lastLocalUpdate:     lastLocalUpdate,
				upstreamLastUpdated: time.Time{},
			})
		}
	}

	return runbooksToDeprecate
}

func copyRunbook(name string) error {
	from := path.Join(upstreamCloneDir, upstreamRunbooksDir, name)
	to := path.Join(downstreamCloneDir, downstreamRunbooksDir, name)

	file, err := os.ReadFile(from)
	if err != nil {
		return fmt.Errorf("failed to read runbook %s: %w", name, err)
	}

	content := replaceContents(string(file))
	return createAndWriteFile(to, content)
}

func replaceContents(content string) string {
	// Replace all 'kubectl' with 'oc'
	content = strings.ReplaceAll(content, "kubectl", "oc")

	// Replace all namespaces
	content = namespaceRegex.ReplaceAllString(content, "$1 openshift-cnv")

	// Remove all US comments
	content = removeTextBetweenTags(content, "<!--USstart-->", "<!--USend-->")

	// Uncomment DS comment - <!--DS: <content>-->
	content = downstreamCommentsRegex.ReplaceAllString(content, "$1")

	// Replace 'KubeVirt' with 'OpenShift Virtualization' when not in code blocks
	content = replaceOnlyInText(content, "KubeVirt", "OpenShift Virtualization")

	// Replace 'Kubernetes' with 'OpenShift Container Platform' when not in code blocks
	content = replaceOnlyInText(content, "Kubernetes", "OpenShift Container Platform")

	// Replace AsciiDoc links with Markdown links
	content = asciiDocLinksRegex.ReplaceAllString(content, "[$2]($1)")

	// Remove multiple (2+) new lines
	content = multipleNewLinesRegex.ReplaceAllString(content, "\n\n")

	content = wrapLines(content, 80)

	// Remove trailing spaces from all lines
	content = trailingSpaceAtEndOfLineRegex.ReplaceAllString(content, "\n")

	// Keep only one empty line at the end of the file
	content = strings.TrimRight(content, "\n")

	return content
}

func wrapLines(content string, maxLineLength int) string {
	var result strings.Builder
	lines := strings.Split(content, "\n")

	inBlockCode := false
	inInlineCode := false

	for _, line := range lines {
		lineLength := 0
		words := strings.SplitAfter(line, " ")

		for _, word := range words {
			if strings.HasPrefix(word, "```") {
				inBlockCode = true
			} else if strings.HasPrefix(word, "`") {
				inInlineCode = true
			}

			if lineLength+len(word) > maxLineLength && !inBlockCode && !inInlineCode {
				result.WriteString("\n")
				lineLength = 0
			}

			if strings.HasSuffix(word, "```") || strings.HasSuffix(word, "``` ") {
				inBlockCode = false
			} else if strings.HasSuffix(word, "`") || strings.HasSuffix(word, "` ") {
				inInlineCode = false
			}

			result.WriteString(word)
			lineLength += len(word)
		}

		result.WriteString("\n")
	}

	return result.String()
}

func removeTextBetweenTags(content, startTag, endTag string) string {
	var result strings.Builder
	startIndex := 0

	for {
		// Find start and end index of the tags
		startIndex = strings.Index(content, startTag)
		endIndex := strings.Index(content, endTag)

		// If both tags exist, remove the text between them
		if startIndex != -1 && endIndex != -1 {
			result.WriteString(content[:startIndex])
			content = content[endIndex+len(endTag):]
		} else {
			result.WriteString(content)
			break
		}
	}

	return result.String()
}

func replaceOnlyInText(text string, old string, new string) string {
	var result strings.Builder
	lines := strings.Split(text, "\n")

	isTitleLine := false
	inBlockCode := false
	inInlineCode := false

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			isTitleLine = true
		} else {
			isTitleLine = false
		}

		words := strings.SplitAfter(line, " ")

		for _, word := range words {
			trimmedWord := strings.TrimSpace(word)

			if !inBlockCode && !inInlineCode {
				if strings.HasPrefix(trimmedWord, "```") {
					inBlockCode = true
				} else if strings.HasPrefix(trimmedWord, "`") {
					inInlineCode = true
				}
			}

			if !inBlockCode && !inInlineCode && !isTitleLine && strings.Contains(word, old) {
				word = strings.ReplaceAll(word, old, new)
				if new == "" {
					word = strings.TrimSuffix(word, " ")
				}
			}

			if inBlockCode && strings.HasSuffix(trimmedWord, "```") {
				inBlockCode = false
			} else if inInlineCode && strings.HasSuffix(trimmedWord, "`") {
				inInlineCode = false
			}

			result.WriteString(word)
		}

		result.WriteString("\n")
	}

	return result.String()
}
func createAndWriteFile(path, content string) error {
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		_, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}
