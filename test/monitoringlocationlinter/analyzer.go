package monitoringlocationlinter

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	prometheusImportPath = `"github.com/prometheus/client_golang/prometheus"`
)

// run is the main assertion function
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if !isInMonitoringDir(pass.Fset.File(file.Pos()).Name()) {
			prometheusLocalName := getPrometheusLocalName(file)
			if prometheusLocalName == "" { // no prometheus import => no use in prometheus in this file; nothing to do here
				continue
			}

			ast.Inspect(file, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if ok && isRegisterMetricsCall(call, prometheusLocalName) {
					pass.Reportf(node.Pos(), "monitoring-location-linter: Prometheus metrics should be registered only under pkg/monitoring.")
					return true
				}

				return true
			})
		}
	}

	return nil, nil
}

// NewAnalyzer returns an Analyzer
func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "monitoringlocationlinter",
		Doc:  "Ensures that in Kubernetes operators projects, monitoring-related practices are implemented within pkg/monitoring directory.",
		Run:  run,
	}
}

// isInMonitoringDir returns true if the file is located under pkg/monitoring directory.
func isInMonitoringDir(filePath string) bool {
	return strings.Contains(filePath, "pkg/monitoring")
}

// getPrometheusLocalName returns the name Prometheus was imported with in the file.
// e.g. import prom "github.com/prometheus/client_golang/prometheus"
// In case Prometheus is not imported, it returns an empty string.
func getPrometheusLocalName(file *ast.File) string {
	for _, imp := range file.Imports {
		if imp.Path.Value == prometheusImportPath {
			if name := imp.Name.String(); name != "<nil>" {
				return name
			}
			return "prometheus"
		}
	}
	return "" // No Prometheus import found; this file does not use Prometheus
}

// isRegisterMetricsCall returns true if a call is a register metrics call.
func isRegisterMetricsCall(call *ast.CallExpr, prometheusLocalName string) bool {
	selectorExpr, ok := call.Fun.(*ast.SelectorExpr)
	if ok && (selectorExpr.Sel.Name == "Register" || selectorExpr.Sel.Name == "MustRegister") {
		if ident, ok := selectorExpr.X.(*ast.Ident); ok {
			if ident.Name == prometheusLocalName {
				return true
			}
		}
	}
	return false
}
