# VirtApiRESTErrorsBurst
<!--Edited by jherrman, 17.10.2022-->

## Meaning

More than 80% of REST calls have failed in `virt-api` in the last 5 minutes.

## Impact

Very high rate of failed REST calls to `virt-api` may lead to slow response and execution of API calls, and potentially to API calls being completely dismissed.

However, currently running virtual machine workloads are not likely to be affected. 

## Diagnosis

1. Set the `NAMESPACE` environment variable as follows:
   ```
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

1. Check how many `virt-api` pods are running on your deployment:
   ```
   $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-api
   ```

1. For each of the `virt-api` pods, display its logs:
   ```
   $ kubectl logs -n  $NAMESPACE <virt-api-pod-name>
   ```

1. Display the states of the `virt-api` pods:
   ```
   $ kubectl describe -n $NAMESPACE <virt-api-pod-name>
   ```

1. Check if any problems occurred with the nodes. For example, they might be in a `NotReady` state:
   ```
   $ kubectl get nodes
   ```

1. Check the status of the `virt-api` deployment to find out more information. The following commands provide the associated events and show if any problems occurred, such as crashing pods or failures to pull images:
   ```
   $ kubectl -n $NAMESPACE get deploy virt-api -o yaml
   ```
   ```
   $ kubectl -n $NAMESPACE describe deploy virt-api
   ```

## Mitigation

If any of the mentioned commands outputs an error, attempt to fix the cause.

The possible fixes may include the following:

- TODO: give some more specific info here!

