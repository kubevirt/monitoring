# HPPOperatorDown

## Meaning

The HPP Operator is responsible for deploying and managing the HPP infrastructure components, such as the daemonset in charge of provisioning hostPath volumes.  
This alert fires when the HPP operator is down.

## Impact

With HPP Operator down, the dependant infrastructure components may not deploy at all or fail to stay in a required, opinionated state.  
As a result, the HPP installation is not 100% operational in the cluster.

## Diagnosis

- Check hostpath-provisioner-operator's pod namespace:
	```bash
	export HPP_NAMESPACE="$(kubectl get deployment -A | grep hostpath-provisioner-operator | awk '{print $1}')"
	```

- Check to see if hostpath-provisioner-operator's pod is currently down.
	```bash
	kubectl -n $HPP_NAMESPACE get pods -l name=hostpath-provisioner-operator
	```
 
- Check hostpath-provisioner-operator's pod logs and describe.
    - `kubectl -n $HPP_NAMESPACE describe pods -l name=hostpath-provisioner-operator`
    - `kubectl -n $HPP_NAMESPACE logs -l name=hostpath-provisioner-operator`

## Mitigation

Please open an issue and attach the artifacts gathered in the Diagnosis section.
