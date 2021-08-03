# OrphanedVirtualMachineImages

## Meaning

This alert fires when a VMI (virt-launcher-* pod) is running on a node that does not have a running `virt-handler-*` pod.  

## Impact

When a node does not have a running `virt-handler` any VMI running on that node is considered orphaned. VMIs being orphaned means that they are no longer manageable.

## Diagnosis

You can confirm the alert by finding which nodes your virt-handler pods are running on.

```bash
kubectl get pods --all-namespaces -o wide -l kubevirt.io=virt-handler
```
*Output:*
```
NAME                 READY   STATUS    RESTARTS   AGE  IP               NODE     NOMINATED NODE   READINESS GATES
virt-handler-vhqsp   1/1     Running   0          4h   10.244.140.80    node02   <none>           <none>
virt-handler-xd8jc   1/1     Running   0          4h   10.244.196.168   node01   <none>           <none>
```

Then check to see which nodes the VMIs are running on. If they are running on a node that a virt-handler pod does not exist on - those VMIs are orphaned.

```bash
kubectl get vmis --all-namespaces
```

*Output:*
```
NAMESPACE   NAME            AGE   PHASE       IP    NODENAME
default     vmi-ephemeral   4s    Scheduled         node02
```

## Mitigation

Check to see if the DaemonSet that controls the `virt-handler-*` pods is healthy.

```bash
kubectl get daemonset virt-handler --all-namespaces
```

*Output:*
```
NAME                    DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR            AGE
virt-handler            2         2         2       2            2           kubernetes.io/os=linux   4h
```

The DaemonSet is considered healthy if Desired, Ready and Available are all the same number.  


### Unhealthy virt-handler DaemonSet
Check the DaemonSets status to see what the issues are with deploying the pods.

`kubectl describe daemonset virt-handler --all-namespaces`


Another way to check the status is to get the object and read through its status object.  

`kubectl get daemonset virt-handler --all-namespaces -o yaml | jq .status`


 You should also check the health of the clusters nodes.

`kubectl get nodes`

### Healthy virt-handler DaemonSet

Check to see if there is a `workloads` placement policy on the kubevirt resource. This will be under `spec.workloads`.

`kubectl get kubevirt kubevirt --all-namespaces -o yaml`

If there is a placement policy you can make adjustments so that the node that is running the VMI is included in the placement policy.

It could also be that there has been a change to a node's taints/tolerations or a pod's scheduling rules. You can read more about [scheduling policies](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/) to find why the virt-handler pod was removed from the node.
