# VirtualMachineSnapshotFailed

## Meaning

This alert fires when a `VirtualMachineSnapshot` has not become ready to use
for more than 15 minutes.

The `kubevirt_vmsnapshot_info` metric reports one series per
`VirtualMachineSnapshot`. The `name` label is the snapshot name and the `vm`
label is the name of its source virtual machine (VM). The alert fires when
the `ready_to_use` label of a series stays `false` for 15 minutes.

By default, a snapshot that does not complete within 5 minutes
(`spec.failureDeadline`) is marked as failed with the condition reason
`snapshot deadline exceeded`. A failed snapshot object remains in the cluster
until it is deleted. If `spec.failureDeadline` is set to 0, the deadline is
disabled and a stuck snapshot stays in progress indefinitely.

## Impact

The snapshot cannot be used to restore the VM or to export its disks.

The VM itself keeps running and is not damaged. When the snapshot fails, the
snapshot controller unfreezes the VM and cleans up the snapshot content. The
cluster and other VMs are not affected.

While the snapshot is in progress, the source VM is protected from deletion.
An online snapshot taken while the VM is running freezes the VM's disks for
the duration of the snapshot, which degrades the VM's I/O.

Because a failed snapshot object is never removed automatically, the alert
keeps firing until the failed snapshot is deleted.

## Diagnosis

Set environment variables from the alert labels:

```bash
$ export NAMESPACE="<namespace label>"
$ export SNAPSHOT_NAME="<name label>"
$ export VM_NAME="<vm label>"
```

1. Get the snapshot status and error:

   ```bash
   $ kubectl get vmsnapshot $SNAPSHOT_NAME -n $NAMESPACE -o yaml
   ```

   Check `status.error.message`, `status.conditions` and `status.phase`. For
   example, the condition reason `snapshot deadline exceeded` means the
   snapshot did not complete within `spec.failureDeadline`.

2. Check whether the VM's volumes can be snapshotted:

   ```bash
   $ kubectl get vm $VM_NAME -n $NAMESPACE -o yaml
   ```

   `status.volumeSnapshotStatuses` lists each volume with `enabled: false`
   and the reason:

   - `Snapshot is not supported for this volumeSource type`: the volume is
     not backed by a persistent volume claim (PVC), a DataVolume or a memory
     dump volume.
   - `No VolumeSnapshotClass`: the storage class of the volume has no
     matching VolumeSnapshotClass.

3. Check the volume snapshots that the Container Storage Interface (CSI)
   driver created for the VM:

   ```bash
   $ kubectl get volumesnapshots -n $NAMESPACE \
     -l snapshot.kubevirt.io/source-vm-name=$VM_NAME
   ```

   Describe each volume snapshot that is not ready to use. Its
   `status.error` reports why the CSI driver could not create the snapshot,
   for example because the storage backend does not support snapshots.

4. Check the logs of the snapshot controller, which runs in the
   `virt-controller` deployment:

   ```bash
   $ export KV_NAMESPACE="$(kubectl get kubevirt -A \
       -o custom-columns="":.metadata.namespace)"
   $ kubectl logs -n $KV_NAMESPACE -l kubevirt.io=virt-controller \
     | grep $SNAPSHOT_NAME
   ```

5. Show the phase and readiness of the snapshot:

<!--USstart-->
   **Note:** You can use the Prometheus metrics explorer to run the query.
<!--USend-->

<!--DS: **Note:** You can run the query in the web console by navigating to
**Observe** -> **Metrics**.-->

   ```promql
   kubevirt_vmsnapshot_info{namespace="<namespace>", name="<name>"}
   ```

## Mitigation

1. Address the cause of the failure:

   - If the condition reason is `snapshot deadline exceeded`, the volume
     snapshots did not complete in time. Check the storage backend and the
     volume snapshots, and consider increasing `spec.failureDeadline` in the
     `VirtualMachineSnapshot`.
   - If a volume's storage class has no matching VolumeSnapshotClass, ask
     the cluster administrator to configure one. KubeVirt selects the
     VolumeSnapshotClass from the CDI StorageProfile, or by matching the
     storage class provisioner, preferring the class annotated with
     `snapshot.storage.kubernetes.io/is-default-class: "true"`.
   - If a volume is not backed by a PVC, a DataVolume or a memory dump
     volume, remove it from the VM before taking a snapshot.

2. Delete the failed snapshot:

   ```bash
   $ kubectl delete vmsnapshot $SNAPSHOT_NAME -n $NAMESPACE
   ```

   By default (`spec.deletionPolicy: Delete`), the snapshot content and the
   volume snapshots are deleted with it.

3. Recreate the snapshot and wait until `status.readyToUse` is `true`.

<!--DS: Create and delete snapshots by following
[Backup and restore by using VM snapshots](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/virtualization/backup-and-restore#virt-backup-restore-snapshots).-->
<!--USstart-->
For more information, see
[Virtual machine snapshots](https://kubevirt.io/user-guide/storage/snapshot_restore_api/).
<!--USend-->

The alert resolves when the snapshot becomes ready to use.

<!--DS: If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
