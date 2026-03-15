# LowReadyVirtAPICount

## Meaning

This alert fires when one or more `virt-api` pods are running, but not
all of them have been in a `Ready` state for the last 10 minutes.

The `virt-api` serves the KubeVirt API. The deployment typically runs two
replicas for high-availability.

## Impact

Reduced capacity or redundancy for the KubeVirt API. If the condition
persists, it can lead to the `NoReadyVirtAPI` alert and API unavailability.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o jsonpath='{.items[].metadata.namespace}')"
   ```

2. Check the status of the `virt-api` pods:

   ```bash
   $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-api
   ```

3. Check the `virt-api` deployment and its events:

   ```bash
   $ kubectl -n $NAMESPACE describe deploy virt-api
   ```

4. Check pod readiness and conditions for non-ready pods:

   ```bash
   $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-api -o wide
   $ kubectl -n $NAMESPACE describe pod -l kubevirt.io=virt-api
   ```

5. If pods are in `CrashLoopBackOff` or to inspect runtime failures, check
   `virt-api` pod logs and look for errors:

   ```bash
   $ kubectl -n $NAMESPACE logs -l kubevirt.io=virt-api
   ```

6. Check for node issues, such as a `NotReady` state:

   ```bash
   $ kubectl get nodes
   ```

## Mitigation

Identify why some `virt-api` pods are not ready (e.g. failed readiness probe,
resource pressure, image pull issues) and resolve the underlying cause.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
