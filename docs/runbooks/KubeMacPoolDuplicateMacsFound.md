# KubeMacPoolDuplicateMacsFound

## Meaning

KubeMacPool is responsible to allocate MAC Addresses and prevent MAC Address conflicts.  
Upon starting KubeMacPool, it scans the cluster for MAC Addresses in VMs set on the managed namespaces.  
In case it finds duplicates it will raise this alert.

## Impact

Duplicate MAC Addresses on the same LAN may cause unexpected and various connectability issues in the network.

## Diagnosis

- Find the `kubemacpool-mac-controller` namespace and pod  
  `kubectl get pod -A -l control-plane=mac-controller-manager --no-headers -o custom-columns=":metadata.namespace,:metadata.name"`
- Extract the confliciting Mac Addresses from the `kubemacpool-mac-controller` pod  
  `kubectl logs -n <kubemacpool-namespace> <kubemacpool-mac-controller-pod-name> | grep "already allocated"`  
  Those lines have all the information that is needed:
    - The Mac Address which has a conflict.
    - The current entry in the data base (VM namespace, name and interface).
    - The entry that rejected due to conflict (VM namespace, name and interface).  

  Example of a log line:
  `"mac address 02:00:ff:ff:ff:ff already allocated to vm/kubemacpool-test/testvm, br1, conflict with: vm/kubemacpool-test/testvm2, br1"`

## Mitigation

- Inspect the duplicates.
- Update the VMs to remove the duplicate MAC Addresses.
- Once all the conflicts are resolved, restart the `kubemacpool-mac-controller` pod  
  `kubectl delete pod -n <kubemacpool-namespace> <kubemacpool-mac-controller-pod-name>`
