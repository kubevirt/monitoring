# CI Metrics Tests
## prom-metrics-linter test
### Test
prom-metrics-linter is a common image for testing that the metrics names follow the naming conventions in all the Kubevirt sub-operators projects.
This linter is based on [promlint](https://github.com/prometheus/client_golang/blob/v1.16.0/prometheus/testutil/promlint/promlint.go), with additional tests customized for Kubevirt.
### Test tools
The test is using GitHub Actions to create a docker image with the linter tests.
The Dockerfile is here: [Dockerfile](https://github.com/kubevirt/monitoring/blob/main/test/metrics/prom-metrics-linter/Dockerfile).\
The action is triggered by version release, here is the [GitHub action](https://github.com/kubevirt/monitoring/blob/main/.github/workflows/prom-metrics-linter.yaml), that builds the image and pushes it to quay.io, with a new tag.

The assumption is that there won't be frequent changes in the prom-metrics-linter image.

### Run the test
The linter image requires 3 parameters, the operator name, the sub-operator name, and a json file with the sub-operator metrics.\
For collecting the metrics in the right format, please use this common help functions: [pkg/metrics/parser](https://github.com/kubevirt/monitoring/tree/main/pkg/metrics/parser).

Then run the linter image:\
`<container runtime> run -i "quay.io/kubevirt/prom-metrics-linter:<linter_image_tag>" --metric-families=<metrics_json_format>" --operator-name=kubevirt --sub-operator-name=<sub-operator-name>`.\
Note: If there is no sub operator, set the sub-operator-name the same as the operator-name.

Example:\
**Containerized Data Importer** sub-operator: [[CI] Add metrics name linter](https://github.com/kubevirt/containerized-data-importer/pull/2774).

