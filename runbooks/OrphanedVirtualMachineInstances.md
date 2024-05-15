# OrphanedVirtualMachineInstances

## Meaning

This alert fires when a virtual machine instance (VMI), or `virt-launcher` pod,
runs on a node that does not have a running `virt-handler` pod. Such a VMI is
called _orphaned_.

## Impact

Orphaned VMIs cannot be managed.

## Diagnosis

1. Check the status of the `virt-handler` pods to view the nodes on which they
are running:

   ```bash
   $ kubectl get pods --all-namespaces -o wide -l kubevirt.io=virt-handler
   ```

2. Check the status of the VMIs to identify VMIs running on nodes that do not
have a running `virt-handler` pod:

   ```bash
   $ kubectl get vmis --all-namespaces
   ```

3. Check the status of the `virt-handler` daemon:

   ```bash
   $ kubectl get daemonset virt-handler --all-namespaces
   ```

   Example output:

   ```text
   NAME                  DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR            AGE
   virt-handler          2         2         2       2            2           kubernetes.io/os=linux   4h
   ```

   The daemon set is considered healthy if the `Desired`, `Ready`, and
   `Available` columns contain the same value.

4. If the `virt-handler` daemon set is not healthy, check the `virt-handler`
daemon set for pod deployment issues:

   ```bash
   $ kubectl get daemonset virt-handler --all-namespaces -o yaml | jq .status
   ```

5. Check the nodes for issues such as a `NotReady` status:

   ```bash
   $ kubectl get nodes
   ```

6. Check the `spec.workloads` stanza of the `KubeVirt` custom resource (CR) for
a workloads placement policy:

   ```bash
   $ kubectl get kubevirt kubevirt --all-namespaces -o yaml
   ```

## Mitigation

If a workloads placement policy is configured, add the node with the VMI to the
policy.

Possible causes for the removal of a `virt-handler` pod from a node include
changes to the node's taints and tolerations or to a pod's scheduling rules.

Try to identify the root cause and resolve the issue.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
See [How Daemon Pods are scheduled](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/#how-daemon-pods-are-scheduled)
for more information.

If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
