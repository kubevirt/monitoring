# LowKVMNodesCount

## Meaning

This alert fires when fewer than two nodes in the cluster have KVM resources.

## Impact

The cluster must have at least two nodes with KVM resources for live migration.

Virtual machines cannot be scheduled or run if no nodes have KVM resources.

## Diagnosis

- Identify the nodes with KVM resources:

  ```bash
  $ kubectl get nodes --o custom-columns=NAME:.metadata.name,KVM:".status.allocatable.devices\.kubevirt\.io/kvm"
  ```

## Mitigation

<!--DS: Install KVM on the nodes without KVM resources.-->
<!--USstart-->
Validate the [hardware virtualization support](https://kubevirt.io/user-guide/operations/installation/#validate-hardware-virtualization-support).

If hardware virtualization is not available, [software emulation](https://github.com/kubevirt/kubevirt/blob/master/docs/software-emulation.md)
can be enabled.
<!--USend-->

<!--DS: If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
