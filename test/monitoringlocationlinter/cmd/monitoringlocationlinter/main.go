package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/monitoring/test/monitoringlocationlinter"
)

func main() {
	singlechecker.Main(monitoringlocationlinter.NewAnalyzer())
}
