# NoReadyVirtOperator 

## Meaning

This alert fires when no `virt-operator` pod in a `Ready` state has been detected for 10 minutes. The virt-operator deployment has a default replica of two `virt-operator` pods.

The `virt-operator` is the first operator to start in a cluster. Its primary responsibilities include the following: 

- Installing, live updating, and live upgrading a cluster

- Monitoring the life cycle of top-level controllers, such as `virt-controller`, `virt-handler`, `virt-launcher`, and managing their reconciliation

- Certain cluster-wide tasks, such as certificate rotation and infrastructure management

## Impact 

This alert indicates a failure at the level of the cluster. 

As a result, critical cluster-wide management functionalities, such as certification rotation, upgrade, and reconciliation of controllers, are currently not available.

Note, however, that `virt-operator` is not directly responsible for virtual machines in the cluster. Therefore, its temporary unavailability does not significantly affect custom workloads.

## Diagnosis


1. Obtain the namespace data of the `virt-operator` deployment
    ```
    $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

2. Obtain the name of the `virt-operator` deployment.
    ```
    $ kubectl -n $NAMESPACE get deploy virt-operator -o yaml
    ```

3. Generate the description of the `virt-operator` deployment.
    ```
    $ kubectl -n $NAMESPACE describe deploy virt-operator
    ```

4. Check for node issues, such as a `NotReady` state:
    ```
    $ kubectl get nodes
    ```

## Mitigation

This alert can have a number of causes, including:

- TBA possible causes!

Verify whether any of these applies to your deployment, and fix it if possible (TBA how?).

If this does not fix the problem, open an issue (TBA where!) and attach to it the artifacts gathered in the Diagnosis section.
