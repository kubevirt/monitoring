<!-- Edited by davozeni, 10 Nov 2022 -->

# CDIStorageProfilesIncomplete

## Meaning

This alert fires when any of the containerized data importer (CDI) storage profiles are incomplete.

When a storage profile is incomplete, CDI cannot automatically infer persistent volume claim (PVC) fields, such as `volumeMode` and  `accessModes`. These PVC fields are needed to successfuly complete a disk request from a user.

## Impact

The creation of a data volume fails.

## Diagnosis

- Find the incomplete storage profile by using the name of the storage class from the associated data volume:
```bash
$ kubectl get storageprofile <storage_class_name>
```

## Mitigation

Refer to the StorageProfile documentation, which explains how you can add the required information to the StorageProfile `Spec` section:
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

<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
