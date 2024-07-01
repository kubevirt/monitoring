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
	"path"
	"time"

	"github.com/go-git/go-git/v5"
)

func findLastCommitDate(repo *git.Repository, dirname, filename string, since *time.Time) (time.Time, error) {
	filepath := path.Join(dirname, filename)
	logOptions := &git.LogOptions{
		FileName: &filepath,
	}
	if since != nil {
		logOptions.Since = since
	}

	commitIter, err := repo.Log(logOptions)
	if err != nil {
		return time.Now(), fmt.Errorf("failed to get log: %w", err)
	}

	commit, err := commitIter.Next()
	if err != nil {
		return time.Now(), fmt.Errorf("failed to get next commit: %w", err)
	}

	lastCommitDate := commit.Committer.When.In(time.UTC)

	return lastCommitDate, nil
}
