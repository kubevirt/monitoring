# HPPNotReady

## Meaning

The hostpath-provisioner (HPP) is used to dynamically provision hostPath volumes to provide storage for PVCs.  
The HPPNotReady alert is indicative of a HPP installation being in a degraded state;  
- Not progressing
- Not available to use

## Impact

HPP is simply not usable, its components are not ready & stopped progressing towards a ready state and hence require action.

## Diagnosis

- Check hostpath-provisioner-operator's pod namespace:
	```bash
	export HPP_NAMESPACE="$(kubectl get deployment -A | grep hostpath-provisioner-operator | awk '{print $1}')"
	```

- Check to see if any of the HPP components are currently not ready.
	```bash
	kubectl -n $HPP_NAMESPACE get all -l k8s-app=hostpath-provisioner
	```
 
- Check the failing corresponding pod logs and describe (substitute <corresponding_pod_name>).
    - `kubectl -n $HPP_NAMESPACE describe pods <corresponding_pod_name>`
    - `kubectl -n $HPP_NAMESPACE logs <corresponding_pod_name>`

## Mitigation

Please open an issue and attach the artifacts gathered in the Diagnosis section.
