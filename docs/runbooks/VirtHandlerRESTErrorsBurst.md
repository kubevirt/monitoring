<!-- Edited by Jiří Herrmann, 7 Nov 2022 -->

# VirtHandlerRESTErrorsBurst

## Meaning

More than 80% of REST calls failed in `virt-handler` in the last 5 minutes. This alert usually indicates that the `virt-handler` pods cannot connect to the API server.

This error is most frequently caused by one of the following problems:

- The API server is overloaded, which causes timeouts. To verify if this is the case, check the metrics of the API server, and view its response times and overall calls.

- The `virt-handler` pod cannot reach the API server. This is commonly caused by DNS issues on the node and networking connectivity issues.

## Impact

Running workloads on the affected node are not impacted, but status updates are not propagated and node-related actions, such as migrations, fail.

## Diagnosis

Check whether `virt-handler` can connect to the API server.

1. Set the `NAMESPACE` environment variable:
    ```
     $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

2. Obtain the name of the `virt-handler` pod:

    ```
    $ kubectl get pods -n $NAMESPACE -l=kubevirt.io=virt-handler
    ```

3. Check the `virt-handler` logs for error messages when connecting to the API server:

    ```
    $ kubectl logs -n  $NAMESPACE <virt-handler-pod-name>
    ```


## Mitigation

If the `virt-handler` cannot connect to the API server, delete the pod to force a restart:

```
$ kubectl delete -n <install-namespace> <virt-handler-pod-name>
```
<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->

