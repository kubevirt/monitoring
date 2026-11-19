# HCOClusterServiceVersionAbnormal

## Meaning

This alert fires when a ClusterServiceVersion (CSV) in the OpenShift
Virtualization install namespace remains in an abnormal phase, such as
`Failed`, `Replacing`, `Pending`, or `Deleting`, for more than 30 minutes.

The CSV is the resource that Operator Lifecycle Manager (OLM) uses to install
and upgrade an operator and its operands, including the Hyperconverged Cluster
Operator (HCO) and the operators it deploys.

This alert is scoped to the namespace where OpenShift Virtualization is
installed. It complements the OLM platform alert `CsvAbnormalOver30Min`, which
lacks the `operator_health_impact` label and therefore does not affect the
OpenShift Virtualization operator health status.

## Impact

OLM cannot complete the installation or upgrade of the affected OpenShift
Virtualization component. Depending on the failing component, parts of the
virtualization management plane might be missing or running an outdated
version. The OpenShift Virtualization Overview shows a `Degraded` status.

## Diagnosis

1. Set the `HCO_NAMESPACE` environment variable:

   ```bash
   $ export HCO_NAMESPACE="$(kubectl get hyperconverged -A --no-headers | awk '{print $1}')"
   ```

2. List the CSVs in the namespace and find the one that matches the `name`
label of the alert:

   ```bash
   $ kubectl -n $HCO_NAMESPACE get csv
   ```

3. Inspect the status of the abnormal CSV, including its `phase`, `reason`,
and conditions:

   ```bash
   $ kubectl -n $HCO_NAMESPACE describe csv <csv-name>
   ```

4. Check the OLM operator log for errors about the CSV or its install plan:

   ```bash
   $ kubectl -n olm logs deployment/olm-operator
   ```

5. If the CSV is stuck while being installed or upgraded, check the related
subscription and install plan:

   ```bash
   $ kubectl -n $HCO_NAMESPACE get subscription
   $ kubectl -n $HCO_NAMESPACE get installplan
   ```

## Mitigation

The mitigation depends on the `reason` reported by the CSV status:

- If the CSV is in the `Failed` phase, check the conditions of the CSV and
the events in the namespace for the underlying cause:

  ```bash
  $ kubectl -n $HCO_NAMESPACE get events --sort-by=.metadata.creationTimestamp
  ```

- If the CSV is stuck in the `Pending` phase, another CSV in the same
namespace might be in a `Failed` state, or the install plan might be waiting
for approval:

  ```bash
  $ kubectl -n $HCO_NAMESPACE patch installplan <installplan-name> --type merge -p '{"spec":{"approved":true}}'
  ```

- If the CSV is stuck in the `Replacing` phase for a long time, the previous
version might not be removable. Check whether the operators owned by the CSV
still have resources, and consider opening a support case with the gathered
artifacts.

<!--DS: If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
