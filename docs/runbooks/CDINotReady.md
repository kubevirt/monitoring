# CDINotReady

## Meaning

The CDINotReady alert is indicative of a CDI installation being in a degraded state;  
- Not progressing
- Not available to use

## Impact

CDI is simply not usable - users cannot build Virtual Machine Disks on PVCs using CDI's DataVolumes.  
It's components are not ready & stopped progressing towards a ready state and hence require action.

## Diagnosis

- Check cdi-operator's pod namespace:
	```bash
	export CDI_NAMESPACE="$(kubectl get deployment -A | grep cdi-operator | awk '{print $1}')"
	```

- Check to see if any of the CDI components are currently not ready.
	```bash
	kubectl -n $CDI_NAMESPACE get deploy -l cdi.kubevirt.io
	```
 
- Check the failing deployments' corresponding pod logs and describe (substitute <corresponding_pod_name>).
    - `kubectl -n $CDI_NAMESPACE describe pods <corresponding_pod_name>`
    - `kubectl -n $CDI_NAMESPACE logs <corresponding_pod_name>`

## Mitigation

Please open an issue and attach the artifacts gathered in the Diagnosis section.
