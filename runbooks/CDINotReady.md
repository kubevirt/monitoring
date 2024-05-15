# CDINotReady

## Meaning

This alert fires when the Containerized Data Importer (CDI) is in a degraded
state:

- Not progressing
- Not available to use

## Impact

CDI is not usable, so users cannot build virtual machine disks on persistent
volume claims (PVCs) using CDI's data volumes. CDI components are not ready, and
they stopped progressing towards a ready state.

## Diagnosis

1. Set the `CDI_NAMESPACE` environment variable:

   ```bash
   $ export CDI_NAMESPACE="$(kubectl get deployment -A | grep cdi-operator | awk '{print $1}')"
   ```

2. Check the CDI deployment for components that are not ready:

   ```bash
   $ kubectl -n $CDI_NAMESPACE get deploy -l cdi.kubevirt.io
   ```

3. Check the details of the failing pod:

   ```bash
   $ kubectl -n $CDI_NAMESPACE describe pods <pod>
   ```

4. Check the logs of the failing pod:

   ```bash
   $ kubectl -n $CDI_NAMESPACE logs <pod>
   ```

## Mitigation

Try to identify the root cause and resolve the issue.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
