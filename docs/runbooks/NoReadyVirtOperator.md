# NoReadyVirtOperator 

## Meaning

This alert fires when no `virt-operator` pod in the `Ready` state has been detected for 10 minutes. The virt-operator deployment has a default replica of 2 `virt-operator` pods.

The `virt-operator` is the first Kubernetes operator coming alive in a KubeVirt cluster. Its primary responsibilities include: 

- Installation, live updating, and live upgrading a KubeVirt cluster

- Monitoring the life-cycle of top-level controllers, such as `virt-controller`, `virt-handler`, `virt-launcher`, etc. and managing their reconciliation

In addition, `virt-operator` is responsible for cluster-wide tasks, such as certificate rotation and certain infrastructure management. 

## Impact 

This alert indicates a failure at the level of the KubeVirt cluster. 

As a result, critical cluster-wide management functionalities, such as certification rotation, KubeVirt upgrade, and reconciliation of KubeVirt controllers, are currently not available.

Note, however, that `virt-operator` is not directly responsible for virtual machines in the cluster. Therefore, its temporary unavailability does not significantly affect custom workloads.

## Diagnosis

1. Check the status of the `virt-operator` deployment to find out more information. The following commands provide the associated events and show if any specific problems occurred.
    ```
    $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    $ kubectl -n $NAMESPACE get deploy virt-operator -o yaml
    $ kubectl -n $NAMESPACE describe deploy virt-operator
    ```

2. Check if any problems occurred with the nodes. For example, they might be in a `NotReady` state:
    ```
    $ kubectl get nodes
    ```

## Mitigation

This alert can have a number of causes, including:

- TBA possible causes!

Verify whether any of these applies to your deployment, and fix it if possible (TBA how?).

If this does not fix the problem, open an issue (TBA where!) and attach to it the artifacts gathered in the Diagnosis section.
