# LowVirtControllersCount

## Meaning

This alert fires when a low number of `virt-controller` devices is detected. More than one `virt-controller` pod should be available to ensure high availability. The current default number of replicas is 2.

A `virt-controller` handles monitoring the custom resource definitions (CRDs) of a virtual machine instance (VMI) and managing the associated pods. The device creates pods for VMIs and manages the life-cycle of the pods.

Therefore, `virt-controller` devices are critical for all cluster-wide virtualization functionality.

## Impact

The responsiveness of KubeVirt may become negatively affected. For example, certain requests may be missed.

In addition, if another `virt-launcher` instance terminates unexpectedly, KubeVirt may become completely unresponsive.

## Diagnosis
1. Set the NAMESPACE environment variable as follows:
    ```
    $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

2. Verify that any running `virt-controller` are available:
    ```
    $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-controller
    ```

3. Check whether any `virt-controller` pods have terminated unexpectedly or are in a `NotReady` state:
    ```
    $ kubectl -n $NAMESPACE logs virt-launcher-<unique-id>
    ```
    ```
    $ kubectl -n $NAMESPACE describe pod/virt-launcher-<unique-id>
    ```

## Mitigation

This alert can have a variety of causes, including:

- Not enough memory on the cluster
- Nodes are down
- The API server is overloaded. For example, the scheduler may be under heavy load and therefore not completely available.
- Networking issues

Verify whether any of these applies to your deployment, and fix it if possible. (TODO: Add how?)

If this does not fix the problem, open an issue (TBA where!) and attach to it the artifacts gathered in the Diagnosis section.