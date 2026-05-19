# Dashboard Review Guidelines

Applies to: `dashboards/**`

- Verify that PromQL queries reference valid metrics or recording rules
  from `docs/metrics.md`, the kubevirt operator source repositories, or
  standard Kubernetes/Prometheus ecosystem metrics.
- Check that panel titles and descriptions are clear and accurate.
- Ensure dashboard variables (e.g., `$namespace`) are used consistently
  and match Kubernetes label conventions.
- Review Grafana JSON for proper formatting and no hardcoded values
  (namespaces, cluster names, etc.).
