# VirtHandlerRESTErrorsBurst

## Meaning

More than 80% of REST calls failed in `virt-handler` in the last 5 minutes.

The `virt-handler` has likely fully lost connectivity to the API server. Running workloads on the affected node are not impacted, but status updates are not propagated and node-related actions, such as migrations, will fail.

## Diagnosis

This error is most frequently caused by one of the following problems:

- The API server is overloaded, which causes timeouts. To verify if this is the case, check the metrics of the API server, and view its response times and overall calls.

- The `virt-handler` pod cannot reach the API server. This is commonly caused by DNS issues on the node and networking connectivity issues.

Check whether `virt-handler` can connect to the API server.

1. List all `virt-handler` pods:
    ```
     $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

2. Obtain the name of the `virt-handler` pod:

    ```
    $ kubectl get pods -n $NAMESPACE -l=kubevirt.io=virt-handler
    ```

3. Review the logs for `virt-handler` and check its connectivity to the API server:

    ```
    $ kubectl logs -n  $NAMESPACE <virt-handler-pod-name>
    ```


## Mitigation
If the `virt-handler` cannot connect to the API server, delete the pod to force a restart:

```
$ kubectl delete -n <kubevirt-install-namespace> <virt-handler-pod-name>
```


