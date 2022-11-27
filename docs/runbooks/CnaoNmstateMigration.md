# CnaoNmstateMigration
<!-- Edited by Jiří Herrmann, 8 Nov 2022 -->

## Meaning

<!--USstart-->
This alert fires when a `kubernetes-nmstate` deployment is detected and the Kubernetes NMState Operator is not installed. This alert only affects OKD 1.6.

The Cluster Network Addons Operator (CNAO) does not support `kubernetes-nmstate` deployments in OKD 1.7 and later.
<!--USend-->
<!--DS: This alert fires when a `kubernetes-nmstate` deployment is detected and the Kubernetes NMState Operator is not installed. This alert only affects {VirtProductName} 4.10.-->

<!--DS: The Cluster Network Addons Operator (CNAO) does not support `kubernetes-nmstate` deployments in {VirtProductName} 4.11 and later.-->

## Impact

<!--DS: You cannot upgrade to {VirtProductName} 4.11.-->
<!--USstart-->
You cannot upgrade to OKD 1.7.
<!--USend-->

## Mitigation

Install the Kubernetes NMState Operator from the OperatorHub. CNAO automatically transfers the `kubernetes-nmstate` deployment to the Operator. 