# SSPFailingToReconcile

## Meaning

This alert indicates that while the operator's pod is up, it's reconcile cycle is consistently failing.
It may mean a number of things, including failing to update the resources it is responsible for,
failing to deploy the Template Validator or failing to deploy or update the common-templates.

## Impact

With SSP Operator not reconciling, the dependant components may not deploy at all and/or changes in the components are not reconciled, as a result the common templates and template validator may not be updated or reset in case they fail.

## Diagnosis

- Check the logs of the ssp-operator pod for errors:
    - `export NAMESPACE="$(kubectl get deployment -A | grep ssp-operator | awk '{print $1}')"`
    - `kubectl -n $NAMESPACE describe pods -l control-plane=ssp-operator`
    - `kubectl -n $NAMESPACE logs -l control-plane=ssp-operator`
- Check if the Template Validator is up, and if not check it's logs for errors.
    - `export NAMESPACE="$(kubectl get deployment -A | grep ssp-operator | awk '{print $1}')"`
    - `kubectl -n $NAMESPACE get pods -l name=virt-template-validator`
    - `kubectl -n $NAMESPACE describe pods -l name=virt-template-validator`
    - `kubectl -n $NAMESPACE logs -l name=virt-template-validator`    

 
## Mitigation

Please open an issue and attach the artifacts gathered in the Diagnosis section.