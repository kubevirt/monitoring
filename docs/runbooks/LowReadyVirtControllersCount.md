# LowReadyVirtControllersCount

## Meaning

This alert fires when one or more `virt-controller` pods are running, but none of these pods have been in a `Ready` state for the last 5 minutes. 

A `virt-controller` device monitors the custom resource definitions (CRDs) of a virtual machine instance (VMI) and manages the associated pods. The device create pods for VMIs and manages the lifecycle of the pods. The device is critical for cluster-wide virtualization functionality.

## Impact

This alert indicates that a cluster-level failure might occur, which would cause actions related to VM life-cycle management to fail. This notably includes launching a new VMI or shutting down an existing VMI.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

  ```
  $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
  ```

2. Verify a `virt-controller` device is available:
  ```
  $ kubectl get deployment -n $NAMESPACE virt-controller -o jsonpath='{.status.readyReplicas}'
  ```

## Mitigation
1. Check the status of the `virt-controller` deployment to find out more information. The following commands provide the associated events and show if any problems occurred, such as crashing pods or failures to pull images:
    ```
    $ kubectl -n $NAMESPACE get deploy virt-controller -o yaml
    ```
    ```
    $ kubectl -n $NAMESPACE describe deploy virt-controller
    ```
2. Check if any problems occurred with the nodes. For example, they might be in a `NotReady` state:
    ```
    $ kubectl get nodes
    ```

<!--CNV: If you cannot resolve the issue, log in to the [Customer Portal](https://access.redhat.com) and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->

<!--KVstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--KVend-->