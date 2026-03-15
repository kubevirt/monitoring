# VirtHandlerDown

## Meaning

This alert fires when all `virt-handler` pods are down (no `virt-handler` pod
has been reported up for 10 minutes).

The `virt-handler` runs on every node that can schedule VMIs. It is
responsible for domain lifecycle, network configuration, and other
node-level operations for virtual machine instances.

## Impact

Virtual machine instances (VMIs) on affected nodes cannot be managed properly.
New VMIs may not start on nodes without a running `virt-handler`, and
existing VMIs may not receive updates or clean shutdowns.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-handler` DaemonSet and pods:

   ```bash
   $ kubectl -n $NAMESPACE get daemonset virt-handler -o yaml
   $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-handler
   ```

3. Check DaemonSet events and pod status:

   ```bash
   $ kubectl -n $NAMESPACE describe daemonset virt-handler
   $ kubectl -n $NAMESPACE describe pod -l kubevirt.io=virt-handler
   ```

4. Check for node issues (e.g. nodes not ready or taints):

   ```bash
   $ kubectl get nodes
   ```

5. If any `virt-handler` pod exists, review its logs:

   ```bash
   $ kubectl -n $NAMESPACE logs <virt-handler-pod-name> --previous
   $ kubectl -n $NAMESPACE logs <virt-handler-pod-name>
   ```

## Mitigation

Identify why `virt-handler` pods are down (e.g. DaemonSet not scheduling, pods
crashing, node issues, image pull failures) and restore the DaemonSet so
`virt-handler` runs on schedulable nodes.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
