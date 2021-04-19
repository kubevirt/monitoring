# KubevirtHyperconvergedClusterOperatorCRModification

## Meaning

Hyperconverged Cluster Operator configures kubevirt and its supporting operators in an opinionated way and overwrites its operands when there is an unexpected change to them. Users are expected to not modify the operands directly. The `HyperConverged` custom resource is the source of truth for the configuration.

This alert is fired when an operand of Hyperconverged Cluster Operator is changed by someone else.


## Impact

When the operands are changed by someone else constantly, it may lead to oscillation in the cluster configuration and instabilities.


## Diagnosis

Check the alert details. `component_name` is the name of the operand that is overwritten.

```
Labels
    alertname=KubevirtHyperconvergedClusterOperatorCRModification
    component_name=kubevirt-kubevirt-hyperconverged
    severity=warning
```
Thanks to the fact that all kubevirt-hyperconverged CRs follow `<type>-kubevirt-hyperconverged` pattern in in their names, it is possible to fetch the type of the CRs from their name. In the example above, the type is `kubevirt` and `kubevirt.kubevirt.io/kubevirt-kubevirt-hyperconverged` is the fully qualified name of the operand. 


## Mitigation

Stop changing operands of HyperConverged Cluster operators directly and use `HyperConverged` objects to configure the cluster. The alert is supposed to resolve after 10 minutes if there isn't a manual intervention to operands in the last 10 minutes. 
