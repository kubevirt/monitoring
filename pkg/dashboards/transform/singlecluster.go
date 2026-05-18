package transform

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strings"
)

var metricRenames = map[string]string{
	"kubevirt_hyperconverged_operator_health_status": "cluster:kubevirt_hco_operator_health_status:count",
	"kubevirt_vmi_memory_used_bytes":                 "vmi:kubevirt_vmi_memory_used_bytes:sum",
	"kubevirt_vmi_phase_count":                       "node:kubevirt_vmi_phase:sum",
}

// Config controls the single-cluster transformation.
type Config struct {
	DashboardPrefix    string
	DisplayNamePrefix  string
	IncludedDashboards []string
}

// ToSingleCluster transforms a multicluster Perses dashboard spec (as JSON)
// into a single-cluster variant by removing cluster references, renaming
// metrics for OCP 4.22+, and updating dashboard metadata.
func ToSingleCluster(specJSON []byte, cfg Config) ([]byte, error) {
	var spec map[string]any
	if err := json.Unmarshal(specJSON, &spec); err != nil {
		return nil, err
	}

	removeClusterVariable(spec)
	rewriteQueries(spec)
	rewriteVariableValues(spec)
	removeClusterTableColumns(spec)
	removeClusterFromTransforms(spec)
	rewriteDashboardLinks(spec, cfg)
	rewriteDisplayName(spec, cfg)
	removeDatasourceRefs(spec)

	return json.Marshal(spec)
}

func removeClusterVariable(spec map[string]any) {
	vars, ok := spec["variables"].([]any)
	if !ok {
		return
	}
	filtered := make([]any, 0, len(vars))
	for _, v := range vars {
		vm, ok := v.(map[string]any)
		if !ok {
			filtered = append(filtered, v)
			continue
		}
		varSpec, _ := vm["spec"].(map[string]any)
		if varSpec == nil {
			filtered = append(filtered, v)
			continue
		}
		if name, _ := varSpec["name"].(string); name == "cluster" {
			continue
		}
		filtered = append(filtered, v)
	}
	spec["variables"] = filtered
}

var (
	// cluster=~"$cluster" with optional surrounding commas/spaces
	reClusterMatcherWithTrailingComma = regexp.MustCompile(`cluster=~"\$cluster"\s*,\s*`)
	reClusterMatcherWithLeadingComma  = regexp.MustCompile(`,\s*cluster=~"\$cluster"`)
	reClusterMatcherAlone             = regexp.MustCompile(`cluster=~"\$cluster"`)
	// Fixes double commas or leading/trailing commas in label matchers
	reDoubleComma   = regexp.MustCompile(`,\s*,`)
	reLeadingComma  = regexp.MustCompile(`\{\s*,`)
	reTrailingComma = regexp.MustCompile(`,\s*\}`)
	// by (cluster, ...) or by (cluster)
	reByClusterPrefix = regexp.MustCompile(`by\s*\(\s*cluster\s*,\s*`)
	reByClusterOnly   = regexp.MustCompile(`\)\s*by\s*\(\s*cluster\s*\)`)
	reSumByCluster    = regexp.MustCompile(`(sum|count|avg|min|max)\s+by\s*\(\s*cluster\s*\)`)
	// on(cluster, ...) or on(cluster)
	reOnClusterPrefix = regexp.MustCompile(`on\s*\(\s*cluster\s*,\s*`)
	reOnClusterOnly   = regexp.MustCompile(`on\s*\(\s*cluster\s*\)\s*group_left\s*\(\s*\)`)
	// name=~"$cluster" (ACM managed cluster labels)
	reNameClusterWithTrailingComma = regexp.MustCompile(`name=~"\$cluster"\s*,\s*`)
	reNameClusterWithLeadingComma  = regexp.MustCompile(`,\s*name=~"\$cluster"`)
	reNameClusterAlone             = regexp.MustCompile(`name=~"\$cluster"`)
	// by(..., cluster, ...) in the middle
	reByMiddleCluster  = regexp.MustCompile(`(by\s*\([^)]*),\s*cluster\s*,([^)]*\))`)
	reByMiddleCluster2 = regexp.MustCompile(`(by\s*\([^)]*),\s*cluster(\s*\))`)
)

func stripClusterFromQuery(query string) string {
	q := query

	// Remove cluster=~"$cluster" matchers (order matters: try with trailing comma first,
	// then leading comma, then standalone)
	q = reClusterMatcherWithTrailingComma.ReplaceAllString(q, "")
	q = reClusterMatcherWithLeadingComma.ReplaceAllString(q, "")
	q = reClusterMatcherAlone.ReplaceAllString(q, "")

	// Clean up malformed label matchers
	q = reDoubleComma.ReplaceAllString(q, ",")
	q = reLeadingComma.ReplaceAllString(q, "{")
	q = reTrailingComma.ReplaceAllString(q, "}")

	// Remove name=~"$cluster" (same order: trailing comma, leading comma, standalone)
	q = reNameClusterWithTrailingComma.ReplaceAllString(q, "")
	q = reNameClusterWithLeadingComma.ReplaceAllString(q, "")
	q = reNameClusterAlone.ReplaceAllString(q, "")
	q = reDoubleComma.ReplaceAllString(q, ",")
	q = reLeadingComma.ReplaceAllString(q, "{")
	q = reTrailingComma.ReplaceAllString(q, "}")

	// Remove by(cluster) entirely from aggregations
	q = reSumByCluster.ReplaceAllString(q, "$1")
	q = reByClusterOnly.ReplaceAllString(q, ")")

	// Remove cluster from by(cluster, ...)
	q = reByClusterPrefix.ReplaceAllString(q, "by (")

	// Remove cluster from by(..., cluster, ...)
	q = reByMiddleCluster.ReplaceAllString(q, "${1}${2}")
	q = reByMiddleCluster2.ReplaceAllString(q, "${1}${2}")

	// Remove on(cluster) group_left()
	q = reOnClusterOnly.ReplaceAllString(q, "")

	// Remove cluster from on(cluster, ...)
	q = reOnClusterPrefix.ReplaceAllString(q, "on(")

	// Apply metric renames
	for old, new := range metricRenames {
		q = strings.ReplaceAll(q, old, new)
	}

	// Clean up empty label matchers {}
	q = strings.ReplaceAll(q, "{}", "")

	return q
}

func rewriteQueries(spec map[string]any) {
	panels, ok := spec["panels"].(map[string]any)
	if !ok {
		return
	}
	for _, p := range panels {
		panel, ok := p.(map[string]any)
		if !ok {
			continue
		}
		panelSpec, _ := panel["spec"].(map[string]any)
		if panelSpec == nil {
			continue
		}
		queries, _ := panelSpec["queries"].([]any)
		for _, q := range queries {
			qm, ok := q.(map[string]any)
			if !ok {
				continue
			}
			qSpec, _ := qm["spec"].(map[string]any)
			if qSpec == nil {
				continue
			}
			plugin, _ := qSpec["plugin"].(map[string]any)
			if plugin == nil {
				continue
			}
			pSpec, _ := plugin["spec"].(map[string]any)
			if pSpec == nil {
				continue
			}
			if query, ok := pSpec["query"].(string); ok {
				pSpec["query"] = stripClusterFromQuery(query)
			}
		}
	}
}

func rewriteVariableValues(spec map[string]any) {
	vars, ok := spec["variables"].([]any)
	if !ok {
		return
	}
	for _, v := range vars {
		vm, ok := v.(map[string]any)
		if !ok {
			continue
		}
		varSpec, _ := vm["spec"].(map[string]any)
		if varSpec == nil {
			continue
		}
		plugin, _ := varSpec["plugin"].(map[string]any)
		if plugin == nil {
			continue
		}
		pSpec, _ := plugin["spec"].(map[string]any)
		if pSpec == nil {
			continue
		}

		// Rewrite matchers in label-values variables
		if matchers, ok := pSpec["matchers"].([]any); ok {
			for i, m := range matchers {
				if ms, ok := m.(string); ok {
					matchers[i] = stripClusterFromQuery(ms)
				}
			}
		}

		// Rewrite PromQL in promql variables (field can be "query" or "expr")
		if query, ok := pSpec["query"].(string); ok {
			pSpec["query"] = stripClusterFromQuery(query)
		}
		if expr, ok := pSpec["expr"].(string); ok {
			pSpec["expr"] = stripClusterFromQuery(expr)
		}

		// Rewrite static list values (VMStatusVariableJoinExpr)
		if values, ok := pSpec["values"].([]any); ok {
			for i, val := range values {
				if vm, ok := val.(map[string]any); ok {
					if v, ok := vm["value"].(string); ok {
						vm["value"] = stripClusterFromQuery(v)
					}
					values[i] = vm
				}
			}
		}

		// Rewrite default value
		if dv, ok := varSpec["defaultValue"].(map[string]any); ok {
			if sv, ok := dv["singleValue"].(string); ok {
				dv["singleValue"] = stripClusterFromQuery(sv)
			}
		}
	}
}

func removeClusterTableColumns(spec map[string]any) {
	panels, ok := spec["panels"].(map[string]any)
	if !ok {
		return
	}
	for _, p := range panels {
		panel, ok := p.(map[string]any)
		if !ok {
			continue
		}
		panelSpec, _ := panel["spec"].(map[string]any)
		if panelSpec == nil {
			continue
		}
		plugin, _ := panelSpec["plugin"].(map[string]any)
		if plugin == nil {
			continue
		}
		if kind, _ := plugin["kind"].(string); kind != "Table" {
			continue
		}
		pSpec, _ := plugin["spec"].(map[string]any)
		if pSpec == nil {
			continue
		}
		cols, ok := pSpec["columnSettings"].([]any)
		if !ok {
			continue
		}
		filtered := make([]any, 0, len(cols))
		for _, c := range cols {
			cm, ok := c.(map[string]any)
			if !ok {
				filtered = append(filtered, c)
				continue
			}
			if name, _ := cm["name"].(string); name == "cluster" {
				continue
			}
			filtered = append(filtered, c)
		}
		pSpec["columnSettings"] = filtered
	}
}

func removeClusterFromTransforms(spec map[string]any) {
	panels, ok := spec["panels"].(map[string]any)
	if !ok {
		return
	}
	for _, p := range panels {
		panel, ok := p.(map[string]any)
		if !ok {
			continue
		}
		panelSpec, _ := panel["spec"].(map[string]any)
		if panelSpec == nil {
			continue
		}
		plugin, _ := panelSpec["plugin"].(map[string]any)
		if plugin == nil {
			continue
		}
		pSpec, _ := plugin["spec"].(map[string]any)
		if pSpec == nil {
			continue
		}
		transforms, ok := pSpec["transforms"].([]any)
		if !ok {
			continue
		}
		for _, t := range transforms {
			tm, ok := t.(map[string]any)
			if !ok {
				continue
			}
			tSpec, _ := tm["spec"].(map[string]any)
			if tSpec == nil {
				continue
			}
			cols, ok := tSpec["columns"].([]any)
			if !ok {
				continue
			}
			filtered := make([]any, 0, len(cols))
			for _, c := range cols {
				if s, ok := c.(string); ok && s == "cluster" {
					continue
				}
				filtered = append(filtered, c)
			}
			tSpec["columns"] = filtered
		}
	}
}

func rewriteDashboardLinks(spec map[string]any, cfg Config) {
	panels, ok := spec["panels"].(map[string]any)
	if !ok {
		return
	}
	for _, p := range panels {
		panel, ok := p.(map[string]any)
		if !ok {
			continue
		}
		panelSpec, _ := panel["spec"].(map[string]any)
		if panelSpec == nil {
			continue
		}

		// Rewrite panel-level links (stat panels)
		if links, ok := panelSpec["links"].([]any); ok {
			filtered := rewriteLinks(links, cfg)
			if len(filtered) > 0 {
				panelSpec["links"] = filtered
			} else {
				delete(panelSpec, "links")
			}
		}

		// Rewrite dataLinks in table column settings
		plugin, _ := panelSpec["plugin"].(map[string]any)
		if plugin == nil {
			continue
		}
		pSpec, _ := plugin["spec"].(map[string]any)
		if pSpec == nil {
			continue
		}
		cols, ok := pSpec["columnSettings"].([]any)
		if !ok {
			continue
		}
		for _, c := range cols {
			cm, ok := c.(map[string]any)
			if !ok {
				continue
			}
			dl, ok := cm["dataLink"].(map[string]any)
			if !ok {
				continue
			}
			if url, ok := dl["url"].(string); ok {
				newURL := rewriteLinkURL(url, cfg)
				if newURL == "" {
					delete(cm, "dataLink")
				} else {
					dl["url"] = newURL
				}
			}
		}
	}
}

func rewriteLinks(links []any, cfg Config) []any {
	filtered := make([]any, 0, len(links))
	for _, l := range links {
		lm, ok := l.(map[string]any)
		if !ok {
			continue
		}
		url, _ := lm["url"].(string)
		newURL := rewriteLinkURL(url, cfg)
		if newURL == "" {
			continue
		}
		lm["url"] = newURL
		filtered = append(filtered, lm)
	}
	return filtered
}

func rewriteLinkURL(url string, cfg Config) string {
	// Rewrite acm- prefix to configured prefix before filtering
	if cfg.DashboardPrefix != "" {
		url = strings.ReplaceAll(url, "dashboard=acm-", "dashboard="+cfg.DashboardPrefix+"-")
	}

	// Remove links to dashboards not in the included set
	if len(cfg.IncludedDashboards) > 0 {
		if !linkTargetsIncludedDashboard(url, cfg.IncludedDashboards) {
			return ""
		}
	}

	// Remove var-cluster parameter
	reVarCluster := regexp.MustCompile(`[&?]var-cluster=[^&]*`)
	url = reVarCluster.ReplaceAllString(url, "")
	url = strings.ReplaceAll(url, "?&", "?")
	url = strings.TrimSuffix(url, "?")
	url = strings.TrimSuffix(url, "&")

	return url
}

func linkTargetsIncludedDashboard(rawURL string, included []string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	values := u.Query()["dashboard"]
	if len(values) == 0 {
		return true
	}
	for _, v := range values {
		for _, name := range included {
			if v == name {
				return true
			}
		}
	}
	return false
}

func rewriteDisplayName(spec map[string]any, cfg Config) {
	if cfg.DisplayNamePrefix == "" {
		return
	}
	display, ok := spec["display"].(map[string]any)
	if !ok {
		return
	}
	name, ok := display["name"].(string)
	if !ok {
		return
	}
	display["name"] = strings.Replace(name, "Virtualization /", cfg.DisplayNamePrefix+" /", 1)
}

func removeDatasourceRefs(spec map[string]any) {
	panels, ok := spec["panels"].(map[string]any)
	if !ok {
		return
	}
	for _, p := range panels {
		panel, ok := p.(map[string]any)
		if !ok {
			continue
		}
		panelSpec, _ := panel["spec"].(map[string]any)
		if panelSpec == nil {
			continue
		}
		queries, _ := panelSpec["queries"].([]any)
		for _, q := range queries {
			qm, ok := q.(map[string]any)
			if !ok {
				continue
			}
			qSpec, _ := qm["spec"].(map[string]any)
			if qSpec == nil {
				continue
			}
			plugin, _ := qSpec["plugin"].(map[string]any)
			if plugin == nil {
				continue
			}
			pSpec, _ := plugin["spec"].(map[string]any)
			if pSpec == nil {
				continue
			}
			delete(pSpec, "datasource")
		}
	}

	vars, ok := spec["variables"].([]any)
	if !ok {
		return
	}
	for _, v := range vars {
		vm, ok := v.(map[string]any)
		if !ok {
			continue
		}
		varSpec, _ := vm["spec"].(map[string]any)
		if varSpec == nil {
			continue
		}
		plugin, _ := varSpec["plugin"].(map[string]any)
		if plugin == nil {
			continue
		}
		pSpec, _ := plugin["spec"].(map[string]any)
		if pSpec == nil {
			continue
		}
		delete(pSpec, "datasource")
	}
}
