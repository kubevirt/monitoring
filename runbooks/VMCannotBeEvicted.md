# VMCannotBeEvicted

## Meaning

This alert fires when the eviction strategy of a VM is set to "LiveMigration" but the VM is not migratable.

## Impact

Such a none-migratable VMs block node eviction which can impact operations such as node drain, updates, etc..

## Diagnosis

Check the eviction strategy and the Migratable status of the VMI

- `kubectl get vmis  -o yaml `  
   Search for evictionStrategy field. For example "evictionStrategy: LiveMigrate"

- `kubectl get vmis  -o wide`  
  Look at the "LIVE-MIGRATABLE" column of the output
 
## Mitigation
Change the VM to be migratable or set the evictionStrategy to shutdown
