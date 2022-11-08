<!-- Edited by Jiří Herrmann, 8 Nov 2022 -->

# HPPSharingPoolPathWithOS

## Meaning

This alert fires when the hostpath provisioner (HPP) is sharing a file system with other critical components, such as `kubelet` or the operating system (OS).

HPP dynamically provisions hostpath volumes to provide storage for persistent volume claims (PVCs). In the described scenario, HPP persistent volumes might cause node disk pressure.


## Impact

The affected node might have degraded performance and stability.

## Diagnosis

Check which HPP pool is being shared with the OS by reviewing the DaemonSet logs.

1. Configure the `HPP_NAMESPACE` environment variable:
```bash
$ export HPP_NAMESPACE="$(kubectl get deployment -A | grep hostpath-provisioner-operator | awk '{print $1}')"
```

2. Find the Container Storage Interface (CSI) DaemonSet pods:
```bash
$ kubectl -n $HPP_NAMESPACE get pods | grep hostpath-provisioner-csi
```
 
3. Check the CSI pod logs to learn which pool and path are being shared with the OS:
```bash
$ kubectl -n $HPP_NAMESPACE logs <csi_daemonset_pod> -c hostpath-provisioner
```

The relevant log lines look similar to the following:
```  
I0208 15:21:03.769731       1 utils.go:221] pool (legacy, csi-data-dir/csi), shares path with OS which can lead to node disk pressure
```


## Mitigation

Using the data obtained in the Diagnosis section, try to prevent the pool path from being shared with the OS. The specific steps vary based on the node and other circumstances.

<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->