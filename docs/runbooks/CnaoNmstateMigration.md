<!-- Edited by Jiří Herrmann, 8 Nov 2022 -->

# CnaoNmstateMigration

## Meaning

This alert fires when the `nmstate` API is deployed, but the Kubernetes NMState Operator is not installed.

In a future update, the `Cluster-network-addons-operator` (CNAO) will stop supporting `kubernetes-nmstate` deployments. Instead of using the CNAO to deploy `kubernetes-nmstate`, it is recommended to install the stand-alone Kubernetes NMState Operator.

## Impact

In a future update, using the CNAO to deploy `kubernetes-nmstate` will cause upgrading your cluster to fail.

## Mitigation

Install the stand-alone Kubernetes NMState Operator.

Afterwards, the CNAO automatically transfers the `kubernetes-nmstate` deployment to the Kubernetes NMState Operator.

