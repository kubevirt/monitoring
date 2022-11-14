# CDIStorageProfilesIncomplete
<!-- Edited by davozeni, 10 Nov 2022 -->

## Meaning

This alert fires when a Containerized Data Importer (CDI) storage profile is incomplete.

If a storage profile is incomplete, the CDI cannot infer persistent volume claim (PVC) fields, such as `volumeMode` and  `accessModes`, which are required to create a virtual machine (VM) disk.

## Impact

The CDI cannot create a VM disk on the PVC.

## Diagnosis

- Identify the incomplete storage profile:
```bash
$ kubectl get storageprofile <storage_class>
```

## Mitigation

- Add the missing storage profile information as in the following example:
```bash
$ kubectl patch storageprofile local --type=merge -p '{"spec": {"claimPropertySets": [{"accessModes": ["ReadWriteOnce"], "volumeMode": "Filesystem"}]}}'
```

<!--USstart-->
See [Empty profiles](https://github.com/kubevirt/containerized-data-importer/blob/main/doc/storageprofile.md#empty-storage-profile) and
[User defined profiles](https://github.com/kubevirt/containerized-data-importer/blob/main/doc/storageprofile.md#user-defined-storage-profile) for more details about storage profiles.
<!--USend-->
