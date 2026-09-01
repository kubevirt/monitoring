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
	"bytes"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/grafana/regexp"
	"k8s.io/klog/v2"

	"github.com/kubevirt/monitoring/tools/runbook-sync-downstream/pkg/transform"
)

var (
	runbookRegex = regexp.MustCompile(`.*\.md`)
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

	content := transform.ReplaceContents(string(file))
	return createAndWriteFile(to, content)
}

// generatedFileMatches reports whether the freshly generated runbook file at
// repoDir/filePath is byte-identical to want (the content already on the PR
// branch). It returns an error if the generated file cannot be read.
func generatedFileMatches(repoDir, filePath string, want []byte) (bool, error) {
	got, err := os.ReadFile(path.Join(repoDir, filePath))
	if err != nil {
		return false, fmt.Errorf("failed to read generated file %s: %w", filePath, err)
	}

	return bytes.Equal(got, want), nil
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
