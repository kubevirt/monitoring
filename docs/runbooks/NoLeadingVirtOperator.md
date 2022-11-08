# NoLeadingVirtOperator 

## Meaning

This alert fires when the system detects that in the past 10 minutes, no `virt-operator` pod has been holding the leader lease, although one or more virt-operator pods are in the `Ready` state. The alert suggests that there is no operating `virt-operator` pod. 

The `virt-operator` is the first operator to start in a cluster. Its primary responsibilities include the following: 

- Installing, live updating, and live upgrading a cluster

- Monitoring the life cycle of top-level controllers, such as `virt-controller`, `virt-handler`, `virt-launcher`, and managing their reconciliation

- Certain cluster-wide tasks, such as certificate rotation and infrastructure management

The `virt-operator` deployment has a default replica of 2 pods, with one leader pod holding a leader lease, indicating an operating `virt-operator` pod. 

## Impact

This alert indicates a failure at the level of the cluster. As a result, critical cluster-wide management functionalities, such as certification rotation, upgrade, and reconciliation of controllers, might be currently not available.

## Diagnosis


1. Configure the `NAMESPACE` environment variable:
```bash
$ `export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"`
```

2. Obtain the details of the `virt-operator` pods:
```bash
$ `kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-operator`
```

3. Inspect `virt-operator` pod logs:
```bash
$ `kubectl -n $NAMESPACE logs <pod>
```

Log messages that contain `Started leading` or `acquire leader` indicate the leader status of a given `virt-operator` pod.

***Leader pod example:***
```bash
$ kubectl -n $NAMESPACE logs <pod-name> |grep lead
{"component":"virt-operator","level":"info","msg":"Attempting to acquire leader status","pos":"application.go:400","timestamp":"2021-11-30T12:15:18.635387Z"}
I1130 12:15:18.635452       1 leaderelection.go:243] attempting to acquire leader lease <namespace>/virt-operator...
I1130 12:15:19.216582       1 leaderelection.go:253] successfully acquired lease <namespace>/virt-operator
{"component":"virt-operator","level":"info","msg":"Started leading","pos":"application.go:385","timestamp":"2021-11-30T12:15:19.216836Z"}
```
***Non-leader pod example:***
```bash
$ kubectl -n $NAMESPACE logs <pod-name> |grep lead
{"component":"virt-operator","level":"info","msg":"Attempting to acquire leader status","pos":"application.go:400","timestamp":"2021-11-30T12:15:20.533696Z"}
I1130 12:15:20.533792       1 leaderelection.go:243] attempting to acquire leader lease <namespace>/virt-operator...
```

4. Obtain the details of the affected `virt-operator` pods:
```bash
$ `kubectl -n $NAMESPACE describe pod <virt-operator>
```

## Mitigation

Based on the information obtained during Diagnosis, try to find and resolve the cause of the issue.

<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->