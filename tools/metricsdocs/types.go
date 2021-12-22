package main

import (
	"os"
)

var projectsInfo = []*projectInfo{
	{"KUBEVIRT", "kubevirt", "docs/metrics.md"},
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

type releaseData struct {
	org      string
	projects []*project

	outFile *os.File
}
