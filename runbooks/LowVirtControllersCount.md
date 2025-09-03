# LowVirtControllersCount

## Meaning

This alert fires when a low number of `virt-controller` pods is detected. At
least one `virt-controller` pod must be available in order to ensure high
availability. The default number of replicas is 2.

A `virt-controller` device monitors the custom resource definitions (CRDs) of a
virtual machine instance (VMI) and manages the associated pods. The device
create pods for VMIs and manages the lifecycle of the pods. The device is
critical for cluster-wide virtualization functionality.

## Impact

The responsiveness of KubeVirt might become negatively affected. For example,
certain requests might be missed.

In addition, if another `virt-controller` instance terminates unexpectedly,
KubeVirt might become completely unresponsive.

## Diagnosis

1. Verify that running `virt-controller` pods are available:

   ```bash
   $ kubectl -n kubevirt get pods -l kubevirt.io=virt-controller
   ```

2. Check the `virt-controller` logs for error messages:

   ```bash
   $ kubectl -n kubevirt logs <virt-controller>
   ```

3. Obtain the details of the `virt-controller` pod to check for status conditions
such as unexpected termination or a `NotReady` state.

   ```bash
   $ kubectl -n kubevirt describe pod/<virt-controller>
   ```

## Mitigation

This alert can have a variety of causes, including:

- Not enough memory on the cluster
- Nodes are down
- The API server is overloaded. For example, the scheduler might be under a
heavy load and therefore not completely available.
- Networking issues

Identify the root cause and fix it, if possible.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
