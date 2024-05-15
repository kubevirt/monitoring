# CDIDataVolumeUnusualRestartCount

## Meaning

This alert fires when a `DataVolume` object restarts more than three times.

## Impact

Data volumes are responsible for importing and creating a virtual machine disk
on a persistent volume claim. If a data volume restarts more than three times,
these operations are unlikely to succeed. You must diagnose and resolve the
issue.

## Diagnosis

1. Find Containerized Data Importer (CDI) pods with more than three restarts:

   ```bash
   $ kubectl get pods --all-namespaces -l app=containerized-data-importer -o=jsonpath='{range .items[?(@.status.containerStatuses[0].restartCount>3)]}{.metadata.name}{"/"}{.metadata.namespace}{"\n"}'
   ```

2. Obtain the details of the pods:

   ```bash
   $ kubectl -n <namespace> describe pods <pod>
   ```

3. Check the pod logs for error messages:

   ```bash
   $ kubectl -n <namespace> logs <pod>
   ```

## Mitigation

Delete the data volume, resolve the issue, and create a new data volume.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
