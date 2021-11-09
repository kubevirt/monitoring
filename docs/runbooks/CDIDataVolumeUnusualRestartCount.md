# CDIDataVolumeUnusualRestartCount

## Meaning

The `.RestartCount` field for DataVolumes counts the number of times that the CDI ephemeral workload pod was restarted;  
In the case of an HTTP import for example, the `.RestartCount` field will watch the amount of times the CDI importer pod was restarted.  

This alert fires when the restart count of any DataVolume is > 3.

## Impact

On the occasion of reaching more than 3 restarts, it is very unlikely that the DataVolume's operation (build a VM disk on PVC) will converge to success and thus action/investigation needs to be taken.

## Diagnosis

- Find the erronous DataVolume's name & namespace:
	```bash
	kubectl get dv -A -o json | jq -r '.items[] | select(.status.restartCount>3)' | jq '.metadata.name, .metadata.namespace'
	```

- Substitute `<dv_name>`, `<dv_namespace>` with the output from the previous command, to find it's corresponding worker pod (resides in same namespace as DataVolume):
	```bash
	kubectl get pods -n <dv_namespace> -o json | jq -r '.items[] | select(.metadata.ownerReferences[] | select(.name=="<dv_name>")).metadata.name'
	```
 
- Substitute `<worker_pod>` with the above output to check the failing deployments' corresponding pod logs and describe:
    - `kubectl -n <dv_namespace> describe pods <worker_pod>`
    - `kubectl -n <dv_namespace> logs <worker_pod>`

## Mitigation

In some cases, the error could be simply a URL typo for example, in such case you will see a 404 in the artifacts collected above.  
You could then start over by:
- Deleting the DV
- Correcting the URL on the DV manifest
- Creating it again

In other cases, please open an issue and attach the artifacts gathered in the Diagnosis section.
