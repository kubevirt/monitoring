# CDIStorageProfilesIncomplete

## Meaning

An incomplete StorageProfile means that CDI will not be able to automatically infer PVC fields such as volumeMode, accessMode for the user's disk request.  

This alert fires when any of the StorageProfiles are incomplete, and thus not inferrable.

## Impact

DataVolume's creation operation (build a VM disk on PVC) is not succeeding.

## Diagnosis

- Find the storage profile that couldn't be fully populated by CDI using the name of your desired storage class from the DataVolume:
	```bash
	kubectl get storageprofile <your_storage_class_name>
	```

- Follow the instructions in the 'Mitigation' section to populate the missing info.

## Mitigation

Please refer to the StorageProfile documentation, which states how one can provide the needed information in the StorageProfile spec section:  
[Empty profiles](https://github.com/kubevirt/containerized-data-importer/blob/main/doc/storageprofile.md#empty-storage-profile)  
[User defined profiles](https://github.com/kubevirt/containerized-data-importer/blob/main/doc/storageprofile.md#user-defined-storage-profile)

An example of providing the missing info:
#### Before
```yaml
apiVersion: cdi.kubevirt.io/v1beta1
kind: StorageProfile
metadata:
  name: local
spec: {}
status:
  provisioner: kubernetes.io/no-provisioner
  storageClass: local
```
#### Addition
```bash
kubectl patch storageprofile local --type=merge -p '{"spec": {"claimPropertySets": [{"accessModes": ["ReadWriteOnce"], "volumeMode": "Filesystem"}]}}'
```
#### After
```yaml
apiVersion: cdi.kubevirt.io/v1beta1
kind: StorageProfile
metadata:
  name: local
spec:
  claimPropertySets:
  - accessModes:
    - ReadWriteOnce
    volumeMode: Filesystem
status:
  claimPropertySets:
  - accessModes:
    - ReadWriteOnce
    volumeMode: Filesystem
  provisioner: kubernetes.io/no-provisioner
  storageClass: local
```
