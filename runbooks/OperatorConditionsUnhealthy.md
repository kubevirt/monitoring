# OperatorConditionsUnhealthy [Deprecated]

This alert has been deprecated; it does not indicate a genuine issue. If
triggered, it may be safely ignored and silenced.

## Meaning

This alert triggers when the operator conditions or its secondary resources
are in an error or warning state.

## Impact

Resources maintained by the operator might not be functioning correctly.

## Diagnosis

Check the operator conditions:

```bash
kubectl get <CR> <CR_OBJECT> -n <namespace> -o jsonpath='{.status.conditions}'
```

For example:

```bash
kubectl get HyperConverged kubevirt-hyperconverged -n kubevirt -o jsonpath='{.status.conditions}'
```

## Mitigation

Based on the information obtained during the diagnosis procedure, try to
identify the root cause within the operator or any of its secondary resources,
and resolve the issue.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->

<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
