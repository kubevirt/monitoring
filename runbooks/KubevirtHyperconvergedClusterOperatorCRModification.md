# KubevirtHyperconvergedClusterOperatorCRModification

## Meaning

Hyperconverged Cluster Operator configures clusters in an opinionated way and overwrites its operands when there is an unexpected change to them. Users are expected to not modify the operands directly. The `HyperConverged` object is the source of truth for the configuration.

This alert fires when an operand of Hyperconverged Cluster Operator is changed by someone else.


## Impact

When the operands are changed by someone else constantly, it may lead to oscillation in the cluster configuration and instabilities.


## Diagnosis

Check the alert details. `component_name` contains the type and name of the operand that is overwritten.

```
Labels
    alertname=KubevirtHyperconvergedClusterOperatorCRModification
    component_name=kubevirt-kubevirt-hyperconverged
    severity=warning
```
In the example above, `kubevirt/kubevirt-hyperconverged` object is the operand for which the alert is triggered.


## Mitigation

Stop changing operands of HyperConverged Cluster operators directly and use `HyperConverged` objects to configure the cluster. The alert is supposed to resolve after 10 minutes if there isn't a manual intervention to operands in the last 10 minutes. 
