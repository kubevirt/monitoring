# LowVirtOperatorCount 

## Meaning

This alert is triggered when there is only one virt-operator pod that is running in `Ready` state in the past 10 minutes. 

## Impact

The virt-operator is the first k8s operator coming alive in a KubeVirt cluster. Its primary responsibilities are installation, live-update, live-upgrade of a KubeVirt cluster, monitoring the life-cycle of top-level controllers, such as virt-controller, virt-handler, virt-launcher, etc. and manage their reconciliation. In addition, virt-operator is responsible for cluster-wide tasks, such as certificate rotation and some infrastructure management, etc.  Note that virt-operator is not directly responsible for virtual machines in the cluster, its temporal unavailability should not affect the custom workloads. 

More than one virt-operator pod should be in the Ready state when KubeVirt is deployed with high-availability. The virt-operator deployment has a default replica of 2 pods.

This alter indicates the downgrade of the availability of virt-operator.

## Diagnosis

Check the states of virt-operator pods:
- `export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"`
- `kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-operator`

Check in-depth virt-operator pods that are in trouble:
- `kubectl -n $NAMESPACE logs <pod-name>`.
- `kubectl -n $NAMESPACE describe pod <pod-name>`.

## Mitigation

There can be several reasons, try to identify the root cause and fix it. If you cannot fix it, please open an issue and attach the artifacts gathered in the Diagnosis section.
