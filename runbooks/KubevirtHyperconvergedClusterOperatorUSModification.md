# KubevirtHyperconvergedClusterOperatorUSModification

## Meaning

The HyperConverged Cluster Operator configures kubevirt and its supporting operators in an opinionated way and overwrites
its operands when there is an unexpected change to them. Users are expected to not modify the operands directly. If 
such a change is a must, and is not supported by HCO API, it is possible to use special annotations to force HCO to set
the required changes in the required operator. For more details, see
[here](https://github.com/kubevirt/hyperconverged-cluster-operator/blob/main/docs/cluster-configuration.md#jsonpatch-annotations)

However, this method may be risky and considered unsafe. 

HCO fires this alert  when one of these annotations are in use.

## Impact

Just an information about using of unsafe modifications.


## Diagnosis
Check the alert details. `annotation_name` refers to specific unsafe jsonpatch annotation in the HyperConverged resource.

```
Labels
    alertname=KubevirtHyperconvergedClusterOperatorUSModification
    annotation_name=kubevirt.kubevirt.io/jsonpatch
    severity=info
```
In the example above, the kind of specific jsonpatch annotation is `kubevirt.kubevirt.io/jsonpatch`.

## Mitigation
Make sure that the change that is defined in the jsonpatch annotation is required, and that using the jsonpatch annotation
is the only way to make this change. If the HyperConverged API supports this modification, please prefer modifying the
HyperConverged resource.

If you believe this jsonpatch annotation is generic and can benefit others, please apply an RFE to add a new API to
support this modification, by filing a bug in [Red Hat bugzilla](https://bugzilla.redhat.com/).

If there is no other way, consider this alert as information only.
