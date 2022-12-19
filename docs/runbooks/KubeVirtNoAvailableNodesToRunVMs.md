# KubeVirtNoAvailableNodesToRunVMs
<!-- Edited by machadovilaca, 19 Dec 2022 -->

## Meaning

The KubeVirtNoAvailableNodesToRunVMs alert is triggered when all nodes in the Kubernetes cluster are missing hardware virtualization or CPU virtualization extensions. This means that the cluster does not have the necessary hardware support to run virtual machines (VMs).

## Impact

If this alert is triggered, it means that VMs will not be able to run on the cluster. This can have a significant impact on the operations of the cluster, as VMs may be used for critical applications or services.

## Diagnosis

To diagnose the cause of this alert, the following steps can be taken:

1. Check the hardware configuration of the nodes in the cluster. Make sure that all nodes have hardware virtualization or CPU virtualization extensions enabled.

2. Check the node labels in the cluster. Make sure that nodes with the necessary hardware support are labeled as such, so that VMs can be scheduled to run on these nodes.

## Mitigation

To mitigate the impact of this alert, the following steps can be taken:

1. Add nodes to the cluster that have hardware virtualization or CPU virtualization extensions enabled.

2. If the issue persists after attempting the above steps, contact support for further assistance.
