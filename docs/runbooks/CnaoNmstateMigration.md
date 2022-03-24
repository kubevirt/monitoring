# CnaoNmstateMigration

## Meaning

Cluster-network-addons-operator (CNAO) will stop support of kubernetes-nmstate deployment.
This alert fires when nmstate is deployed, but there is no Kubernetes NMState Operator installed.

## Impact

Having kubernetes-nmstate deployed via CNAO will block future upgrade.

## Mitigation

Please, install standalone Kubernetes NMState Operator. After Kubernetes NMState Operator is installed, CNAO will hand over the kubernetes-nmstate deployment to the NMState Operator automatically.
