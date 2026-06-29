# KubevirtHyperconvergedClusterOperatorNMOInUseAlert [Deprecated]

This alert has been deprecated; it does not indicate a genuine issue. If
triggered, it may be safely ignored and silenced.

## Meaning

<!--DS: This alert fires when _integrated_ Node Maintenance Operator (NMO)
custom resources (CRs) are detected. This alert only affects {VirtProductName}
4.10.-->

<!--DS: The Node Maintenance Operator is not included with {VirtProductName}
4.11.0 or later. Instead, the Operator is installed from OperatorHub.-->

<!--DS: The presence of `NodeMaintenance` CRs belonging to the
`nodemaintenance.kubevirt.io` API group indicates that the node specified in
`spec.nodeName` was put into maintenance mode. The target node has been cordoned
off and drained.-->

<!--USstart-->
This alert fires when _integrated_ Node Maintenance Operator (NMO) custom
resources (CRs) are detected. This alert only affects OKD 1.6.

The presence of `NodeMaintenance` CRs belonging to the
`nodemaintenance.kubevirt.io` API group indicates that the node specified in
`spec.nodeName` was put into maintenance mode. The target node has been
[cordoned off](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#cordon)
and [drained](https://kubernetes.io/docs/tasks/administer-cluster/safely-drain-node/#use-kubectl-drain-to-remove-a-node-from-service).
<!--USend-->

## Impact

<!--DS: You cannot upgrade to {VirtProductName} 4.11.-->
<!--USstart-->
You cannot upgrade to OKD 1.7.
<!--USend-->

## Diagnosis

1. Check the `kubevirt-hyperconverged` resource to see whether it is upgradeable:

   ```bash
   $ kubectl get hco kubevirt-hyperconverged -o json | jq -r '.status.conditions[] | select(.type == "Upgradeable")'
   ```

   Example output:

   ```json
   {
     "lastTransitionTime": "2022-05-26T09:23:21Z",
     "message": "NMO custom resources have been found",
     "reason": "UpgradeBlocked",
     "status": "False",
     "type": "Upgradeable"
   }
   ```

2. Check for a ClusterServiceVersion (CSV) warning event such as the following:

   ```text
   Warning  NotUpgradeable      2m12s (x5 over 2m50s)   kubevirt-hyperconvergedNode
   Maintenance Operator custom resources nodemaintenances.nodemaintenance.kubevirt.io
   have been found.
   Please remove them to allow upgrade. You can use NMO standalone operator if
   keeping the node(s) under maintenance is still required.
   ```

3. Check for NMO CRs belonging to the `nodemaintenance.kubevirt.io` API group:

   ```bash
   $ kubectl get nodemaintenances.nodemaintenance.kubevirt.io
   ```

   Example output:

   ```text
   NAME                   AGE
   nodemaintenance-test   5m33s
   ```

## Mitigation

Remove all NMO CRs belonging to the
`nodemaintenance.nodemaintenance.kubevirt.io/` API group. After the integrated
NMO resources are removed, the alert is cleared and you can upgrade.

If a node must remain in maintenance mode during upgrade, install the Node
Maintenance Operator from OperatorHub. Then, create an NMO CR belonging to the
`nodemaintenance.nodemaintenance.medik8s.io/v1beta1` API group and version for
the node.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
See the [HCO cluster configuration documentation](https://github.com/kubevirt/hyperconverged-cluster-operator/blob/main/docs/cluster-configuration.md#enablecommonbootimageimport-feature-gate)
for more information.

If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
