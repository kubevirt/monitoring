# VirtOperatorDown 

## Meaning

The virt-operator is the first k8s operator coming alive in a KubeVirt cluster. Its primary responsibilities are installation, live-update, live-upgrade of a KubeVirt cluster, monitoring the life-cycle of top-level controllers, such as virt-controller, virt-handler, virt-launcher, etc. and manage their reconciliation. In addition, virt-operator is responsible for cluster-wide tasks, such as certificate rotation and some infrastructure management, etc.

Note that virt-operator is not directly responsible for virtual machines in the cluster, its temporal unavailability should not affect the custom workloads. 

This alert is triggered when there was no virt-operator pod detected in the `Running` state in the past 5 minutes. The virt-operator deployment has a default replica of 2 pods.

## Impact

This alert indicates a failure at the level of the KubeVirt cluster. Critical cluster-wide management functionalities such as certification rotation, KubeVirt upgrade, and reconciliation of KubeVirt controllers, are not available for the time being.

## Diagnosis

Check if the output of the following command is `0`. 
- `export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"`
- `kubectl get deployment -n $NAMESPACE virt-operator -o jsonpath='{.status.readyReplicas}'` 

Check if the following command shows that there is virt-operator pod in the `Running` state.
- `export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"`
- `kubectl get pods -n $NAMESPACE -l=kubevirt.io=virt-operator`

## Mitigation

- Check the status of the virt-operator deployment to find out more information. The following commands will provide the associated events and show if there are any specific issues.
  - `export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"`
  - `kubectl -n $NAMESPACE get deploy virt-operator -o yaml`
  - `kubectl -n $NAMESPACE describe deploy virt-operator`
- Check the status of the virt-operator pods for further information: 
  - `export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"`
  - `kubectl get pods -n $NAMESPACE -l=kubevirt.io=virt-operator`
- Check if there are issues with the nodes for control-plane and masters. For example, if they are in a NotReady state.
  - `kubectl get nodes`

