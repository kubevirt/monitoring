# LowKVMNodesCount

## Meaning

Low number of nodes with KVM resource available.

Alert will fire if less than 2 nodes in the cluster have KVM resource available.

## Impact

Virtual Machine will not be scheduled and run if there are no nodes with KVM resource available.

Virtual Machine can not be migrated if less than 2 nodes in the cluster have KVM resource available.

## Diagnosis

Verify that nodes have KVM resource available.
```
kubectl get nodes -o jsonpath='{.items[*].status.allocatable}' | grep devices.kubevirt.io/kvm
```

## Mitigation

[Validate hardware virtualization support](https://kubevirt.io/user-guide/operations/installation/#validate-hardware-virtualization-support)

If hardware virtualization is not available, then a [software emulation fallback](https://github.com/kubevirt/kubevirt/blob/master/docs/software-emulation.md) can be enabled.

