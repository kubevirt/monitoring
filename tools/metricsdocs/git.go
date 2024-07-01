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
	"log"
	"os"
	"os/exec"
)

func (p *project) gitCheckoutUpstream() error {
	_, err := os.Stat(p.repoDir)
	if err == nil {
		_, err := gitCommand("-C", p.repoDir, "status")
		if err == nil {
			// checkout already exists, updating
			return p.gitUpdateFromUpstream()
		}
	}

	return p.gitCloneUpstream()
}

func (p *project) gitUpdateFromUpstream() error {
	_, err := gitCommand("-C", p.repoDir, "checkout", "main")
	if err != nil {
		log.Printf("INFO: repository %s does not have 'main' branch, falling back to 'master'", p.name)
		_, err = gitCommand("-C", p.repoDir, "checkout", "master")
		if err != nil {
			return err
		}
	}

	_, err = gitCommand("-C", p.repoDir, "pull")
	if err != nil {
		return err
	}
	return nil
}

func (p *project) gitCloneUpstream() error {
	// start fresh because checkout doesn't exist or is corrupted
	os.RemoveAll(p.repoDir)
	err := os.MkdirAll(p.repoDir, 0755)
	if err != nil {
		return err
	}

	// add upstream remote branch
	_, err = gitCommand("clone", p.repoUrl, p.repoDir)
	if err != nil {
		return err
	}

	_, err = gitCommand("-C", p.repoDir, "config", "diff.renameLimit", "999999")
	if err != nil {
		return err
	}

	return nil
}

func (p *project) gitSwitchToBranch(branch string) error {
	_, err := gitCommand("-C", p.repoDir, "checkout", branch)
	if err != nil {
		return err
	}

	return nil
}

func gitCommand(arg ...string) (string, error) {
	log.Printf("executing 'git %v", arg)
	cmd := exec.Command("git", arg...)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("ERROR: git command output: %s : %s ", string(bytes), err)
		return "", err
	}
	return string(bytes), nil
}
