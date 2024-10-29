# VMOutdatedMachineType

## Meaning
This alert fires when there are virtual machines (VMs) using machine types
that are no longer supported.
Machine types define the virtualized hardware configuration for a VM.

## Impact

Using an outdated machine type can have severe consequences.
It may prevent certain features from functioning correctly, reduce performance,
or pose security risks. In some cases, the VM might fail to restart,
especially if the machine type is no longer supported by the underlying
virt-launcher version. This can cause service disruptions and affect the
availability of workloads running in these VMs.

## Diagnosis

To list the virtual machines using outdated machine types,
run the following command:

```bash
$ kubectl get vmis -o jsonpath='{.items[?(@.spec.template.spec.domain.machine.type=="outdated")].metadata.name}'
```
Replace "outdated" with the actual machine type flagged by the alert.

## Mitigation

Ensure that all virtual machines are updated to use a supported machine type.
This can be done by editing the VM definition and setting the
spec.template.spec.domain.machine.type field to a supported machine type.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
