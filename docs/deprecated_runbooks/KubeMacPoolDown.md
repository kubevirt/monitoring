# KubeMacPoolDown [Deprecated]

**Note:** Starting from 4.14, this runbook was replaced by [KubemacpoolDown runbook](http://kubevirt.io/monitoring/runbooks/KubemacpoolDown.html).

## Meaning

`KubeMacPool` is down. `KubeMacPool` is responsible for allocating MAC addresses
and preventing MAC address conflicts.

## Impact

If `KubeMacPool` is down, `VirtualMachine` objects cannot be created.

## Diagnosis

1. Set the `KMP_NAMESPACE` environment variable:

   ```bash
   $ export KMP_NAMESPACE="$(kubectl get pod -A --no-headers -l \
      control-plane=mac-controller-manager | awk '{print $1}')"
   ```

2. Set the `KMP_NAME` environment variable:

   ```bash
   $ export KMP_NAME="$(kubectl get pod -A --no-headers -l \
      control-plane=mac-controller-manager | awk '{print $2}')"
   ```

3. Obtain the `KubeMacPool-manager` pod details:

   ```bash
   $ kubectl describe pod -n $KMP_NAMESPACE $KMP_NAME
   ```

4. Check the `KubeMacPool-manager` logs for error messages:

   ```bash
   $ kubectl logs -n $KMP_NAMESPACE $KMP_NAME
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
