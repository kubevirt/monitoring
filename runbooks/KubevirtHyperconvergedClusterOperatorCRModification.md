# KubeVirtCRModified

## Meaning

This alert fires when an operand of the HyperConverged Cluster Operator (HCO) is
changed by someone or something other than HCO.

HCO configures KubeVirt and its supporting operators in an opinionated way and
overwrites its operands when there is an unexpected change to them. Users must
not modify the operands directly. The `HyperConverged` custom resource is the
source of truth for the configuration.

## Impact

Changing the operands manually causes the cluster configuration to fluctuate and
might lead to instability.

## Diagnosis

Check the `component_name` in the alert details to determine the operand that is
being changed.

In the following example, the operand kind is `kubevirt` and the operand name is
`kubevirt-kubevirt-hyperconverged`:

```text
Labels
   alertname=KubeVirtCRModified
   component_name=kubevirt/kubevirt-kubevirt-hyperconverged
   severity=warning
```

## Mitigation

Do not change the HCO operands directly. Use `HyperConverged` objects to
configure the cluster.

The alert resolves itself after 10 minutes if the operands are not changed
manually.
