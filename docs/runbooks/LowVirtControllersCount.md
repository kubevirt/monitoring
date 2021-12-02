# LowVirtControllersCount

## Meaning

More than one virt-controller pod should be ready to ensure high-availability. The current default number of replicas is 2.

## Impact
This impacts availability of the virt-controller, which is responsible of managing Kubevirt's resources (e.g. VMIs) and their life-cycles, hence effects Kubevirt's responsiveness as a whole.
With low number of virt-controller replicas the cluster might be less responsive (some requests may be missed, etc'), but it is not guaranteed. 
The more severe danger is if another virt-lancher instance crashes Kubevirt might be totally unresponsive.

## Diagnosis
- Set the environment variable NAMESPACE

    ```
    export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
    ```

- Check to see if there are any running virt-controller pods:
    ```
    kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-controller
    ```

- Further check not ready or crashing virt-controller pods:
    ```
    kubectl -n $NAMESPACE logs virt-launcher-<unique id>
    ```
    ```
    kubectl -n $NAMESPACE describe pod/virt-launcher-<unique id>
    ```

## Mitigation

There can be several reasons. Like:

- Not enough memory on the cluster
- Nodes are down
- API server is overloaded (e.g. Scheduler has a lot of work and therefore is not 100% available)
- Networking issues

Try to identify the root cause and fix it.

In other cases, please open an issue and attach the artifacts gathered in the Diagnosis section.