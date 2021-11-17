# KubeMacPoolDown

## Meaning

KubeMacPool is responsible to allocate MAC Addresses and prevent MAC Address conflicts.  

## Impact

If KubeMacPool-manager pod is down then virtualMachine objects cannot be created.

## Diagnosis

- Find Kubemacpool-manager's pod namespace and name:
	```
	export KMP_NAMESPACE="$(kubectl get pod -A --no-headers -l control-plane=mac-controller-manager | awk '{print $1}')"
 	export KMP_NAME="$(kubectl get pod -A --no-headers -l control-plane=mac-controller-manager | awk '{print $2}')"
	```
- Check Kubemacpool-manager pod's logs and describe.
  - `kubectl describe pod -n $KMP_NAMESPACE $KMP_NAME`
  - `kubectl logs -n $KMP_NAMESPACE $KMP_NAME`

## Mitigation

- Please open an issue and attach the artifacts gathered in the Diagnosis section.
