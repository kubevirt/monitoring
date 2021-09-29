# SSPCommonTemplatesModificationReverted

## Meaning

This alert indicates that one or more common templates were changed by the ssp operator as part of its reconciliation procedure.

## Impact

If any changes were made to a common template manually or by a script they were reverted by the operator.

## Diagnosis

The common templates should never be modified manually, find out if you are misusing the common templates by making such modifications.
 
## Mitigation

Do not modify common templates. If necessary, copy the one you need and modify the copy.
Refer to the [templates documentation](https://kubevirt.io/user-guide/virtual_machines/templates/) to find out more.

