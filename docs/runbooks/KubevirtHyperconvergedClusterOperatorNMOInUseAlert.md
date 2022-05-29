# KubevirtHyperconvergedClusterOperatorNMOInUseAlert

## Meaning
Starting from version 4.11.0, Node Maintenance Operator (NMO) is no longer included in OpenShift Virtualization.  
This alert is being fired if both of the following conditions are met:
- OpenShift version is 4.11 or higher.
- The integrated NMO is being used, i.e. one or more custom resources of Kind `NodeMaintenance` and API group of `nodemaintenance.kubevirt.io` have been found on the cluster.

Existence of such a custom resource indicates that the targeted node specified in `spec.nodeName` was put into a maintenance mode, which means the node is [cordoned](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#cordon) and server-side [drained](https://kubernetes.io/docs/tasks/administer-cluster/safely-drain-node/#use-kubectl-drain-to-remove-a-node-from-service).

## Impact
Upgrade to OpenShift Virtualization v4.11 is blocked until all integrated NMO custom resources (`nodemaintenance.kubevirt.io`) are removed.

## Diagnosis
- `Upgradeable` condition in HyperConverged CR is false:
```bash
$ kubectl get hco kubevirt-hyperconverged -o json | jq -r '.status.conditions[] | select(.type == "Upgradeable")'
{
  "lastTransitionTime": "2022-05-26T09:23:21Z",
  "message": "NMO custom resources have been found",
  "reason": "UpgradeBlocked",
  "status": "False",
  "type": "Upgradeable"
}
```

- A warning event is emitted to the ClusterServiceVersion (CSV):
```bash
  Warning  NotUpgradeable      2m12s (x5 over 2m50s)   kubevirt-hyperconverged     Node Maintenance Operator custom resources nodemaintenances.nodemaintenance.kubevirt.io have been found.
  Please remove them to allow upgrade. You can use NMO standalone operator if keeping the node(s) under maintenance is still required. 
```

- NMO custom resource(s) from API group `nodemaintenance.kubevirt.io` exists:
```bash
$ kubectl get nodemaintenances.nodemaintenance.kubevirt.io
NAME                   AGE
nodemaintenance-test   5m33s 
```
## Mitigation
1. If the targeted node(s) should remain in maintenance mode throughout the upgrade of OpenShift Virtualization, please install the standalone Node Maintenance Operator from OperatorHub and create a corresponding custom resource targeting the same node(s), with API group and version of `nodemaintenance.nodemaintenance.medik8s.io/v1beta1`. Then, remove all previous NMO CRs from API Group `nodemaintenance.nodemaintenance.kubevirt.io/v1beta1`.
2. If the maintenance mode on these nodes is no longer required, just remove the old NMO CRs. You can install the standalone Node Maintenance Operator independently of OpenShift Virtualization.

Once all old (`nodemaintenance.nodemaintenance.kubevirt.io`) NMO custom resources are removed, the alert is cleared and the upgrade to OpenShift Virtualization v4.11 is no longer blocked by this reason.