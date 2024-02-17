# monitoring-location-linter

Monitoring Location Linter is a Golang linter designed specifically for Kubernetes operators. 
Its purpose is to ensure monitoring-related practices are implemented only under pkg/monitoring directory. 

## Installation

```shell
go install github.com/kubevirt/monitoring/test/monitoringlocationlinter/cmd/monitoringlocationlinter@latest
```

## Usage
Once installed, you can run the Monitoring Location Linter against your Kubernetes operator project by using the following command:

```shell
monitoringlocationlinter  ./...
```


## Linter Rules
Currently, the linter verifies that all metrics registration (using [Prometheus metrics registration functions](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#Registerer)) 
are implemented within pkg/monitoring directory.

Example:
```go
package main

import "github.com/prometheus/client_golang/prometheus"

func main() {
	prometheus.Register(prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gauge_name",
			Help: "gauge_help",
		}))
}
```

