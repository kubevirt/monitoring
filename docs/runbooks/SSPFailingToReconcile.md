# SSPFailingToReconcile
<!--apinnick, Nov 2022-->

## Meaning

This alert fires when the reconcile cycle of the Scheduling, Scale and Performance (SSP) Operator fails repeatedly, although the SSP Operator is running.

## Impact

Dependent components might not be deployed. Changes in the components might not be reconciled. As a result, the common templates and/or the Template Validator might not be updated or reset if they fail.

## Diagnosis

1. Export the `NAMESPACE` environment variable:
  ```bash
  $ export NAMESPACE="$(kubectl get deployment -A | grep ssp-operator | awk '{print $1}')"
  ```
2. Obtain the details of the `ssp-operator` pods:
  ```bash
  $ kubectl -n $NAMESPACE describe pods -l control-plane=ssp-operator
  ```
3. Check the `ssp-operator` logs for errors:
  ```bash
  $ kubectl -n $NAMESPACE logs --tail=-1 -l control-plane=ssp-operator
  ```
4. Obtain the status of the `virt-template-validator` pods:
  ```bash
  $ kubectl -n $NAMESPACE get pods -l name=virt-template-validator
  ```
5. Obtain the details of the `virt-template-validator` pods:
  ```bash
  $ kubectl -n $NAMESPACE describe pods -l name=virt-template-validator
  ```
6. Check the `virt-template-validator` logs for errors:
  ```bash
  $ kubectl -n $NAMESPACE logs --tail=-1 -l name=virt-template-validator
  ```
 
## Mitigation

<!--CNV: If you cannot resolve the issue, log in to the [Customer Portal](https://access.redhat.com) and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->

<!--KVstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--KVend-->