# LowVirtOperatorCount 

## Meaning

The virt-operator is the first k8s operator coming to alive in a KubeVirt cluster. With the primary responsible for the installation, live-update, live-upgrade of a KubeVirt cluster, it monitors the life-cycle of top-level controllers such as virt-controller, virt-handler, virt-launcher, etc. and manages their reconciliation. In addition, virt-operator is responsible for cluster-wide tasks such as certificate rotation and some infrastructure management, etc.

Note that virt-operator is not directly responsible for virtual machines in the cluster, its temporal unavailability should not affect the custom workloads. 

More than one virt-operator pod should be in the Ready state when KubeVirt is deployed with high-availability. The virt-operator deployment has a default replica of 2 pods.

This alert is triggered when there is only one virt-operator pod is running in `Ready` state in the past 5 minutes. 

## Impact

This alter indicates the downgrade of the availability of virt-operator.

## Diagnosis

Check the states of virt-operator pods:
- `kubectl -n kubevirt get pods -l kubevirt.io=virt-operator`

Check in-depth virt-operator pods that are in trouble:
- `kubectl -n kubevirt logs <pod-name>`.
- `kubectl -n kubevirt describe pod <pod-name>`.

## Mitigation

There can be several reasons, identify the root cause and fix it.

