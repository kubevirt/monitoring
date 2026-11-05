# KubeVirtVMIOrdinalPodInterfaceName

## Meaning

This alert fires when a running virtual machine instance (VMI) still uses the
deprecated *ordinal* naming scheme (`net1`, `net2`, and so on) for its secondary
network pod interfaces.

VMs started on KubeVirt versions prior to <!--USstart-->v1.0<!--USend--><!--DS: 4.14--> use this ordinal scheme, while
newer VMs use a hashed scheme (for example, `pod914f438d88d`). The condition
must hold for 5 minutes before the alert fires.

## Impact

This is a `warning` alert with no impact on operator health. The affected VM
continues to run normally for now.

The ordinal scheme is being discontinued: on KubeVirt versions after <!--USstart-->v1.10<!--USend--><!--DS: 5.1-->, live
migration of these VMs breaks. A VM left on the ordinal scheme can therefore
block node maintenance and cluster upgrades once the cluster is upgraded past
<!--USstart-->v1.10<!--USend--><!--DS: 5.1-->. The VM remains in this state until it is restarted on KubeVirt <!--USstart-->v1.0<!--USend--><!--DS: 4.14--> or
later, or live-migrated on KubeVirt <!--USstart-->v1.10<!--USend--><!--DS: 5.1-->.

As a minor, additional limitation, a VM on the ordinal naming scheme does not
support NIC hot-unplug.

## Diagnosis

1. Identify the affected VMI name and namespace from the alert description.

2. Inspect the pod interface names reported in the VMI status. Secondary
   interfaces named `net` followed by digits (for example, `net1`) use the
   deprecated ordinal scheme:

   ```bash
   $ kubectl get vmi <name> -n <namespace> \
       -o jsonpath='{range .status.interfaces[*]}{.name}{"="}{.podInterfaceName}{"\n"}{end}'
   ```

## Mitigation

On KubeVirt <!--USstart-->v1.0<!--USend--><!--DS: 4.14--> or later, restart the VM. The recreated VMI adopts the hashed
naming scheme:

```bash
$ virtctl restart <name> -n <namespace>
```

Note: on KubeVirt <!--USstart-->v1.10<!--USend--><!--DS: 5.1-->, you can
live-migrate the VM once instead, without downtime. Check that the VMI is
live-migratable:

```bash
$ kubectl get vmi <name> -n <namespace> \
    -o jsonpath='{.status.conditions[?(@.type=="LiveMigratable")].status}{"\n"}'
```

If it reports `True`, live-migrate it:

```bash
$ virtctl migrate <name> -n <namespace>
```

These commands act on a VirtualMachine. For a standalone VMI, recreate the VMI
by deleting it and re-applying its manifest.

After the VM adopts the hashed scheme, the
`kubevirt_vmi_vnic_info{pod_interface_name=~"net[0-9]+"}` series clears and the
alert resolves.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Red Hat Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
