# CDIOperatorDown

## Meaning

This alert fires when the Containerized Data Importer (CDI) Operator is down.
The CDI Operator deploys and manages the CDI infrastructure components, such as
data volume and persistent volume claim (PVC) controllers. These controllers
help users build virtual machine disks on PVCs.

## Impact

The CDI components might fail to deploy or to stay in a required state. The CDI
installation might not function correctly.

## Diagnosis

1. Set the `CDI_NAMESPACE` environment variable:

   ```bash
   $ export CDI_NAMESPACE="$(kubectl get deployment -A | grep cdi-operator | awk '{print $1}')"
   ```

2. Check whether the `cdi-operator` pod is currently running:

   ```bash
   $ kubectl -n $CDI_NAMESPACE get pods -l name=cdi-operator
   ```

3. Obtain the details of the `cdi-operator` pod:

   ```bash
   $ kubectl -n $CDI_NAMESPACE describe pods -l name=cdi-operator
   ```

4. Check the log of the `cdi-operator` pod for errors:

   ```bash
   $ kubectl -n $CDI_NAMESPACE logs -l name=cdi-operator
   ```

## Mitigation

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the Diagnosis procedure.-->

<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
