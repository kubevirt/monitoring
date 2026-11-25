# CDIDataVolumeUnusualRestartCount

## Meaning

This alert fires when a `DataVolume` object restarts more than three times.

The alert fires once for each namespace with such restarts, and its
`namespace` label names the namespace of the CDI pod that restarted. For a
cross-namespace clone, that namespace holds the source PVC and the clone source
pod, so the target data volume can be in another namespace. Alerts from older
CDI versions don't have this label.

## Impact

Data volumes are responsible for importing and creating a virtual machine disk
on a persistent volume claim. If a data volume restarts more than three times,
these operations are unlikely to succeed. You must diagnose and resolve the
issue.

## Diagnosis

1. Set `NAMESPACE` to the value of the alert's `namespace` label:

   ```bash
   $ export NAMESPACE=<namespace>
   ```

2. Find the Containerized Data Importer (CDI) pods with a container that
   restarted more than three times. The alert counts the `importer`,
   `cdi-upload-server` and `cdi-clone-source` containers:

   ```bash
   $ kubectl get pods -n "$NAMESPACE" -l app=containerized-data-importer -o json | jq -r '.items[] | select(any(.status.containerStatuses[]?; (.name == "importer" or .name == "cdi-upload-server" or .name == "cdi-clone-source") and .restartCount > 3)) | .metadata.name'
   ```

   If the alert has no `namespace` label, search all namespaces instead, then
   set `NAMESPACE` to the namespace of the pod you want to look at:

   ```bash
   $ kubectl get pods --all-namespaces -l app=containerized-data-importer -o json | jq -r '.items[] | select(any(.status.containerStatuses[]?; (.name == "importer" or .name == "cdi-upload-server" or .name == "cdi-clone-source") and .restartCount > 3)) | "\(.metadata.name)/\(.metadata.namespace)"'
   $ export NAMESPACE=<namespace of the pod>
   ```

3. Obtain the details of the pods:

   ```bash
   $ kubectl -n "$NAMESPACE" describe pods <pod>
   ```

   To find the data volume, look at the PVC the pod works on. Its owner
   references name the data volume when it has one. For a clone source pod,
   the `cdi.kubevirt.io/storage.ownerRef` annotation holds the namespace and
   name of the target PVC, which can be in another namespace. If that PVC has
   no data volume owner, check the PVC's own events instead.

4. Check the pod logs for error messages:

   ```bash
   $ kubectl -n "$NAMESPACE" logs <pod>
   ```

## Mitigation

Delete the data volume, resolve the issue, and create a new data volume.

<!--DS: If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
