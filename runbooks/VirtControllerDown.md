# VirtControllerDown

## Meaning
No running `virt-controller` pod has been detected for 5 minutes.

## Impact
Any actions related to virtual machine (VM) lifecycle management fail. This
notably includes launching a new virtual machine instance (VMI) or shutting down
an existing VMI.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-controller` deployment:

   ```bash
   $ kubectl get deployment -n $NAMESPACE virt-controller -o yaml
   ```

3. Review the logs of the `virt-controller` pod:

   ```bash
   $ kubectl get logs <virt-controller>
   ```

## Mitigation

This alert can have a variety of causes, including the following:

- Node resource exhaustion
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

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
