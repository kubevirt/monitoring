# SSPCommonTemplatesModificationReverted

## Meaning

This alert indicates that one or more common templates were changed by the ssp operator as part of its reconciliation procedure.
If any changes were made to a common template manually or by a script they were reverted by the operator.

## Impact

The changes made manually on the templates are lost.
## Diagnosis

The common templates should never be modified manually, find out if you are misusing the common templates by making such modifications.

- Check the logs of the SSP Operator pods:
	- `export NAMESPACE="$(kubectl get deployment -A | grep ssp-operator | awk '{print $1}')"`
	- `kubectl -n $NAMESPACE logs --tail=-1 -l control-plane=ssp-operator`
	- Look for lines similar to this to find the names of the restored Common Templates:
	  `Changes reverted in common template: <template name>`
 
## Mitigation

Do not modify common templates. If necessary, copy the one you need and modify the copy.
Refer to the [templates documentation](https://kubevirt.io/user-guide/virtual_machines/templates/) to find out more.

