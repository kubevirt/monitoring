package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/perses/perses/go-sdk/dashboard"
	"github.com/stolostron/multicluster-observability-addon/pkg/perses/dashboards/virtualization"

	"github.com/kubevirt/monitoring/pkg/dashboards/transform"

	k8syaml "sigs.k8s.io/yaml"
)

type persesResource struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Spec       json.RawMessage   `json:"spec"`
}

type builderEntry struct {
	name         string
	fn           func(string, string) (dashboard.Builder, error)
	multicluster bool // true = source is multicluster; apply ToSingleCluster transform
}

func includedDashboardNames(builders []builderEntry, prefix string) []string {
	names := make([]string, len(builders))
	for i, b := range builders {
		names[i] = prefix + "-" + b.name
	}
	return names
}

func (b builderEntry) generateSpec(project string, cfg transform.Config) (json.RawMessage, error) {
	db, err := b.fn(project, "rbac-query-proxy-datasource")
	if err != nil {
		return nil, fmt.Errorf("building dashboard: %w", err)
	}
	specJSON, err := json.Marshal(db.Dashboard.Spec)
	if err != nil {
		return nil, fmt.Errorf("marshalling spec: %w", err)
	}
	if b.multicluster {
		return transform.ToSingleCluster(specJSON, cfg)
	}
	return specJSON, nil
}

func wrapAsResource(name, namespace string, spec json.RawMessage) ([]byte, error) {
	resource := persesResource{
		APIVersion: "perses.dev/v1alpha1",
		Kind:       "PersesDashboard",
		Metadata: map[string]string{
			"name":      name,
			"namespace": namespace,
		},
		Spec: spec,
	}

	resourceJSON, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("marshalling resource: %w", err)
	}
	return k8syaml.JSONToYAML(resourceJSON)
}

func main() {
	prefix := flag.String("prefix", "cnv", "dashboard ID prefix")
	displayPrefix := flag.String("display-prefix", "OpenShift Virtualization", "display name prefix (replaces 'Virtualization')")
	outDir := flag.String("out", "dashboards/perses", "output directory for generated YAML")
	namespace := flag.String("namespace", "openshift-cnv", "namespace for PersesDashboard CRs")
	flag.Parse()

	builders := []builderEntry{
		{"node-memory-overview", virtualization.BuildNodeMemoryOverview, true},
		{"virtual-machines-inventory", virtualization.BuildVMInventory, true},
		{"virtual-machines-utilization", virtualization.BuildVMUtilization, true},
		{"virtual-machines-service-level", virtualization.BuildVMServiceLevel, true},
		{"virtual-machines-by-time-in-status", virtualization.BuildVMByTimeInStatus, true},
		{"virtual-machines-top-consumers", virtualization.BuildTopConsumers, true},
	}

	cfg := transform.Config{
		DashboardPrefix:    *prefix,
		DisplayNamePrefix:  *displayPrefix,
		IncludedDashboards: includedDashboardNames(builders, *prefix),
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	for _, b := range builders {
		outputSpec, err := b.generateSpec(*namespace, cfg)
		if err != nil {
			log.Fatalf("Failed to generate %s: %v", b.name, err)
		}

		dashboardName := *prefix + "-" + b.name
		yamlBytes, err := wrapAsResource(dashboardName, *namespace, outputSpec)
		if err != nil {
			log.Fatalf("Failed to generate YAML for %s: %v", b.name, err)
		}

		outFile := filepath.Join(*outDir, dashboardName+".yaml")
		if err := os.WriteFile(outFile, yamlBytes, 0o644); err != nil {
			log.Fatalf("Failed to write %s: %v", outFile, err)
		}
		fmt.Printf("Generated: %s\n", outFile)
	}
}
