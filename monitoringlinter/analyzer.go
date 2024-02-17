package monitoringlinter

import (
	"go/ast"
	"path"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	prometheusImportPath      = `"github.com/prometheus/client_golang/prometheus"`
	operatorMetricsImportPath = `"github.com/machadovilaca/operator-observability/pkg/operatormetrics"`
	operatorRulesImportPath   = `"github.com/machadovilaca/operator-observability/pkg/operatorrules"`
)

// NewAnalyzer returns an Analyzer.
func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "monitoringlinter",
		Doc: "Ensures that in Kubernetes operators projects, monitoring related practices are implemented " +
			"within pkg/monitoring directory, using operator-observability packages.",
		Run: run,
	}
}

// run is the main assertion function.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		prometheusLocalName, prometheusIsImported := getPackageLocalName(file, prometheusImportPath)
		operatorMetricsLocalName, OperatorMetricsIsImported := getPackageLocalName(file, operatorMetricsImportPath)
		operatorRulesLocalName, operatorRulesIsImported := getPackageLocalName(file, operatorRulesImportPath)

		if !prometheusIsImported && !OperatorMetricsIsImported && !operatorRulesIsImported {
			continue // no monitoring related packages are imported => nothing to do in this file;
		}

		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			selectorExpr, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			ident, ok := selectorExpr.X.(*ast.Ident)
			if !ok {
				return true
			}

			methodPackage := ident.Name
			methodName := selectorExpr.Sel.Name

			if prometheusIsImported && methodPackage == prometheusLocalName {
				checkPrometheusMethodCall(methodName, pass, node)
				return true

			}

			if !isMonitoringDir(pass.Fset.File(file.Pos()).Name()) {
				if OperatorMetricsIsImported && methodPackage == operatorMetricsLocalName {
					checkOperatorMetricsMethodCall(methodName, pass, node)
					return true
				}

				if operatorRulesIsImported && methodPackage == operatorRulesLocalName {
					checkOperatorRulesMethodCall(methodName, pass, node)
					return true
				}
			}

			return true
		})
	}

	return nil, nil
}

// checkPrometheusMethodCall checks if Prometheus method call should be reported.
func checkPrometheusMethodCall(methodName string, pass *analysis.Pass, node ast.Node) {
	if methodName == "Register" || methodName == "MustRegister" {
		pass.Reportf(node.Pos(), "monitoring-linter: metrics should be registered only within pkg/monitoring directory, "+
			"using operator-observability packages.")
	}
}

// checkOperatorMetricsMethodCall checks if operatormetrics method call should be reported.
func checkOperatorMetricsMethodCall(methodName string, pass *analysis.Pass, node ast.Node) {
	if methodName == "RegisterMetrics" {
		pass.Reportf(node.Pos(), "monitoring-linter: metrics should be registered only within pkg/monitoring directory.")
	}
}

// checkOperatorRulesMethodCall checks if operatorrules method call should be reported.
func checkOperatorRulesMethodCall(methodName string, pass *analysis.Pass, node ast.Node) {
	if methodName == "RegisterAlerts" || methodName == "RegisterRecordingRules" {
		pass.Reportf(node.Pos(), "monitoring-linter: alerts and recording rules should be registered only within pkg/monitoring directory.")
	}
}

// getPackageLocalName returns the name a package was imported with in the file.
// e.g. import prom "github.com/prometheus/client_golang/prometheus"
func getPackageLocalName(file *ast.File, importPath string) (string, bool) {
	for _, imp := range file.Imports {
		if imp.Path.Value == importPath {
			if name := imp.Name.String(); name != "<nil>" {
				return name, true
			}
			pathWithoutQuotes := strings.Trim(importPath, `"`)
			return path.Base(pathWithoutQuotes), true
		}
	}
	return "", false
}

func isMonitoringDir(filePath string) bool {
	return strings.Contains(filePath, "pkg/monitoring")
}
