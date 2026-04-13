# LowReadyVirtHandlerCount

## Meaning

This alert fires when one or more `virt-handler` pods are running, but not
all of them have been in a `Ready` state for the last 10 minutes.

The `virt-handler` runs on every node that can schedule VMIs (as a
DaemonSet). Each node has one `virt-handler` pod.

## Impact

Some nodes may have a running but not ready `virt-handler`. VMIs running on those
nodes might not be fully managed (e.g. domain updates, network or storage
changes). If the condition persists, it can lead to the `NoReadyVirtHandler`
alert for affected nodes.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-handler` pods:

   ```bash
   $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-handler -o wide
   ```

3. For pods that are running but not ready, inspect pod conditions and events:

   ```bash
   $ kubectl -n $NAMESPACE describe pod -l kubevirt.io=virt-handler
   ```

4. If pods are in `CrashLoopBackOff` or to inspect runtime failures, check
   non-ready `virt-handler` pod logs and look for errors:

   ```bash
   $ kubectl -n $NAMESPACE logs -l kubevirt.io=virt-handler
   ```

   Note: With multiple pods (DaemonSet), `-l` streams one pod's logs; use a
   pod name from step 2 to target a specific non-ready pod.

5. If needed, check the `virt-handler` DaemonSet and its events:

   ```bash
   $ kubectl -n $NAMESPACE describe daemonset virt-handler
   ```

6. Check for node issues on nodes where `virt-handler` is not ready:

   ```bash
   $ kubectl get nodes
   ```

## Mitigation

Identify why some `virt-handler` pods are not ready (e.g. failed readiness
probe, resource pressure, node issues) and resolve the underlying cause so
all schedulable nodes have a ready `virt-handler`.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
