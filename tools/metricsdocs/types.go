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
	"os"
)

var projectsInfo = []*projectInfo{
	{"KUBEVIRT", "kubevirt", "docs/observability/metrics.md"},
	{"CDI", "containerized-data-importer", "doc/metrics.md"},
	{"NETWORK_ADDONS", "cluster-network-addons-operator", "docs/metrics.md"},
	{"SSP", "ssp-operator", "docs/metrics.md"},
	{"NMO", "node-maintenance-operator", "docs/metrics.md"},
	{"HPPO", "hostpath-provisioner-operator", "docs/metrics.md"},
	{"HPP", "hostpath-provisioner", "docs/metrics.md"},
	{"HCO", "hyperconverged-cluster-operator", "docs/metrics.md"},
}

type projectInfo struct {
	short          string
	name           string
	metricsDocPath string
}

type project struct {
	short   string
	name    string
	version string

	repoDir        string
	repoUrl        string
	metricsDocPath string
}

type Metric struct {
	Operator    string
	Name        string
	Kind        string
	Type        string
	Description string
}

type TemplateOperator struct {
	Name string
	Link string
}

type releaseData struct {
	org      string
	projects []*project

	outFile *os.File
}
