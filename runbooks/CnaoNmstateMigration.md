# CnaoNmstateMigration

## Meaning

This alert fires when a `kubernetes-nmstate` deployment is detected and the
Kubernetes NMState Operator is not installed. This alert only affects OpenShift
Virtualization 4.10.

The Cluster Network Addons Operator (CNAO) does not support `kubernetes-nmstate`
deployments in OpenShift Virtualization 4.11 and later.

## Impact

You cannot upgrade your cluster to OpenShift Virtualization 4.11.

## Mitigation

Install the Kubernetes NMState Operator from the OperatorHub. CNAO automatically
transfers the `kubernetes-nmstate` deployment to the Operator.

Afterwards, you can upgrade to OpenShift Virtualization 4.11.

<!--DS: If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
