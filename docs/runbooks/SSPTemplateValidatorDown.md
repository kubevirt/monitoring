# SSPTemplateValidatorDown

## Meaning

The Template Validator is responsible for validating the Virtual Machines do not violate the templates they are assigned to.
This alert fires when all the Template Validator's pods are down.

## Impact

With all Template Validator's pods down the Vms will not be validated against the templates they are associated with.

## Diagnosis

- Check virt-template-validator's operator pod namespace:
	```
	export NAMESPACE="$(kubectl get deployment -A | grep ssp-operator | awk '{print $1}')"
	```

- Check to see if virt-template-validator's pod is currently down.
	```
	kubectl -n $NAMESPACE get pods -l name=virt-template-validator
	```
 
- Check virt-template-validator's pods logs and describe.
    - `kubectl -n $NAMESPACE describe pods -l name=virt-template-validator`
    - `kubectl -n $NAMESPACE logs -l name=virt-template-validator`

## Mitigation

Please open an issue and attach the artifacts gathered in the Diagnosis section.
