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

// codeScanner tracks whether a word-by-word scan of a Markdown document is
// currently inside a fenced code block or an inline code span. Block state is
// deliberately carried across lines, since a fenced block spans several.
type codeScanner struct {
	inBlockCode  bool
	inInlineCode bool

	// fence is the marker that opened the current block, so that a shorter
	// run of backticks inside the block does not close it.
	fence       string
	lineStarted bool
	atLineStart bool
	openedHere  bool
}

func (c *codeScanner) inCode() bool {
	return c.inBlockCode || c.inInlineCode
}

func (c *codeScanner) beginLine() {
	c.lineStarted = false
}

func (c *codeScanner) enter(word string) {
	c.atLineStart = word != "" && !c.lineStarted
	if c.atLineStart {
		c.lineStarted = true
	}
	c.openedHere = false

	if c.inCode() {
		return
	}

	if c.atLineStart {
		if marker := fenceMarker(word); marker != "" {
			c.inBlockCode = true
			c.fence = marker
			c.openedHere = true

			return
		}
	}

	if strings.HasPrefix(word, "`") {
		c.inInlineCode = true
	}
}

func (c *codeScanner) exit(word string) {
	if c.inBlockCode {
		if !c.openedHere && c.atLineStart && closesFence(word, c.fence) {
			c.inBlockCode = false
			c.fence = ""
		}

		return
	}

	// Trim trailing punctuation so a closing backtick immediately followed
	// by e.g. "." or "," (as in "`virt-api`,") is still detected as the end
	// of an inline code span. Otherwise the state stays stuck open for the
	// rest of the document: replaceOnlyInText then leaves every later term
	// unconverted, and wrapLines leaves every later line unwrapped.
	if c.inInlineCode && strings.HasSuffix(strings.TrimRight(word, ".,;:!?)"), "`") {
		c.inInlineCode = false
	}
}

// fenceMarker returns the run of backticks that opens a fenced code block
func fenceMarker(word string) string {
	length := 0
	for length < len(word) && word[length] == '`' {
		length++
	}

	if length < 3 || strings.Contains(word[length:], "`") {
		return ""
	}

	return word[:length]
}

func closesFence(word, fence string) bool {
	if fence == "" || len(word) < len(fence) {
		return false
	}

	return word == strings.Repeat("`", len(word))
}

func wrapLines(content string, maxLineLength int) string {
	var result strings.Builder
	lines := strings.Split(content, "\n")

	var code codeScanner

	for _, line := range lines {
		code.beginLine()

		lineLength := 0
		words := strings.SplitAfter(line, " ")

		for _, word := range words {
			trimmedWord := strings.TrimSpace(word)

			code.enter(trimmedWord)

			if lineLength+len(word) > maxLineLength && !code.inCode() {
				result.WriteString("\n")
				lineLength = 0
			}

			code.exit(trimmedWord)

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

	var code codeScanner

	for _, line := range lines {
		code.beginLine()

		if strings.HasPrefix(line, "#") {
			isTitleLine = true
		} else {
			isTitleLine = false
		}

		words := strings.SplitAfter(line, " ")

		for _, word := range words {
			trimmedWord := strings.TrimSpace(word)

			code.enter(trimmedWord)

			if !code.inCode() && !isTitleLine && strings.Contains(word, old) {
				word = strings.ReplaceAll(word, old, new)
				if new == "" {
					word = strings.TrimSuffix(word, " ")
				}
			}

			code.exit(trimmedWord)

			result.WriteString(word)
		}

		result.WriteString("\n")
	}

	return result.String()
}
