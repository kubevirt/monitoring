<!-- Edited by Jiří Herrmann, 8 Nov 2022 -->

# VirtOperatorRESTErrorsBurst 

## Meaning

This alert fires when more than 80% of the REST calls in the `virt-operator` pods failed in the last 5 minutes. This usually indicates that the `virt-operator` pods cannot connect to the API server. 

This error is frequently caused by one of the following problems:

- The API server is overloaded, which causes timeouts. To verify if this is the case, check the metrics of the API server, and view its response times and overall calls.

- The `virt-operator` pod cannot reach the API server. This is commonly caused by DNS issues on the node and networking connectivity issues.

## Impact

Cluster-level actions, such as upgrading and controller reconciliation, probably do not function. 

However, customer workloads, such as virtual machines (VMs) and VM instances (VMIs), are not likely to be affected.

## Diagnosis

1. Set the `NAMESPACE` environment variable:
```
$ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
```

2. Check the status of the `virt-operator` pods:
```
$ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-operator
```

3. Check the `virt-operator` logs for error messages when connecting to the API server:
```
$ kubectl -n $NAMESPACE logs <virt-operator>
```

4. Obtain the details of the `virt-operator` pod:
```
$ kubectl -n $NAMESPACE describe pod <virt-operator>
```

## Mitigation

If the `virt-operator` pod cannot connect to the API server, delete the pod to force a restart:

```
$ kubectl delete -n $NAMESPACE <virt-operator>
```

<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->