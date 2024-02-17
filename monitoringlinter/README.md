# monitoring-linter

Monitoring Linter is a Golang linter designed to ensure that in Kubernetes operator projects, 
monitoring-related practices are implemented within the `pkg/monitoring` directory using [operator-observability](https://github.com/machadovilaca/operator-observability/tree/main) methods.
It verifies that all metrics, alerts and recording rules registrations are centralized in this directory.
The use of [Prometheus registration methods](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#Registerer) is restricted across the entire project.

## Installation

```shell
go install github.com/kubevirt/monitoring/monitoringlinter/cmd/monitoringlinter@latest
```

## Usage
Once installed, you can run the Monitoring Linter against your Kubernetes operator project by using the following command:

```shell
monitoringlinter  ./...
```


## Linter Rules

- The following Prometheus registration methods calls are restricted across all the project:
    ```go
      prometheus.Register()
      prometheus.MustRegister()
    ```

- The following operator-observability methods calls are allowed only within `pkg/monitoring` directory:
    ```go
      operatormetrics.RegisterMetrics()
      operatorrules.RegisterAlerts()
      operatorrules.RegisterRecordingRules()
    ```

- Examples for calls that are allowed across all the project, as they are not registering metrics, alerts or recording rules:
    ```go
      prometheus.NewHistogram()
      operatormetrics.ListMetrics()
    ```