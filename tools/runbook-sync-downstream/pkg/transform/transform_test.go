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

import "testing"

// replaceOnlyInText re-joins lines with a trailing "\n", so its output always
// ends with a newline. The inputs below intentionally omit the trailing newline
// to keep the expectations focused on the replacement behavior rather than that
// artifact (which ReplaceContents later trims away).
func TestReplaceOnlyInText(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "replaces term in plain text",
			in:   "KubeVirt objects cannot send API calls.",
			want: "OpenShift Virtualization objects cannot send API calls.\n",
		},
		{
			name: "does not replace term in a title line",
			in:   "# KubeVirt is great",
			want: "# KubeVirt is great\n",
		},
		{
			name: "does not replace term inside inline code",
			in:   "Run the `KubeVirt` command.",
			want: "Run the `KubeVirt` command.\n",
		},
		{
			// Regression test: the punctuation-trimming rule must apply to
			// inline code only. A fenced-code line ending in "```." is code
			// content, not a valid closing delimiter, so block code stays open
			// and "KubeVirt" inside it must be left untouched.
			name: "keeps block code open when a fence line ends in punctuation",
			in: "```bash\n" +
				"KubeVirt```.\n" +
				"KubeVirt\n" +
				"```",
			want: "```bash\n" +
				"KubeVirt```.\n" +
				"KubeVirt\n" +
				"```\n",
		},
		{
			// Regression test: a closing backtick followed by punctuation
			// ("`virt-api-.*`.") used to leave the inline-code state stuck
			// open, so no later term was ever replaced.
			name: "replaces term after inline code whose closing backtick is followed by punctuation",
			in: "The recording rule counts pods matching `virt-api-.*`.\n" +
				"KubeVirt objects cannot send API calls.",
			want: "The recording rule counts pods matching `virt-api-.*`.\n" +
				"OpenShift Virtualization objects cannot send API calls.\n",
		},
		{
			// End-to-end regression against the reported VirtAPIDown runbook:
			// the "## Impact" line was left unconverted because of the stuck
			// inline-code state introduced earlier by "`virt-api-.*`.".
			name: "replaces term in the Impact section of a realistic runbook",
			in: "# VirtAPIDown\n" +
				"\n" +
				"## Meaning\n" +
				"\n" +
				"No running `virt-api` pod has been detected for 10 minutes.\n" +
				"\n" +
				"The recording rule counts pods in `Running` phase matching `virt-api-.*`.\n" +
				"\n" +
				"## Impact\n" +
				"\n" +
				"KubeVirt objects cannot send API calls.",
			want: "# VirtAPIDown\n" +
				"\n" +
				"## Meaning\n" +
				"\n" +
				"No running `virt-api` pod has been detected for 10 minutes.\n" +
				"\n" +
				"The recording rule counts pods in `Running` phase matching `virt-api-.*`.\n" +
				"\n" +
				"## Impact\n" +
				"\n" +
				"OpenShift Virtualization objects cannot send API calls.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceOnlyInText(tt.in, "KubeVirt", "OpenShift Virtualization")
			if got != tt.want {
				t.Errorf("replaceOnlyInText() mismatch\n in:   %q\n got:  %q\n want: %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestWrapLines(t *testing.T) {
	const longSentence = "This line is quite long and absolutely should be wrapped because it exceeds the limit."

	const wrappedSentence = "This line is quite long and absolutely should be wrapped because it exceeds the \n" +
		"limit."

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "wraps plain prose at the limit",
			in:   longSentence,
			want: wrappedSentence + "\n",
		},
		{
			name: "does not wrap inside a fenced code block",
			in: "```bash\n" +
				"$ oc get pods -A -l kubevirt.io=virt-launcher --field-selector=status.phase=Failed --no-headers\n" +
				"```",
			want: "```bash\n" +
				"$ oc get pods -A -l kubevirt.io=virt-launcher --field-selector=status.phase=Failed --no-headers\n" +
				"```\n",
		},
		{
			// Regression test: a closing backtick followed by punctuation
			// ("`virt-api`,") used to leave the inline-code state stuck open.
			// Because inInlineCode suppresses wrapping, every later line in the
			// document was then emitted unwrapped.
			name: "keeps wrapping after inline code whose closing backtick is followed by punctuation",
			in: "Without a ready `virt-api`, requests cannot be served.\n" +
				longSentence,
			want: "Without a ready `virt-api`, requests cannot be served.\n" +
				wrappedSentence + "\n",
		},
		{
			// Regression test: the punctuation-trimming rule must apply to
			// inline code only. A fenced-code line ending in "```." is code
			// content, not a valid closing delimiter, so block code stays open
			// and the long line inside it must not be wrapped.
			name: "keeps block code open when a fence line ends in punctuation",
			in: "```bash\n" +
				"$ echo hello```.\n" +
				longSentence + "\n" +
				"```",
			want: "```bash\n" +
				"$ echo hello```.\n" +
				longSentence + "\n" +
				"```\n",
		},
		{
			// End-to-end regression against the NoReadyVirtAPI runbook: the
			// "## Impact" paragraph was emitted unwrapped because
			// "`virt-api`," in the Meaning section stuck the state open.
			name: "wraps the Impact section of a realistic runbook",
			in: "# NoReadyVirtAPI\n" +
				"\n" +
				"## Meaning\n" +
				"\n" +
				"Without a ready `virt-api`, API requests for virtual machines\n" +
				"cannot be served.\n" +
				"\n" +
				"## Impact\n" +
				"\n" +
				longSentence,
			want: "# NoReadyVirtAPI\n" +
				"\n" +
				"## Meaning\n" +
				"\n" +
				"Without a ready `virt-api`, API requests for virtual machines\n" +
				"cannot be served.\n" +
				"\n" +
				"## Impact\n" +
				"\n" +
				wrappedSentence + "\n",
		},
		{
			// A complete triple-backtick span in the middle of a line is
			// inline code, not a fence: the info string of a real fence may
			// not contain a backtick. Treating it as a fence used to leave
			// block state open and suppress all later wrapping.
			name: "treats a mid-line triple-backtick span as inline code",
			in: "A span such as ```foo```, stays inline.\n" +
				longSentence,
			want: "A span such as ```foo```, stays inline.\n" +
				wrappedSentence + "\n",
		},
		{
			// Trailing text on a line inside a code block must not close it,
			// even when that text ends in three backticks.
			name: "does not close a block on trailing text ending in backticks",
			in: "```bash\n" +
				"$ echo hello```\n" +
				longSentence + "\n" +
				"```",
			want: "```bash\n" +
				"$ echo hello```\n" +
				longSentence + "\n" +
				"```\n",
		},
		{
			// A bare opening fence must not be closed by the very word that
			// opened it.
			name: "handles an opening fence with no info string",
			in: "```\n" +
				longSentence + "\n" +
				"```\n" +
				longSentence,
			want: "```\n" +
				longSentence + "\n" +
				"```\n" +
				wrappedSentence + "\n",
		},
		{
			// A longer opening fence is only closed by a run at least as
			// long, so a shorter run inside the block is content.
			name: "closes a long fence only on a run at least as long",
			in: "````text\n" +
				"```\n" +
				longSentence + "\n" +
				"````\n" +
				longSentence,
			want: "````text\n" +
				"```\n" +
				longSentence + "\n" +
				"````\n" +
				wrappedSentence + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapLines(tt.in, 80)
			if got != tt.want {
				t.Errorf("wrapLines() mismatch\n in:   %q\n got:  %q\n want: %q", tt.in, got, tt.want)
			}
		})
	}
}
