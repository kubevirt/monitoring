# NoReadyVirtAPI

## Meaning

This alert fires when no `virt-api` pod in a `Ready` state has been detected
for 10 minutes.

The `virt-api` serves the KubeVirt API.
Without a ready `virt-api`, API requests for virtual machines
and other KubeVirt resources cannot be served.

## Impact

KubeVirt API is effectively unavailable. Users and controllers cannot perform
API operations such as creating, updating, or deleting virtual machine
instances (VMIs) or other KubeVirt resources.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-api` pods:

   ```bash
   $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-api
   ```

3. Check the `virt-api` deployment and events:

   ```bash
   $ kubectl -n $NAMESPACE describe deploy virt-api
   ```

4. Review logs of any `virt-api` pod that is running but not ready:

   ```bash
   $ kubectl -n $NAMESPACE logs <virt-api-pod-name> --previous
   $ kubectl -n $NAMESPACE logs <virt-api-pod-name>
   ```

5. Check for node issues:

   ```bash
   $ kubectl get nodes
   ```

## Mitigation

Identify the root cause (e.g. all replicas crashing, readiness probe failures,
node or resource issues) and restore at least one ready `virt-api` pod.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
