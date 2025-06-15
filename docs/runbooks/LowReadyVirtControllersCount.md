# LowReadyVirtControllersCount

## Meaning

This alert fires when one or more `virt-controller` pods are running, but none
of these pods has been in the `Ready` state for the last 5 minutes.

A `virt-controller` device monitors the custom resource definitions (CRDs) of a
virtual machine instance (VMI) and manages the associated pods. The device
create pods for VMIs and manages the lifecycle of the pods. The device is
critical for cluster-wide virtualization functionality.

## Impact

This alert indicates that a cluster-level failure might occur, which would cause
actions related to VM lifecycle management to fail. This notably includes
launching a new VMI or shutting down an existing VMI.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Verify a `virt-controller` device is available:

   ```bash
   $ kubectl get deployment -n $NAMESPACE virt-controller -o jsonpath='{.status.readyReplicas}'
   ```

3. Check the status of the `virt-controller` deployment:

   ```bash
   $ kubectl -n $NAMESPACE get deploy virt-controller -o yaml
   ```

4. Obtain the details of the `virt-controller` deployment to check for status
conditions, such as crashing pods or failures to pull images:

   ```bash
   $ kubectl -n $NAMESPACE describe deploy virt-controller
   ```

5. Check if any problems occurred with the nodes. For example, they might be in
a `NotReady` state:

   ```bash
   $ kubectl get nodes
   ```

## Mitigation

This alert can have multiple causes, including the following:

- The cluster has insufficient memory.
- The nodes are down.
- The API server is overloaded. For example, the scheduler might be under a
heavy load and therefore not completely available.
- There are network issues.

Try to identify the root cause and resolve the issue.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
