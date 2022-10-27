# KubeMacPoolDown
<!--apinnick, Oct. 2022-->

## Meaning

`KubeMacPool` is down. `KubeMacPool` is responsible for allocating MAC addresses and preventing MAC address conflicts.

## Impact

If `KubeMacPool` is down, `VirtualMachine` objects cannot be created. 

## Diagnosis

1. Set the `KMP_NAMESPACE` environment variable:
  ```bash
  $ export KMP_NAMESPACE="$(kubectl get pod -A --no-headers -l \
    control-plane=mac-controller-manager | awk '{print $1}')"
  ```

1. Set the `KMP_NAME` environment variable:
  ```bash
  $ export KMP_NAME="$(kubectl get pod -A --no-headers -l \
    control-plane=mac-controller-manager | awk '{print $2}')"
  ```

1. Obtain the `KubeMacPool-manager` pod details:
  ```bash
  $ kubectl describe pod -n $KMP_NAMESPACE $KMP_NAME
  ```

1. Check the `KubeMacPool-manager` logs for error messages:
  ```bash
  $ kubectl logs -n $KMP_NAMESPACE $KMP_NAME
  ```

## Mitigation

<!--CNV: If you cannot resolve the issue, log in to the [Customer Portal](https://access.redhat.com) and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->

<!--KVstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--KVend-->


