<!-- Edited by Jiří Herrmann, 7 Nov 2022 -->

# LowReadyVirtOperatorsCount 

## Meaning

This alert fires when one or more `virt-operator` pods are running, but none of these pods have been in a `Ready` state for the last 10 minutes. 

The `virt-operator` is the first operator to start in a cluster. Its primary responsibilities include the following: 

- Installing, live updating, and live upgrading a cluster

- Monitoring the life cycle of top-level controllers, such as `virt-controller`, `virt-handler`, `virt-launcher`, and managing their reconciliation

- Certain cluster-wide tasks, such as certificate rotation and infrastructure management


## Impact

This alert indicates that a cluster-level failure might occur, and that critical cluster-wide management functionalities, such as certification rotation, upgrade, and reconciliation of controllers, might become unavailable. Such a state would also trigger the `NoReadyVirtOperator` alert.

The `virt-operator` deployment has a default replica of two `virt-operator` pods.

Note, however, that `virt-operator` is not directly responsible for virtual machines in the cluster. Therefore, its temporary unavailability does not significantly affect custom workloads.

## Diagnosis

1. Set the `NAMESPACE` environment variable:
    ```
    $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

2. Obtain the name of the `virt-operator` deployment:
    ```
    $ kubectl -n $NAMESPACE get deploy virt-operator -o yaml
    ```

3. Obtain the details of the `virt-operator` deployment:
    ```
    $ kubectl -n $NAMESPACE describe deploy virt-operator
    ```

4. Check for node issues, such as a `NotReady` state:
    ```
    $ kubectl get nodes
    ```
    
## Mitigation

<!--CNV: If you cannot resolve the issue, log in to the [Customer Portal](https://access.redhat.com) and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->

<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->