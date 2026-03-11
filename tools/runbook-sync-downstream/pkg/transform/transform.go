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

package transform

import (
	"strings"

	"github.com/grafana/regexp"
)

var (
	namespaceRegex                = regexp.MustCompile(`(namespace:|-n|--namespace) (kubevirt(?:-hyperconverged)?)`)
	downstreamCommentsRegex       = regexp.MustCompile(`(?s)<!--DS: (.*?)-->`)
	multipleNewLinesRegex         = regexp.MustCompile(`\n\n+`)
	asciiDocLinksRegex            = regexp.MustCompile(`link:(https://[^[\]]+)\[([^[\]]+)\]`)
	trailingSpaceAtEndOfLineRegex = regexp.MustCompile(`[ ]+\n`)
)

func ReplaceContents(content string) string {
	content = strings.ReplaceAll(content, "kubectl", "oc")

	content = namespaceRegex.ReplaceAllString(content, "$1 openshift-cnv")

	content = removeTextBetweenTags(content, "<!--USstart-->", "<!--USend-->")

	content = downstreamCommentsRegex.ReplaceAllString(content, "$1")

	content = replaceOnlyInText(content, "KubeVirt", "OpenShift Virtualization")

	content = replaceOnlyInText(content, "Kubernetes", "OpenShift Container Platform")

	content = asciiDocLinksRegex.ReplaceAllString(content, "[$2]($1)")

	content = multipleNewLinesRegex.ReplaceAllString(content, "\n\n")

	content = wrapLines(content, 80)

	content = trailingSpaceAtEndOfLineRegex.ReplaceAllString(content, "\n")

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
		startIndex = strings.Index(content, startTag)
		endIndex := strings.Index(content, endTag)

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
