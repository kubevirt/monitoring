# VirtControllerDown

## Meaning
No running virt-controller was detected for the last 5 min. Virt-controller deployment has a default replica of 2 pods.

## Impact
Complete failure in VM lifecycle management i.e. launching a new VM instance or shutting down an existing VM instance.

## Diagnosis

- Set the environment variable NAMESPACE

    ```
    export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

- Observe the status of the virt-controller deployment
    ```
    kubectl get deployment -n $NAMESPACE virt-controller -o yaml
    ```
- Observe the logs of the kubecontroller manager pod, to see why it cannot create the virt-controller pods.
    ```
    kubectl get logs <virt-controller-pod>
    ```

    *Note:* virt-controller pod should have name similar to virt-controller-7888c64d66-dzc9p.\
    There could be several pods that run virt-controller

## Mitigation
There can be several reasons. Like:

- Node resource exhaustion
- Not enough memory on the cluster
- Nodes are down
- API server is overloaded (e.g. Scheduler has a lot of work and therefore is not 100% available)
- Networking issues

Try to identify the root cause and fix it.

In other cases, please open an issue and attach the artifacts gathered in the Diagnosis section.