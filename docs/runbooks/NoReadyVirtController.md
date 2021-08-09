# NoReadyVirtController

## Meaning

From a high-level perspective the virt-controller is responsible for all the cluster wide virtualization functionality.

This controller is responsible for monitoring the VMI (CRDs) and managing the associated pods. The controller will make sure to create and manage the life-cycle of the pods associated to the VMI objects.

A VMI object will always be associated to a pod during it's life-time, however, due to i.e. migration of a VMI the pod instance might change over time.

This alert is triggered when no ready virt controllers were detected for the last 5 min.

## Impact
Complete failure in VM lifecycle management i.e. launching a new VM instance or shutting down an existing VM instance.


## Diagnosis

- Set the environment variable NAMESPACE

    ```
    export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

- Run
    ```
    kubectl get deployment -n $NAMESPACE virt-controller -o jsonpath='{.status.readyReplicas}'
    ```

## Mitigation
- Check the status of the virt-controller deployment to find out more information. The following commands will provide the associated events and show if there are any issues with pulling an image, crashing pod, etc. 
    ```
    kubectl -n $NAMESPACE get deploy virt-controller -o yaml
    ```
    ```
    kubectl -n $NAMESPACE describe deploy virt-controller
    ```

- Get virt-controller pods
    ```
    get pods -n $NAMESPACE |grep virt-controller
    ```

- Inspect the logs for each virt-controller
    ```
    kubectl logs -n $NAMESPACE <virt-controller-XYZ>
    ```

- Check if there are issues with the nodes. For example, if they are in a NotReady state.
    ```
    kubectl get nodes
    ```