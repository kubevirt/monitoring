# CDIStorageProfilesIncomplete

## Meaning

This alert fires when a Containerized Data Importer (CDI) storage profile is
incomplete.

If a storage profile is incomplete, the CDI cannot infer persistent volume claim
(PVC) fields, such as `volumeMode` and  `accessModes`, which are required to
create a virtual machine (VM) disk.

## Impact

The CDI cannot create a VM disk on the PVC.

## Diagnosis

- Identify the incomplete storage profile:

  ```bash
  $ kubectl get storageprofile <storage_class>
  ```

## Mitigation

- Add the missing storage profile information:

  ```bash
  $ kubectl patch storageprofile local --type=merge -p '{"spec": \
    {"claimPropertySets": [{"accessModes": ["ReadWriteOnce"], \
    "volumeMode": "Filesystem"}]}}'
  ```

<!--USstart-->
See [Empty profiles](https://github.com/kubevirt/containerized-data-importer/blob/main/doc/storageprofile.md#empty-storage-profile)
and [User defined profiles](https://github.com/kubevirt/containerized-data-importer/blob/main/doc/storageprofile.md#user-defined-storage-profile)
for more details about storage profiles.
<!--USend-->

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
