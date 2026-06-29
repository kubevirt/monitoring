# SSPDown

## Meaning

This alert fires when all the Scheduling, Scale and Performance (SSP) Operator
pods are down.

The SSP Operator is responsible for deploying and reconciling the common
templates and the Template Validator.

## Impact

Dependent components might not be deployed. Changes in the components might not
be reconciled. As a result, the common templates and/or the Template Validator
might not be updated or reset if they fail.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get deployment -A | grep ssp-operator | awk '{print $1}')"
   ```

2. Check the status of the `ssp-operator` pods.

   ```bash
   $ kubectl -n $NAMESPACE get pods -l control-plane=ssp-operator
   ```

3. Obtain the details of the `ssp-operator` pods:

   ```bash
   $ kubectl -n $NAMESPACE describe pods -l control-plane=ssp-operator
   ```

4. Check the `ssp-operator` logs for error messages:

   ```bash
   $ kubectl -n $NAMESPACE logs --tail=-1 -l control-plane=ssp-operator
   ```

## Mitigation

Try to identify the root cause and resolve the issue.
<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
