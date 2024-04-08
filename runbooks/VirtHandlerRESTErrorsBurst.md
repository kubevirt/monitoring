# VirtHandlerRESTErrorsBurst

## Meaning

For the last 10 minutes or longer, over 80% of the REST calls made to `virt-handler` pods have failed.

This alert usually indicates that the `virt-handler` pods cannot connect to the API server.

This error is frequently caused by one of the following problems:

- The API server is overloaded, which causes timeouts. To verify if this is the case, check the metrics of the API server, and view its response times and overall calls.

- The `virt-handler` pod cannot reach the API server. This is commonly caused by DNS issues on the node and networking connectivity issues.

## Impact

Status updates are not propagated and node-related actions, such as migrations, fail. However, running workloads on the affected node are not impacted.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-handler` pod:

   ```bash
   $ kubectl get pods -n $NAMESPACE -l=kubevirt.io=virt-handler
   ```

3. Check the `virt-handler` logs for error messages when connecting to the API server:

   ```bash
   $ kubectl logs -n  $NAMESPACE <virt-handler>
   ```

## Mitigation

If the `virt-handler` cannot connect to the API server, delete the pod to force a restart:

```bash
$ kubectl delete -n $NAMESPACE <virt-handler>
```

<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->

