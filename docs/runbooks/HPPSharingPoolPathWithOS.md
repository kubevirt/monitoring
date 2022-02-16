# HPPSharingPoolPathWithOS

## Meaning

The hostpath-provisioner (HPP) is used to dynamically provision hostPath volumes to provide storage for PVCs.  
If HPP is sharing a filesystem with other critical components (kubelet / OS) its possible that HPP PVs will cause node disk pressure.
## Impact

Possible degraded node operation.

## Diagnosis

- Check which HPP pool is being shared with OS by looking at DaemonSet logs
	```bash
	export HPP_NAMESPACE="$(kubectl get deployment -A | grep hostpath-provisioner-operator | awk '{print $1}')"
	```

- Find CSI DaemonSet pods
	```bash
	kubectl -n $HPP_NAMESPACE get pods | grep hostpath-provisioner-csi
	```
 
- Check CSI pod logs to find out which pool & path are being shared with OS;  
The relevant log lines will look like:  
`I0208 15:21:03.769731       1 utils.go:221] pool (legacy, csi-data-dir/csi), shares path with OS which can lead to node disk pressure`  
of course with your chosen pool name/pool path instead.
    ```bash
	kubectl -n $HPP_NAMESPACE logs <csi_daemonset_pod> -c hostpath-provisioner
	```

## Mitigation

Using the data obtained in the diagnosis section, it is possible to prevent the pool path from being shared with OS.  
The specification on how to achieve that varies and is generally handled uniquely per-case according to node constraints.
