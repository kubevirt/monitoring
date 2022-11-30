# KubeVirtComponentExceedsRequestedCPU
<!-- Edited by apinnick, Nov 2022-->

## Meaning

This alert fires when a component's CPU usage exceeds the requested limit.

## Impact

Usage of CPU resources is not optimal and the node might be overloaded.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the component's CPU request limit:

   ```bash
   $ kubectl -n $NAMESPACE get deployment <component> -o yaml | grep requests: -A 2
   ```

3. Check the actual CPU usage by using a PromQL query:

   ```  
node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace="$NAMESPACE",container="<component>"}
   ```

See the [Prometheus documentation](https://prometheus.io/docs/prometheus/latest/querying/basics/) for more information.

## Mitigation

<!--DS: Update the CPU request limit in the `HCO` custom resource.-->
<!--USstart-->
Update the CPU resource request in the `KubeVirt` custom resource as in the following example:

```yaml
spec:
  customizeComponents:
    patches:
    - type:
      resourceName: <component>
      resourceType: < Deployment|DaemonSet >
      type: strategic
      patch: '{"spec":{"template":{"spec":{"containers":[{"name":"<component>","resources":{"requests":{"cpu":" <cpu_request> "}}}]}}}}'
```
<!--USend-->