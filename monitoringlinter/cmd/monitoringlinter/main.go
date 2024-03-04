package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/kubevirt/monitoring/monitoringlinter"
)

func main() {
	singlechecker.Main(monitoringlinter.NewAnalyzer())
}
