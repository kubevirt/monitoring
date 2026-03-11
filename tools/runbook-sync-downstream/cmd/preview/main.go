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
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kubevirt/monitoring/tools/runbook-sync-downstream/pkg/transform"
)

func main() {
	outputDir := flag.String("output-dir", "", "directory to write transformed runbook files (defaults to stdout)")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: runbook-preview [-output-dir DIR] FILE [FILE...]")
		os.Exit(1)
	}

	if *outputDir != "" {
		if err := os.MkdirAll(*outputDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create output directory: %v\n", err)
			os.Exit(1)
		}
	}

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read %s: %v\n", f, err)
			os.Exit(1)
		}

		transformed := transform.ReplaceContents(string(content))

		if *outputDir != "" {
			outPath := filepath.Join(*outputDir, filepath.Base(f))
			if err := os.WriteFile(outPath, []byte(transformed+"\n"), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", outPath, err)
				os.Exit(1)
			}
			fmt.Printf("wrote %s\n", outPath)
		} else {
			fmt.Printf("--- %s (downstream preview) ---\n%s\n\n", f, transformed)
		}
	}
}
