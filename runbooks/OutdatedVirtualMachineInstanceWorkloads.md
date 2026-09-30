# OutdatedVirtualMachineInstanceWorkloads

## Meaning

This alert fires when running virtual machine instances (VMIs) in outdated
`virt-launcher` pods are detected 24 hours after the KubeVirt control plane has
been updated.

## Impact

Outdated VMIs might not have access to new KubeVirt features.

Outdated VMIs will not receive the security fixes associated with the
`virt-launcher` pod update.

## Diagnosis

1. Identify the outdated VMIs:

   ```bash
   $ kubectl get vmi -l kubevirt.io/outdatedLauncherImage --all-namespaces
   ```

2. Check the `KubeVirt` custom resource (CR) to determine whether
`workloadUpdateMethods` is configured in the `workloadUpdateStrategy` stanza:

   ```bash
   $ kubectl get kubevirt --all-namespaces -o yaml
   ```

3. Check each outdated VMI to determine whether it is live-migratable:

   ```bash
   $ kubectl get vmi <vmi> -o yaml
   ```

   Example output:

   ```yaml
   apiVersion: kubevirt.io/v1
   kind: VirtualMachineInstance
   ...
     status:
       conditions:
       - lastProbeTime: null
         lastTransitionTime: null
         message: cannot migrate VMI which does not use masquerade to connect to the pod network
         reason: InterfaceNotLiveMigratable
         status: "False"
         type: LiveMigratable
   ```

## Mitigation

### Configuring automated workload updates

<!--DS: Update the `HyperConverged` CR to enable automatic workload updates.-->
<!--USstart-->
Update the `KubeVirt` CR to enable automatic workload updates.

See [Updating KubeVirt Workloads](https://kubevirt.io/user-guide/operations/updating_and_deletion/#updating-kubevirt-workloads)
for more information.
<!--USend-->

### Stopping a VM associated with a non-live-migratable VMI

If a VMI is not live-migratable and if `runStrategy: always` is set in the
corresponding `VirtualMachine` object, you can update the VMI by manually
stopping the virtual machine (VM):

  ```bash
  $ virtctl stop --namespace <namespace> <vm>
  ```

A new VMI spins up immediately in an updated `virt-launcher` pod to replace the
stopped VMI. This is the equivalent of a restart action.

Note: Manually stopping a _live-migratable_ VM is destructive and not
recommended because it interrupts the workload.

### Migrating a live-migratable VMI

If a VMI is live-migratable, you can update it by creating a
`VirtualMachineInstanceMigration` object that targets a specific running VMI.
The VMI is migrated into an updated `virt-launcher` pod.

1. Create a `VirtualMachineInstanceMigration` manifest and save it as
`migration.yaml`:

   ```yaml
   apiVersion: kubevirt.io/v1
   kind: VirtualMachineInstanceMigration
   metadata:
     name: <migration_name>
     namespace: <namespace>
   spec:
     vmiName: <vmi_name>
   ```

2. Create a `VirtualMachineInstanceMigration` object to trigger the migration:

   ```bash
   $ kubectl create -f migration.yaml
   ```

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
