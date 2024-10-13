# DeprecatedMachineType

## Meaning
This alert fires when one or more Virtual Machines (VMs)  
are running with machine types that have been marked as deprecated (no  
longer supported).

## Impact
Running VMs on deprecated machine types can cause serious issues,  
including:
- Loss of certain features or degraded functionality.
- Reduced performance or unexpected behavior.
- Security vulnerabilities due to lack of ongoing support.
- Failure to restart VMs if the machine type is unsupported by the  
  current or upgraded virt-launcher version.

These issues can lead to service disruptions and affect workload  
availability during or after cluster upgrades.

## Diagnosis
This alert does **not** specify which exact machine types are deprecated  
on your VMs. Instead, it indicates that some VMs or VMIs are using  
machine types marked as deprecated by the cluster nodes.

To identify which machine types are deprecated on your nodes, inspect  
the following file on each node:  
`/var/lib/kubevirt-node-labeller/capabilities.xml`

This XML file contains the node's virtualization capabilities,  
including a list of guest machine types. Deprecated machine types are  
marked with the attribute `deprecated="yes"`.

Example excerpt from the XML:
```xml
<machine maxCpus='710' deprecated='yes'>pc-q35-rhel8.6.0</machine>
<machine maxCpus='4096'>pc-q35-rhel9.6.0</machine>
```
In this example, pc-q35-rhel8.6.0 is deprecated, while pc-q35-rhel9.6.0 is
supported.

### Next steps for admins
- Compare the machine types used by your VMs with the deprecated types
  listed in this file.
- Plan and perform migrations of VMs using deprecated machine types to
  supported ones before upgrading the cluster.

## Mitigation
Update affected VMs to use a supported machine type. You can:
- Edit VM definitions individually by modifying the
  `spec.template.spec.domain.machine.type` field.
- Or, for a smoother and cleaner update of multiple VMs, use the
  `kubevirt-api-lifecycle-automation` tool to transition all deprecated VMs
  in one operation. This approach ensures consistent, automated migration
  and reduces manual errors or downtime during cluster upgrades.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/
  virtualization)
<!--USend-->
