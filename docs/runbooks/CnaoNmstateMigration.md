<!-- Edited by Jiří Herrmann, 8 Nov 2022 -->

# CnaoNmstateMigration

## Meaning

This alert fires when the `nmstate` API is deployed, but the Kubernetes NMState Operator is not installed.

The Cluster Network Addons Operator (CNAO) does not support `kubernetes-nmstate` deployments in OpenShift Virtualization 4.11 and later.

This alert only affects OpenShift Virtualization 4.10.

## Impact

You cannot upgrade your cluster to OpenShift Virtualization 4.11.

## Mitigation

Install the Kubernetes NMState Operator from the OperatorHub. CNAO automatically transfers the `kubernetes-nmstate` deployment to the Operator. 

Afterwards, you can upgrade to OpenShift Virtualization 4.11.