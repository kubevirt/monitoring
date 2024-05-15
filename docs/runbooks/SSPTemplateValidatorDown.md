# SSPTemplateValidatorDown

## Meaning

This alert fires when all the Template Validator pods are down.

The Template Validator checks virtual machines (VMs) to ensure that they do not
violate their templates.

## Impact

VMs are not validated against their templates. As a result, VMs might be created
with specifications that do not match their respective workloads.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get deployment -A | grep ssp-operator | awk '{print $1}')"
   ```

2. Obtain the status of the `virt-template-validator` pods:

   ```bash
   $ kubectl -n $NAMESPACE get pods -l name=virt-template-validator
   ```

3. Obtain the details of the `virt-template-validator` pods:

   ```bash
   $ kubectl -n $NAMESPACE describe pods -l name=virt-template-validator
   ```

4. Check the  `virt-template-validator` logs for error messages:

   ```bash
   $ kubectl -n $NAMESPACE logs --tail=-1 -l name=virt-template-validator
   ```

## Mitigation

Try to identify the root cause and resolve the issue.
<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
