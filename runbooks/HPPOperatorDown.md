# HPPOperatorDown
<!-- Edited by Jiří Herrmann, 8 Nov 2022 -->

## Meaning

This alert fires when the hostpath provisioner (HPP) Operator is down.

The HPP Operator deploys and manages the HPP infrastructure components, such as
the daemon set that provisions hostpath volumes.  

## Impact

The HPP components might fail to deploy or to remain in the required state. As a
result, the HPP installation might not work correctly in the cluster.

## Diagnosis

1. Configure the `HPP_NAMESPACE` environment variable:

   ```bash
   $ HPP_NAMESPACE="$(kubectl get deployment -A | grep hostpath-provisioner-operator | awk '{print $1}')"
   ```

2. Check whether the `hostpath-provisioner-operator` pod is currently running:

   ```bash
   $ kubectl -n $HPP_NAMESPACE get pods -l name=hostpath-provisioner-operator
   ```

3. Obtain the details of the `hostpath-provisioner-operator` pod:

   ```bash
   $ kubectl -n $HPP_NAMESPACE describe pods -l name=hostpath-provisioner-operator
   ```

4. Check the log of the `hostpath-provisioner-operator` pod for errors:

   ```bash
   $ kubectl -n $HPP_NAMESPACE logs -l name=hostpath-provisioner-operator
   ```

## Mitigation

Based on the information obtained during Diagnosis, try to find and resolve the
cause of the issue.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
