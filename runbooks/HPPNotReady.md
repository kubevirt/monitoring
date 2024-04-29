<!-- Edited by Jiří Herrmann, 8 Nov 2022 -->

# HPPNotReady

## Meaning

This alert fires when a hostpath provisioner (HPP) installation is in a degraded
state.

The HPP dynamically provisions hostpath volumes to provide storage for
persistent volume claims (PVCs).  

## Impact

HPP is not usable. Its components are not ready and they are not progressing
towards a ready state.

## Diagnosis

1. Set the `HPP_NAMESPACE` environment variable:

   ```bash
   $ export HPP_NAMESPACE="$(kubectl get deployment -A | grep hostpath-provisioner-operator | awk '{print $1}')"
   ```

2. Check for HPP components that are currently not ready:

   ```bash
   $ kubectl -n $HPP_NAMESPACE get all -l k8s-app=hostpath-provisioner
   ```

3. Obtain the details of the failing pod:

   ```bash
   $ kubectl -n $HPP_NAMESPACE describe pods <pod>
   ```

4. Check the logs of the failing pod:

   ```bash
   $ kubectl -n $HPP_NAMESPACE logs <pod>
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
