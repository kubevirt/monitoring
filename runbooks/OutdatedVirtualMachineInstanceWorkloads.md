# OutdatedVirtualMachineInstanceWorkloads

## Meaning

There are VMIs running in outdated virt-launcher pods 24 hours after the KubeVirt control plane update has completed.

## Impact

VMIs that have not been updated to run in the most recent virt-launcher pod may not have access to new KubeVirt features and will not have received any security fixes associated with the virt-launcher pod update.

## Diagnosis

You can identify the VMIs that are out-of-date by using the `kubevirt.io/outdatedLauncherImage` label as a label selector when listing VMIs. Below is an example of a command that will list all out-of-date VMIs across all namespaces within the cluster.

```bash
kubectl get vmi -l “kubevirt.io/outdatedLauncherImage” --all-namespaces
```

## Mitigation

### Check if automated workload updates are enabled

Check the KubeVirt CR used to install KubeVirt to see if the `workloadUpdateStrategy` is configured on the KubeVirt CR's spec.

Below is an example with automated workload updates enabled. With this configuration, workloads which are live migratable will always be live migrated and non live migratable workloads will be evicted which causes a restart when `RunStrategy: Always` is used on the corresponding VM definition

```
apiVersion: kubevirt.io/v1
kind: KubeVirt
metadata:
  name: kubevirt
  namespace: kubevirt
spec:
  imagePullPolicy: IfNotPresent
  workloadUpdateStrategy:
    workloadUpdateMethods:
      - LiveMigrate
      - Evict
    batchEvictSize: 10
    batchEvictInterval: "1m"
```

If automated workload updates are not enabled, consider enabling them. This will cause the KubeVirt control plane to automatically update the VMI workloads using the methods defined in the `workloadUpdateMethods` field.

More information about automated workload updates can be found in our user guide [here](https://kubevirt.io/user-guide/operations/updating_and_deletion/#updating-kubevirt-workloads)

### Automated workload updates are enabled but VMIs are still out-of-date

Identify the VMIs that are out of date. Check each impacted VMI to see if they are live migratable or not by looking for the ```LiveMigratable``` condition within the VMI's status.

If a VMI is not live migratable, and the `evict` method is not chosen as a `workloadUpdateMethods` value on the KubeVirt CR, then the only path forward is to stop the VMI. If the VMI is controlled by a corresponding VM with `runStrategy: always` set, then a new VMI will immediately spin up in an updated virt-launcher pod to replace the stopped VMI. This is the equivalent of a restart action.

If the VMI is live migratable and the migration is failing, then the previous option of stopping the VMI is still possible here, but it will result in the workload being interrupted.

Below is an example of how to stop a VM named `my-vm` in namespace `my-namespace` using `virtctl`

```bash
virctl stop --namespace my-namespace my-vm
```

### Manually updating VMIs

To manually update VMIs, either manually create migration objects for migratable VMIs (non distructive), or manually stop VMIs (distructive) which are not live migratable.

With live migration, the VMIs will be migrated into an updated virt-launcher container.

With stop (and potentially restart when the VMI is controlled by a VM), any replacement VMIs will be started in updated virt-launcher containers. 

Below is an example of how to manually execute a live migration by posting a VirtualMachineInstanceMigration object to the cluster that targets a specific running VM. In this example, the VM is called `my-vm` which is running in the namespace `my-namespace`.

```bash
cat << EOF > migration.yaml
apiVersion: kubevirt.io/v1
kind: VirtualMachineInstanceMigration
metadata:
  name: my-vm-migration-job
  namespace: my-namespace
spec:
  vmiName: my-vm
EOF

kubectl create -f migration.yaml

```

