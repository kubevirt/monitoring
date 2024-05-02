# HPPSharingPoolPathWithOS

## Meaning

This alert fires when the hostpath provisioner (HPP) shares a file system with
other critical components, such as `kubelet` or the operating system (OS).

HPP dynamically provisions hostpath volumes to provide storage for persistent
volume claims (PVCs).

## Impact

A shared hostpath pool puts pressure on the node's disks. The node might have
degraded performance and stability.

## Diagnosis

1. Configure the `HPP_NAMESPACE` environment variable:

   ```bash
   $ export HPP_NAMESPACE="$(kubectl get deployment -A | grep hostpath-provisioner-operator | awk '{print $1}')"
   ```

2. Obtain the status of the `hostpath-provisioner-csi` daemon set pods:

   ```bash
   $ kubectl -n $HPP_NAMESPACE get pods | grep hostpath-provisioner-csi
   ```

3. Check the `hostpath-provisioner-csi` logs to identify the shared pool and path:

   ```bash
   $ kubectl -n $HPP_NAMESPACE logs <csi_daemonset> -c hostpath-provisioner
   ```

   Example output:

   ```text
   I0208 15:21:03.769731       1 utils.go:221] pool (<legacy, csi-data-dir>/csi), shares path with OS which can lead to node disk pressure
   ```

## Mitigation

Using the data obtained in the Diagnosis section, try to prevent the pool path
from being shared with the OS. The specific steps vary based on the node and
other circumstances.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
