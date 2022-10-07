# VirtHandlerRESTErrorsHigh

## Meaning

More than 5% of REST calls failed in `virt-handler` in the last 60 minutes. This is most likely because `virt-handler` has partially lost connection to the API server.

## Impact

Node-related actions, such as starting and migrating workloads, are delayed on the node that `virt-handler` is running on. Running workloads is not affected, but reporting their current status may be delayed.

## Diagnosis

This error is most frequently caused by one of the following problems:

- The apiserver is overloaded, which causes timeouts. To verify if this is the case, check the metrics of the API server, and view its response times and overall calls.

- The `virt-handler` pod cannot reach the API server. This is commonly caused by DNS issues on the node and networking connectivity issues.

Check whether `virt-handler` can connect to the API server:

1. List all the `virt-handler` pods:
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
