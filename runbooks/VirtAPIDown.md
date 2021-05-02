# VirtAPIDown

## Meaning

This alert fires when all KubeVirt API servers are down.

## Impact

Without KubeVirt API servers, no API call around KubeVirt entities can be made anymore.

## Diagnosis

- Set the environment variable `NAMESPACE`
	```
	export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
	```

- Check to see if there are any running virt-api pods.
	```
	kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-api
	```

## Mitigation

There can be several reasons for virt-api pods to be down, identify the root cause and fix it.

- Check the status of the virt-api deployment to find out more information. The following commands will provide the associated events and show if there are any issues with pulling an image, crashing pod, etc. 
	- `kubectl -n $NAMESPACE get deploy virt-api -o yaml`
    - `kubectl -n $NAMESPACE describe deploy virt-api`
- Check if there are issues with the nodes. For example, if they are in a NotReady state.
	- `kubectl get nodes`
