# NoReadyVirtHandler

## Meaning

This alert fires when no `virt-handler` pod in a `Ready` state has been
detected for 5 minutes.

The `virt-handler` runs on every node that can schedule VMIs (as a
DaemonSet). It is responsible for domain lifecycle and node-level operations
for virtual machine instances.

## Impact

This is a warning level alert. Virtual machine instances may experience
minor delays until at least one `virt-handler` becomes ready.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-controller` deployment and pods:

   ```bash
   $ kubectl -n $NAMESPACE get deployment virt-controller
   $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-controller -o wide
   ```

3. Check deployment and pod events:

   ```bash
   $ kubectl -n $NAMESPACE describe deployment virt-controller
   $ kubectl -n $NAMESPACE describe pod -l kubevirt.io=virt-controller
   ```

4. Review logs of any running but not ready `virt-controller` pod:

   ```bash
   $ kubectl -n $NAMESPACE logs <virt-controller-pod-name> --previous
   $ kubectl -n $NAMESPACE logs <virt-controller-pod-name>
   ```

5. Check for cluster-wide node or scheduling issues:

   ```bash
   $ kubectl get nodes
   ```

## Mitigation

Identify the root cause (e.g. DaemonSet not scheduling, all pods crashing or
failing readiness, node or image issues) and restore at least one ready
`virt-handler` pod on a schedulable node.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
