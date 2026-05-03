# VirtAPIDown

## Meaning

The `virt-api` deployment has fewer than the expected number of
replicas available for 5 minutes.

## Impact

This is a warning level alert. KubeVirt objects may experience degraded
API performance.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-api` pods:

   ```bash
   $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-controller
   ```

3. Check the status of the `virt-api` deployment:

   ```bash
   $ kubectl -n $NAMESPACE get deploy virt-api -o yaml
   ```

4. Check the `virt-api` deployment details for issues such as crashing pods or
image pull failures:

   ```bash
   $ kubectl -n $NAMESPACE describe deploy virt-api
   ```

5. Check for issues such as nodes in a `NotReady` state:

   ```bash
   $ kubectl get nodes
   ```

## Mitigation

Try to identify the root cause and resolve the issue.
<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
