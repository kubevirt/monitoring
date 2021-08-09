# LowReadyVirtControllersCount

## Meaning

From a high-level perspective the virt-controller has all the cluster wide virtualization functionality.

This controller is responsible for monitoring the VMI (CRDs) and managing the associated pods. The controller will make sure to create and manage the life-cycle of the pods associated to the VMI objects.

A VMI object will always be associated to a pod during it's life-time, however, due to i.e. migration of a VMI the pod instance might change over time.

This alert is triggered when some virt controllers are running but not ready in the past 5 minutes.

## Impact
In these circumstances the virt-controller will become a single point of failure.

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
- Check the status of the virt-controller deployment to find out more information.   
  The following commands will provide the associated events and show if there are any issues with pulling an image, crashing pod, etc. 
    ```
    kubectl -n $NAMESPACE get deploy virt-controller -o yaml
    ```
    ```
    kubectl -n $NAMESPACE describe deploy virt-controller
    ```
- Check if there are issues with the nodes. For example, if they are in a NotReady state.
    ```
    kubectl get nodes
    ```