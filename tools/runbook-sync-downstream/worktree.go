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
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func getDirCurrentTree(repo *git.Repository, dir string) (*object.Tree, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	commitTree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit tree: %w", err)
	}

	runbooksTree, err := commitTree.Tree(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get runbooks tree: %w", err)
	}

	return runbooksTree, nil
}

// forkBranchFileContent fetches branchName from the fork remote and returns the
// content of filePath as it currently exists on that branch. It returns
// (nil, nil) when the branch or file is absent on the fork, so callers can treat
// a missing remote file as "changed".
func forkBranchFileContent(repo *git.Repository, branchName, filePath string) ([]byte, error) {
	refSpec := config.RefSpec(fmt.Sprintf("+refs/heads/%s:refs/remotes/%s/%s", branchName, forkRemoteName, branchName))

	err := repo.Fetch(&git.FetchOptions{
		RemoteName: forkRemoteName,
		RefSpecs:   []config.RefSpec{refSpec},
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		if errors.Is(err, git.NoMatchingRefSpecError{}) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch fork branch %s: %w", branchName, err)
	}

	ref, err := repo.Reference(plumbing.NewRemoteReferenceName(forkRemoteName, branchName), true)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve fork branch %s: %w", branchName, err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit for fork branch %s: %w", branchName, err)
	}

	f, err := commit.File(filePath)
	if err != nil {
		if errors.Is(err, object.ErrFileNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read %s from fork branch %s: %w", filePath, branchName, err)
	}

	content, err := f.Contents()
	if err != nil {
		return nil, fmt.Errorf("failed to read contents of %s from fork branch %s: %w", filePath, branchName, err)
	}

	return []byte(content), nil
}

func newBranchFromMain(repo *git.Repository, name string) (*git.Worktree, error) {
	w, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(downstreamMainBranch),
		Create: false,
		Keep:   false,
		Force:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to checkout to main: %w", err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
		Create: true,
		Keep:   false,
		Force:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	return w, nil
}
