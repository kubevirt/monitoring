# VMCannotBeEvicted

## Meaning

This alert fires when the eviction strategy of a VM is set to "LiveMigration" but the VM is not migratable.

## Impact

Non-migratable VMs block node eviction. This can impact operations such as node drain, updates, etc..

## Diagnosis

Check the eviction strategy and the Migratable status of the VMI

- `kubectl get vmis  -o yaml `  
   Search for evictionStrategy field. For example "evictionStrategy: LiveMigrate"

- `kubectl get vmis  -o wide`  
  Look at the "LIVE-MIGRATABLE" column of the output. In case the status is "False" you can inspect the VMI
  to understand what is the reason that the VM can't be migrated.  
  
  Run `kubectl get vmis  -o yaml`  and inspect the `conditions` section under the VMI status. For example:
    
```
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
In order to resolve the situation and alert you can either
1. Set the evictionStrategy to shutdown
1. Inspect the reasons that prohibit the VM to be live migrated (as described above) and see whether this can be changed. For example: changing disk type, network configuration etc.
