# VirtHandlerDaemonSetRolloutFailing

## Meaning

Some virt-handler DaemonSet fail to roll out after 15 minutes. This alert suggests that in the cluster, at least one worker node does not have the virt-handler DaemonSet pod successfully rolled out in the given time. 

## Impact

The severity of the alert is as `warning`. This alert does not mean that the failure of rollouts of all virt-handler DaemonSet. Therefore it should not affect normal VM life-cycle if the cluster is not overloaded. 

## Diagnosis

The diagnosis is to show that at least one worker node does not have a virt-handler pod running. One can follow these steps to identify the nodes associated with the failed rollouts:
- List all the pods in the virt-handler DaemonSet with `kubectl get pods -n kubevirt -l=kubevirt.io=virt-handler`
- For each virt-handler pod, find out the name of the worker node the pod is deployed on with `kubectl -n kubevirt get pod <virt-handler-pod-name> -o jsonpath='{.spec.nodeName}'`

## Mitigation

A common reason for this alert is that the nodes associated with the failed rollouts run out of resources. For such case, one can delete some non-DaemonSet pods from the affected nodes. 

