package transform

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestStripClusterFromQuery(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple cluster matcher",
			input: `kubevirt_vm_info{cluster=~"$cluster", status_group="running"}`,
			want:  `kubevirt_vm_info{status_group="running"}`,
		},
		{
			name:  "cluster matcher only",
			input: `kubevirt_vm_info{cluster=~"$cluster"}`,
			want:  `kubevirt_vm_info`,
		},
		{
			name:  "cluster in middle of matchers",
			input: `kubevirt_vmi_info{namespace=~"$namespace",cluster=~"$cluster",phase="running"}`,
			want:  `kubevirt_vmi_info{namespace=~"$namespace",phase="running"}`,
		},
		{
			name:  "sum by cluster",
			input: `sum by (cluster)(kubevirt_vm_info)`,
			want:  `sum(kubevirt_vm_info)`,
		},
		{
			name:  "by cluster and other labels",
			input: `sum by (cluster, namespace, name)(kubevirt_vm_info)`,
			want:  `sum by (namespace, name)(kubevirt_vm_info)`,
		},
		{
			name:  "on(cluster,name,namespace)",
			input: `on(cluster,name,namespace) group_left()`,
			want:  `on(name,namespace) group_left()`,
		},
		{
			name:  "on(cluster) group_left()",
			input: `sum(x) on(cluster) group_left() sum(y)`,
			want:  `sum(x)  sum(y)`,
		},
		{
			name:  "metric rename - health status",
			input: `kubevirt_hyperconverged_operator_health_status{namespace="test"}`,
			want:  `cluster:kubevirt_hco_operator_health_status:count{namespace="test"}`,
		},
		{
			name:  "metric rename - memory used",
			input: `kubevirt_vmi_memory_used_bytes{node=~"$node"}`,
			want:  `vmi:kubevirt_vmi_memory_used_bytes:sum{node=~"$node"}`,
		},
		{
			name:  "metric rename - phase count",
			input: `kubevirt_vmi_phase_count{phase="running"}`,
			want:  `node:kubevirt_vmi_phase:sum{phase="running"}`,
		},
		{
			name:  "name=~$cluster removal",
			input: `acm_managed_cluster_labels{name=~"$cluster"}`,
			want:  `acm_managed_cluster_labels`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripClusterFromQuery(tt.input)
			if got != tt.want {
				t.Errorf("stripClusterFromQuery(%q)\n  got:  %q\n  want: %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToSingleCluster_RemovesClusterVariable(t *testing.T) {
	spec := map[string]any{
		"variables": []any{
			map[string]any{
				"kind": "ListVariable",
				"spec": map[string]any{
					"name": "cluster",
				},
			},
			map[string]any{
				"kind": "ListVariable",
				"spec": map[string]any{
					"name": "namespace",
				},
			},
		},
		"panels":  map[string]any{},
		"layouts": []any{},
	}

	specJSON, _ := json.Marshal(spec)
	cfg := Config{
		DashboardPrefix: "cnv",
	}

	result, err := ToSingleCluster(specJSON, cfg)
	if err != nil {
		t.Fatalf("ToSingleCluster returned error: %v", err)
	}

	var resultSpec map[string]any
	if err := json.Unmarshal(result, &resultSpec); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	vars := resultSpec["variables"].([]any)
	if len(vars) != 1 {
		t.Fatalf("Expected 1 variable, got %d", len(vars))
	}
	varSpec := vars[0].(map[string]any)["spec"].(map[string]any)
	if varSpec["name"] != "namespace" {
		t.Errorf("Expected remaining variable to be 'namespace', got %q", varSpec["name"])
	}
}

func TestToSingleCluster_NoClusterInOutput(t *testing.T) {
	spec := map[string]any{
		"variables": []any{
			map[string]any{
				"kind": "ListVariable",
				"spec": map[string]any{
					"name": "cluster",
					"plugin": map[string]any{
						"kind": "PrometheusLabelValuesVariable",
						"spec": map[string]any{
							"matchers": []any{"kubevirt_vm_info"},
						},
					},
				},
			},
		},
		"panels": map[string]any{
			"0_0": map[string]any{
				"spec": map[string]any{
					"queries": []any{
						map[string]any{
							"spec": map[string]any{
								"plugin": map[string]any{
									"kind": "PrometheusTimeSeriesQuery",
									"spec": map[string]any{
										"query": `sum by (cluster)(kubevirt_vm_info{cluster=~"$cluster"})`,
									},
								},
							},
						},
					},
				},
			},
		},
		"layouts": []any{},
	}

	specJSON, _ := json.Marshal(spec)
	cfg := Config{
		DashboardPrefix: "cnv",
	}

	result, err := ToSingleCluster(specJSON, cfg)
	if err != nil {
		t.Fatalf("ToSingleCluster returned error: %v", err)
	}

	resultStr := string(result)
	if strings.Contains(resultStr, "$cluster") {
		t.Error("Output still contains $cluster")
	}
	if strings.Contains(resultStr, `"name":"cluster"`) {
		t.Error("Output still contains cluster variable")
	}
}

func TestToSingleCluster_PreservesNonClusterSpec(t *testing.T) {
	spec := map[string]any{
		"variables": []any{
			map[string]any{
				"kind": "ListVariable",
				"spec": map[string]any{
					"name": "namespace",
					"plugin": map[string]any{
						"kind": "PrometheusLabelValuesVariable",
						"spec": map[string]any{
							"labelName": "namespace",
						},
					},
				},
			},
		},
		"panels": map[string]any{
			"0_0": map[string]any{
				"spec": map[string]any{
					"queries": []any{
						map[string]any{
							"spec": map[string]any{
								"plugin": map[string]any{
									"kind": "PrometheusTimeSeriesQuery",
									"spec": map[string]any{
										"query": `sum by (namespace)(kubevirt_vm_info{namespace=~"$namespace"})`,
									},
								},
							},
						},
					},
				},
			},
		},
		"layouts": []any{},
	}

	specJSON, _ := json.Marshal(spec)
	cfg := Config{DashboardPrefix: "cnv"}

	result, err := ToSingleCluster(specJSON, cfg)
	if err != nil {
		t.Fatalf("ToSingleCluster returned error: %v", err)
	}

	var resultSpec map[string]any
	if err := json.Unmarshal(result, &resultSpec); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	panels := resultSpec["panels"].(map[string]any)
	panel := panels["0_0"].(map[string]any)
	queries := panel["spec"].(map[string]any)["queries"].([]any)
	query := queries[0].(map[string]any)["spec"].(map[string]any)["plugin"].(map[string]any)["spec"].(map[string]any)["query"].(string)

	want := `sum by (namespace)(kubevirt_vm_info{namespace=~"$namespace"})`
	if query != want {
		t.Errorf("query was modified\n  got:  %q\n  want: %q", query, want)
	}

	vars := resultSpec["variables"].([]any)
	if len(vars) != 1 {
		t.Fatalf("Expected 1 variable, got %d", len(vars))
	}
}

func TestToSingleCluster_RemovesDatasource(t *testing.T) {
	spec := map[string]any{
		"variables": []any{
			map[string]any{
				"kind": "ListVariable",
				"spec": map[string]any{
					"name": "namespace",
					"plugin": map[string]any{
						"kind": "PrometheusLabelValuesVariable",
						"spec": map[string]any{
							"datasource": map[string]any{"kind": "PrometheusDatasource", "name": "rbac-query-proxy-datasource"},
							"labelName":  "namespace",
						},
					},
				},
			},
		},
		"panels": map[string]any{
			"0_0": map[string]any{
				"spec": map[string]any{
					"queries": []any{
						map[string]any{
							"spec": map[string]any{
								"plugin": map[string]any{
									"kind": "PrometheusTimeSeriesQuery",
									"spec": map[string]any{
										"query":      "up",
										"datasource": map[string]any{"kind": "PrometheusDatasource", "name": "rbac-query-proxy-datasource"},
									},
								},
							},
						},
					},
				},
			},
		},
		"layouts": []any{},
	}

	specJSON, _ := json.Marshal(spec)
	cfg := Config{
		DashboardPrefix: "cnv",
	}

	result, err := ToSingleCluster(specJSON, cfg)
	if err != nil {
		t.Fatalf("ToSingleCluster returned error: %v", err)
	}

	if strings.Contains(string(result), "datasource") {
		t.Error("Output still contains datasource references")
	}
}
