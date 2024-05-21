# LowKVMNodesCount

## Meaning

This alert fires when fewer than two nodes in the cluster have KVM resources.

## Impact

The cluster must have at least two nodes with KVM resources for live migration.

Virtual machines cannot be scheduled or run if no nodes have KVM resources.

## Diagnosis

- Identify the nodes with KVM resources:

  ```bash
  $ kubectl get nodes -o jsonpath='{.items[*].status.allocatable}' | grep devices.kubevirt.io/kvm
  ```

## Mitigation

<!--DS: Install KVM on the nodes without KVM resources.-->
<!--USstart-->
Validate the [hardware virtualization support](https://kubevirt.io/user-guide/operations/installation/#validate-hardware-virtualization-support).

If hardware virtualization is not available, [software emulation](https://github.com/kubevirt/kubevirt/blob/master/docs/software-emulation.md)
can be enabled.
<!--USend-->
