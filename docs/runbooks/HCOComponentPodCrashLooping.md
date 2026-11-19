# HCOComponentPodCrashLooping

## Meaning

This alert fires when a container in a pod that belongs to the OpenShift
Virtualization install namespace has been in `CrashLoopBackOff` for more than
15 minutes.

`CrashLoopBackOff` means that the container keeps exiting and restarting.
The affected pod is part of the OpenShift Virtualization control plane, such
as the Hyperconverged Cluster Operator, `virt-api`, `virt-controller`,
`virt-handler`, or one of the other operators deployed by HCO.

This alert is scoped to the install namespace. It complements the platform
alert `KubePodCrashLooping`, which lacks the `operator_health_impact` label
and therefore does not affect the OpenShift Virtualization operator health
status. Pods in user namespaces, such as virtual machine instances or data
volume importers, are not covered by this alert.

## Impact

The crashing component repeatedly fails to start. Depending on the component,
parts of the virtualization management plane might be degraded or unavailable,
and the OpenShift Virtualization Overview shows a `Degraded` status.

## Diagnosis

1. Identify the crashing pod from the `namespace`, `pod`, and `container`
labels of the alert, then set the variables used by the commands below:

   ```bash
   $ export NAMESPACE="<namespace label>"
   $ export POD="<pod label>"
   $ export CONTAINER="<container label>"
   ```

2. Inspect the status of the pod and the reason for the container restarts:

   ```bash
   $ kubectl -n $NAMESPACE describe pod $POD
   ```

3. Check the log of the crashing container, including the previous instance:

   ```bash
   $ kubectl -n $NAMESPACE logs $POD -c $CONTAINER
   $ kubectl -n $NAMESPACE logs $POD -c $CONTAINER --previous
   ```

4. Check the events in the namespace for related errors:

   ```bash
   $ kubectl -n $NAMESPACE get events --sort-by=.metadata.creationTimestamp
   ```

## Mitigation

The mitigation depends on the root cause found in the diagnosis:

- If the container exits because of a configuration problem, such as an
invalid flag or an unreachable dependency, fix the configuration. For
components deployed by HCO, check the relevant section of the `HyperConverged`
custom resource.

- If the container runs out of memory, check whether the node has enough free
resources, and consider adjusting the memory requests and limits.

- If the container image cannot be pulled or fails to start because of an
image problem, verify the image registry and the image digests in the
deployment.

- As a last resort, delete the pod so that it is recreated:

  ```bash
  $ kubectl -n $NAMESPACE delete pod $POD
  ```

  If the pod keeps crash looping after being recreated, gather the diagnosis
artifacts and continue investigating.

<!--DS: If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
