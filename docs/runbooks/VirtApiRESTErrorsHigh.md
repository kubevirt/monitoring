# VirtApiRESTErrorsHigh

## Meaning

Virt-api REST calls are failing in a high rate.
Alert will be fired when more than 5% of the REST calls failed in virt-api for the last hour.

## Impact

High rate of failed REST calls to virt-api could lead to slow response / execution of API calls.

Currently running VM workloads should not be affected. 

## Diagnosis

- Set the environment variable `NAMESPACE`
  ```
  export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
  ```

- Check to see how many running virt-api pods exist.
  ```
  kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-api
  ```

- For these pods, try to see the logs with `kubectl logs` and pod status via `kubectl describe`

## Mitigation
There can be several reasons for virt-api pods to be down, identify the root cause and fix it.

- Check the status of the virt-api deployment to find out more information. The following commands will provide the associated events and show if there are any issues with pulling an image, crashing pod, etc.
    - `kubectl -n $NAMESPACE get deploy virt-api -o yaml`
    - `kubectl -n $NAMESPACE describe deploy virt-api`
- Check if there are issues with the nodes. For example, if they are in a NotReady state.
    - `kubectl get nodes`
