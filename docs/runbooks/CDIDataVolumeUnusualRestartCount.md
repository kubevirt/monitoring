# CDIDataVolumeUnusualRestartCount

## Meaning

This alert fires when a `DataVolume` object restarts more than three times.

The alert fires once for each namespace that has such data volumes, and its
`namespace` label names the namespace of the CDI pod that restarted. For a
cross-namespace clone, that is the namespace of the source PVC, so the target
data volume can be in another namespace. Alerts from older CDI versions don't
have this label.

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

2. Find Containerized Data Importer (CDI) pods with more than three restarts
   in that namespace:

   ```bash
   $ kubectl get pods -n "$NAMESPACE" -l app=containerized-data-importer -o=jsonpath='{range .items[?(@.status.containerStatuses[0].restartCount>3)]}{.metadata.name}{"/"}{.metadata.namespace}{"\n"}'
   ```

   If the alert has no `namespace` label, search all namespaces instead:

   ```bash
   $ kubectl get pods --all-namespaces -l app=containerized-data-importer -o=jsonpath='{range .items[?(@.status.containerStatuses[0].restartCount>3)]}{.metadata.name}{"/"}{.metadata.namespace}{"\n"}'
   ```

3. Obtain the details of the pods. The owner references, labels and annotations
   show which PVC the pod works on, and from there which data volume:

   ```bash
   $ kubectl -n "$NAMESPACE" describe pods <pod>
   ```

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
