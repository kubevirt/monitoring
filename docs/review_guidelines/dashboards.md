# Dashboard Review Guidelines

Applies to: `dashboards/**`

## Grafana dashboards (`dashboards/grafana/`)

- Verify that PromQL queries reference valid metrics or recording rules
  from `docs/metrics.md`, the kubevirt operator source repositories, or
  standard Kubernetes/Prometheus ecosystem metrics.
- Check that panel titles and descriptions are clear and accurate.
- Ensure dashboard variables (e.g., `$namespace`) are used consistently
  and match Kubernetes label conventions.
- Review Grafana JSON for proper formatting and no hardcoded values
  (namespaces, cluster names, etc.).

## Perses dashboards (`dashboards/perses/`)

Perses dashboard YAMLs are **generated artifacts** produced by
`cmd/generate-perses-dashboards/`. Do not review PromQL or panel
logic in the YAML files — review the Go builder code instead.

When reviewing changes to generated Perses dashboards:

- Confirm `make verify-perses-dashboards` passes (CI checks this).
- Review the Go source in `cmd/generate-perses-dashboards/` and
  `pkg/dashboards/` for correctness, not the YAML output.
- For multicluster builders (imported from MCOA), PromQL logic
  is owned upstream — only review the transform layer in
  `pkg/dashboards/transform/`.
- For local builders (`multicluster: false`), review the builder
  function for correct PromQL, variables, and panel structure.
- Verify that new dashboards are added to the `builders` slice
  in `cmd/generate-perses-dashboards/main.go` with the correct
  `multicluster` flag.
