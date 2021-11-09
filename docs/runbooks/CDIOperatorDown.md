# CDIOperatorDown

## Meaning

The CDI Operator is responsible for deploying and managing the CDI infrastructure components, such as the DataVolume/PVC controllers which help users build Virtual Machine Disks on PVCs.  
This alert fires when the CDI operator is down.

## Impact

With CDI Operator down, the dependant infrastructure components may not deploy at all or fail to stay in a required, opinionated state.  
As a result, the CDI installation is not 100% operational in the cluster.

## Diagnosis

- Check cdi-operator's pod namespace:
	```bash
	export CDI_NAMESPACE="$(kubectl get deployment -A | grep cdi-operator | awk '{print $1}')"
	```

- Check to see if cdi-operator's pod is currently down.
	```bash
	kubectl -n $CDI_NAMESPACE get pods -l name=cdi-operator
	```
 
- Check cdi-operator's pod logs and describe.
    - `kubectl -n $CDI_NAMESPACE describe pods -l name=cdi-operator`
    - `kubectl -n $CDI_NAMESPACE logs -l name=cdi-operator`

## Mitigation

Please open an issue and attach the artifacts gathered in the Diagnosis section.
