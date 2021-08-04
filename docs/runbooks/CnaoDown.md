# CnaoDown

## Meaning

Cluster-network-addons-operator (CNAO) deploys additional networking components on top of the cluster.
This alert fires when CNAO operator is down.

## Impact

With CNAO operator down, changes in the components are not reconciled, and as a result might not be properly deployed.

## Diagnosis

- Check CNAO's operator pod namespace:
	```
	export NAMESPACE="$(kubectl get deployment -A | grep cluster-network-addons-operator | awk '{print $1}')"
	```

- Check to see if CNAO's operator pod is down.
	```
	kubectl -n $NAMESPACE get pods -l name=cluster-network-addons-operator
	```
 
- Check CNAO's operator pod logs and describe.
    - `kubectl -n $NAMESPACE describe pods -l name=cluster-network-addons-operator`
    - `kubectl -n $NAMESPACE logs -l name=cluster-network-addons-operator`

## Mitigation

Please open an issue and attach the artifacts gathered in the Diagnosis section.
