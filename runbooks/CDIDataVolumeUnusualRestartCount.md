# CDIDataVolumeUnusualRestartCount
<!-- Edited by apinnick, Nov 2022-->

## Meaning

This alert fires when a `DataVolume` object restarts more than three times.

## Impact

Data volumes are responsible for importing and creating a virtual machine disk on a persistent volume claim. If a data volume restarts more than three times, these operations are unlikely to succeed. You must diagnose and resolve the issue.

## Diagnosis

1. Obtain the name and namespace of the data volume:

   ```bash
   $ kubectl get dv -A -o json | jq -r '.items[] | select(.status.restartCount>3)' | jq '.metadata.name, .metadata.namespace'
   ```

2. Check the status of the pods associated with the data volume:

   ```bash
   $ kubectl get pods -n <namespace> -o json | jq -r '.items[] | select(.metadata.ownerReferences[] | select(.name=="<dv_name>")).metadata.name'
   ```

3. Obtain the details of the pods:

   ```bash
   $ kubectl -n <namespace> describe pods <pod>
   ```

4. Check the pod logs for error messages:

   ```bash
   $ kubectl -n <namespace> describe logs <pod>
   ```

## Mitigation

Delete the data volume, resolve the issue, and create a new data volume.

<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->

