# NoReadyVirtController

## Meaning

This alert fires when no available `virt-controller` devices have been detected for 5 minutes.

A `virt-controller` handles monitoring the custom resource definitions (CRDs) of a virtual machine instance (VMI) and managing the associated pods. The device creates pods for VMIs and manages the life-cycle of the pods.

Therefore, `virt-controller` devices are critical for all cluster-wide virtualization functionality.

## Impact
Any actions related to VM life-cycle management fail. This notably includes:

- Launching a new VM instance

- Shutting down an existing VM instance


## Diagnosis

- Set the `NAMESPACE` environment variable.
    ```
    $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

- Run the following command to verify a `virt-controller` device is available.
    ```
    $ kubectl get deployment -n $NAMESPACE virt-controller -o jsonpath='{.status.readyReplicas}'
    ```

## Mitigation
1. Check the status of the `virt-controller` deployment to find out more information. The following commands provide the associated events and show if any problems occurred, such as failures to pull images, or crashing pods.
    ```
    $ kubectl -n $NAMESPACE get deploy virt-controller -o yaml
    ```
    ```
    $ kubectl -n $NAMESPACE describe deploy virt-controller
    ```

2. Obtain the names of `virt-controller` pods.
    ```
    $ get pods -n $NAMESPACE |grep virt-controller
    ```

3. Inspect the logs for each `virt-controller`.
    ```
    $ kubectl logs -n $NAMESPACE <virt-controller-XYZ>
    ```

4. Check if any problems occurred with the nodes. For example, they might be in a `NotReady` state.
    ```
    $ kubectl get nodes
    ```