# VirtOperatorRESTErrorsHigh 

## Meaning

More than 5% of REST calls in `virt-operator` pods failed in the last 60 minutes. 

This is most likely because `virt-operator` has partially lost connection to the API server.

## Impact

Cluster-level actions, such as upgrading and controller reconciliation, may be delayed. 

However, customer workloads, such as virtual machines (VMs) and VMIs, are not likely to be affected.

## Diagnosis

This error is most frequently caused by one of the following problems:

- The API server is overloaded, which causes timeouts. To verify if this is the case, check the metrics of the API server, and view its response times and overall calls.

- The `virt-operator` pod cannot reach the API server. This is commonly caused by DNS issues on the node and networking connectivity issues.

Check whether `virt-operator` can connect to the API server.

1. List all `virt-operator` pods:
    ```
     $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

1. Obtain the name of the `virt-operator` pod:

    ```
    $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-operator
    ```

1. Review the logs for `virt-operator` and check its connectivity to the API server:

    ```
    $ kubectl -n $NAMESPACE logs <virt-operator-pod-name>
    ```

1. Obtain the description of the `virt-operator` pod:
    
    ```
    $ kubectl -n $NAMESPACE describe pod <virt-operator-pod-name>
    ```

## Mitigation

If `virt-operator` cannot connect to the API server, delete the pod to force a restart:

```
$ kubectl delete -n <kubevirt-install-namespace> <virt-operator-pod-name>
```
