# SSPOperatorDown

## Meaning

The SSP Operator is responsible for deploying and reconciling the common-templates and the template validator.
This alert fires when the SSP operator is down.

## Impact

With SSP Operator down, the dependant components may not deploy at all and/or changes in the components are not reconciled, as a result the common templates and template validator may not be updated or reset in case they fail.

## Diagnosis

- Check ssp-operator's pod namespace:
	```
	export NAMESPACE="$(kubectl get deployment -A | grep ssp-operator | awk '{print $1}')"
	```

- Check to see if ssp-operator's pod is currently down.
	```
	kubectl -n $NAMESPACE get pods -l control-plane=ssp-operator
	```
 
- Check ssp-operator's pod logs and describe.
    - `kubectl -n $NAMESPACE describe pods -l control-plane=ssp-operator`
    - `kubectl -n $NAMESPACE logs --tail=-1 -l control-plane=ssp-operator`

## Mitigation

Please open an issue and attach the artifacts gathered in the Diagnosis section.
