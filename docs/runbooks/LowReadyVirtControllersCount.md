# LowReadyVirtControllersCount

## Meaning

This alert fires when one or more `virt-controller` pods are running, but none of these pods have been in a `Ready` state for the last 5 minutes. 

A `virt-controller` handles monitoring the custom resource definitions (CRDs) of a virtual machine instance (VMI) and managing the associated pods. The device creates pods for VMIs and manages the life-cycle of the pods.

Therefore, `virt-controller` devices are critical for all cluster-wide virtualization functionality. 

## Impact

This alert indicates that a cluster-level failure may occur, which would cause actions related to VM life-cycle management to fail. This notably includes:

- Launching a new VM instance

- Shutting down an existing VM instance

## Diagnosis

1. Set the `NAMESPACE` environment variable as follows:

  ```
  $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
  ```

1. Run the following command to verify a `virt-controller` device is available:
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
1. Check if any problems occurred with the nodes. For example, they might be in a `NotReady` state:
    ```
    $ kubectl get nodes
    ```