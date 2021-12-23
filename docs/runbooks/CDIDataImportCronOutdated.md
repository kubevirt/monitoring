# CDIDataImportCronOutdated

## Meaning

DataImportCron is in charge of recurring polling of latest version disk images as PVCs, commonly into a 'golden image' namespace;   
These are PVCs that always get updated with latest version,  
serving as a trustworthy clone source for created VMs (created from latest image of OS).  

This alert fires when a DataImportCron fails to create a corresponding PVC or keep its corresponding PVC updated (new `:latest` exists, we can't update to it).

**Note**:  
In the golden-image use case, `latest` is simply the latest OS version of the distribution.  
In the non-golden-image case, as described above, we are merely referring to the latest hash of the image that is available.

## Impact

VMs are created from outdated disk images *or* VMs fail to start because there is no source PVC to clone from.

## Diagnosis

### Distinguish between golden images and regular crons
DataImportCrons are widely used for golden images and likely serve as a source to create VM disks from, and thus it is vital for them to be up to date.
- Find the erronous DataImportCron's namespace & name:
	```bash
	kubectl get dataimportcron -A -o json | jq -r '.items[] | select(.status.conditions[] | select(.type == "UpToDate" and .status == "False")) | .metadata.namespace + "/" + .metadata.name'
	```
	Output is returned as namespace/name - golden image crons will reside in the `openshift-virtualization-os-images` namespace.

- With golden image crons, verify a default storge class is configured:
	```bash
	kubectl get sc
	```
	A single storage class should be marked with `(default)` next to its name.

### Diagnosis artifacts
- Substitute `<cron_namespace>`, `<cron_name>` to find the DataImportCron's corresponding DataVolume (resides in same namespace as cron):
	```bash
	kubectl -n <cron_namespace> get dataimportcron <cron_name> -o json | jq .status.lastImportedPVC.name
	```
 
- Substitute `<dv_name>` with the above output to check for error messages:
    ```bash
	kubectl -n <cron_namespace> get dv <dv_name> -o yaml
	```

- Find cdi-operator's pod namespace:
	```bash
	export CDI_NAMESPACE="$(kubectl get deployment -A | grep cdi-operator | awk '{print $1}')"
	```
- Check cdi controller logs for error messages:
	```bash
	kubectl logs -n $CDI_NAMESPACE deployment/cdi-deployment
	```

- Follow the instructions in the 'Mitigation' section.

## Mitigation

A common issue when opting in to golden images auto-polling is not having a default storage class set.  
Ensure you have a default storage class set in the cluster *or*  
if you're using a custom DataImportCron, verify there is an explicit storage class set in the DataImportCron definition:
```bash
$ kubectl get dataimportcron cron-test -o yaml | grep -B 5 storageClassName
          url: docker://.../cdi-func-test-tinycore
      storage:
        resources:
          requests:
            storage: 5Gi
        storageClassName: rook-ceph-block
```

In other cases, please open an issue and attach the artifacts gathered in the Diagnosis section.
