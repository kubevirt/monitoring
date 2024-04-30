# NetworkAddonsConfigNotReady

## Meaning

This alert fires when the `NetworkAddonsConfig` custom resource (CR) of the
Cluster Network Addons Operator (CNAO) is not ready.

CNAO deploys additional networking components on the cluster. This alert
indicates that one of the deployed components is not ready.

## Impact

Network functionality is affected.

## Diagnosis

1. Check the status conditions of the `NetworkAddonsConfig` CR to identify the
deployment or daemon set that is not ready:

   ```bash
   $ kubectl get networkaddonsconfig -o custom-columns="":.status.conditions[*].message
   ```

   Example output:

   ```text
   DaemonSet "cluster-network-addons/macvtap-cni" update is being processed...
   ```

2. Check the component's pod for errors:

   ```bash
   $ kubectl -n cluster-network-addons get daemonset <pod> -o yaml
   ```

3. Check the component's logs:

   ```bash
   $ kubectl -n cluster-network-addons logs <pod>
   ```

4. Check the component's details for error conditions:

   ```bash
   $ kubectl -n cluster-network-addons describe <pod>
   ```

## Mitigation

Try to identify the root cause and resolve the issue.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
