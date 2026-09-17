# CnaoDown

## Meaning

This alert fires when the Cluster Network Addons Operator (CNAO) is down.
The CNAO deploys additional networking components on top of the cluster.

## Impact

If the CNAO is not running, the cluster cannot reconcile changes to virtual
machine components. As a result, the changes might fail to take effect.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get deployment -A | grep cluster-network-addons-operator | awk '{print $1}')"
   ```

2. Check the status of the `cluster-network-addons-operator` pod:

   ```bash
   $ kubectl -n $NAMESPACE get pods -l name=cluster-network-addons-operator
   ```

3. Check the `cluster-network-addons-operator` logs for error messages:

   ```bash
   $ kubectl -n $NAMESPACE logs -l name=cluster-network-addons-operator
   ```

4. Obtain the details of the `cluster-network-addons-operator` pods:

   ```bash
   $ kubectl -n $NAMESPACE describe pods -l name=cluster-network-addons-operator
   ```

## Mitigation

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
