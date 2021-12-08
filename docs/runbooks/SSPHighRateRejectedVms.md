# SSPHighRateRejectedVms

## Meaning

This alert indicates that many Virtual Machines were not valid and therefore were not created or not modified.
It means either a user or a script is continuously trying to create invalid Virtual Machines.

## Impact

There is no immediate impact, however one should see what is the reason for the large amount of rejected Virtual Machines.
A Large amount of Virtual Machines may mean that something is not going to work as planned, for example, a script may be consistently failing to create or modify Virtual Machines

## Diagnosis

The validation and rejection of the Virtual Machines is done by the Template Validator, so one may find more information in the logs of the Virtual Machine pods.

- Check the logs of the Template Validator pods:
	- `export NAMESPACE="$(kubectl get deployment -A | grep ssp-operator | awk '{print $1}')"`
	- `kubectl -n $NAMESPACE logs --tail=-1 -l name=virt-template-validator`
	- Look for lines similar to this to find clues for the source of the invalid Virtual Machines:
	  `{"component":"kubevirt-template-validator","level":"info","msg":"evalution summary for ubuntu-3166wmdbbfkroku0:\nminimal-required-memory applied: FAIL, value 1073741824 is lower than minimum [2147483648]\n\nsucceeded=false","pos":"admission.go:25","timestamp":"2021-09-28T17:59:10.934470Z"}`
 
## Mitigation

Stop the source that creates the invalid Virtual Machines.
