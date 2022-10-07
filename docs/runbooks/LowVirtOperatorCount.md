# LowVirtOperatorCount 

## Meaning

This alert fires when only one `virt-operator` pod in a `Ready` state has been running for the last 60 minutes. 

The `virt-operator` is the first operator to start in a cluster. Its primary responsibilities include the following: 

- Installing, live updating, and live upgrading a cluster

- Monitoring the life cycle of top-level controllers, such as `virt-controller`, `virt-handler`, `virt-launcher`, and managing their reconciliation

- Certain cluster-wide tasks, such as certificate rotation and infrastructure management

## Impact

This alert indicates that `virt-operator` cannot provide high availability (HA) for the deployment.

For HA to work reliably, two or more `virt-operator` pods in a `Ready` state must be available. The `virt-operator` deployment has a default replica of two `virt-operator` pods.

Note, however, that `virt-operator` is not directly responsible for virtual machines in the cluster. Therefore, its decreased availability does not significantly affect custom workloads.


## Diagnosis

1. Obtain the namespace data of the `virt-operator` deployment:
    ```
    $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

1. Check the states of `virt-operator` pods:

    ```
    $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-operator
    ```

1. Review the logs of the affected `virt-operator` pods:
    ```
    $ kubectl -n $NAMESPACE logs <pod-name>
    ```

1.  Generate the descriptions of the affected `virt-operator` pods:

    ```
    $ kubectl -n $NAMESPACE describe pod <pod-name>
    ```

## Mitigation

This alert can have a number of causes, including:

- TBA possible causes!

Verify whether any of these applies to your deployment, and fix it if possible (TBA how?).

If this does not fix the problem, open an issue (TBA where!) and attach to it the artifacts gathered in the Diagnosis section.