# VirtHandlerRESTErrorsHigh
<!-- Edited by Jiří Herrmann, 9 Nov 2022 -->

## Meaning

More than 5% of REST calls failed in `virt-handler` in the last 60 minutes. This alert usually indicates that the `virt-handler` pods have partially lost connection to the API server.

This error is frequently caused by one of the following problems:

- The API server is overloaded, which causes timeouts. To verify if this is the case, check the metrics of the API server, and view its response times and overall calls.

- The `virt-handler` pod cannot reach the API server. This is commonly caused by DNS issues on the node and networking connectivity issues.

## Impact

Node-related actions, such as starting and migrating workloads, are delayed on the node that `virt-handler` is running on. Running workloads are not affected, but reporting their current status might be delayed.

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
$ kubectl delete -n <install-namespace> <virt-handler>
```

<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->