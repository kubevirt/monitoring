# VirtControllerRESTErrorsHigh

## Meaning

More than 5% of REST calls failed in `virt-controller` in the last 60 minutes. 

This is most likely because `virt-controller` has partially lost connection to the API server.


## Impact

Node-related actions, such as starting and migrating, and scheduling virtual machines, are delayed. Running workloads is not affected, but reporting their current status may be delayed.

## Diagnosis

This error is most frequently caused by one of the following problems:

- The API server is overloaded, which causes timeouts. To verify if this is the case, check the metrics of the API server, and view its response times and overall calls.

- The `virt-controller` pod cannot reach the API server. This is commonly caused by DNS issues on the node and networking connectivity issues.

Check whether `virt-controller` can connect to the API server.

1. List all `virt-controller` pods:
    ```
     $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

1. Obtain the name of the `virt-controller` pod:

    ```
    $ kubectl get pods -n $NAMESPACE -l=kubevirt.io=virt-controller
    ```

1. Review the logs for `virt-controller` and check its connectivity to the API server:

    ```
    $ kubectl logs -n  $NAMESPACE <virt-controller-pod-name>
    ```


## Mitigation

If the `virt-controller` pod cannot connect to the API server, delete the pod to force a restart:

```
$ kubectl delete -n <kubevirt-install-namespace> <virt-controller-pod-name>
```