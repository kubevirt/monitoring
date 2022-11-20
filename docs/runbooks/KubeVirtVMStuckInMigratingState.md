# KubeVirtVMStuckInMigratingState
<!-- Edited by apinnick, Nov 2022 -->

## Meaning

This alert fires when a virtual machine (VM) is in a migrating state for more than 5 minutes.

This alert might indicate a problem in the cluster infrastructure, such as network disruptions or insufficient node resources.

## Impact

There is no immediate impact. However, if this alert persists, you should try to resolve the issue.

## Diagnosis

1. Check the node resources:
```bash
$ kubectl get nodes -l node-role.kubernetes.io/worker= -o json | jq '.items | .[].status.allocatable'
```
Example output:
```
{
  "cpu": "5",
  "devices.kubevirt.io/kvm": "1k",
  "devices.kubevirt.io/sev": "0",
  "devices.kubevirt.io/tun": "1k",
  "devices.kubevirt.io/vhost-net": "1k",
  "ephemeral-storage": "33812468066",
  "hugepages-1Gi": "0",
  "hugepages-2Mi": "128Mi",
  "memory": "3783496Ki",
  "pods": "110"
}
```

2. Check the node status conditions:
```bash
$ kubectl get nodes -l node-role.kubernetes.io/worker= -o json | jq '.items | .[].status.conditions'
```
Example output:
```
[
  {
    "lastHeartbeatTime": "2022-10-03T11:13:34Z",
    "lastTransitionTime": "2022-10-03T10:14:20Z",
    "message": "kubelet has sufficient memory available",
    "reason": "KubeletHasSufficientMemory",
    "status": "False",
    "type": "MemoryPressure"
  },
  {
    "lastHeartbeatTime": "2022-10-03T11:13:34Z",
    "lastTransitionTime": "2022-10-03T10:14:20Z",
    "message": "kubelet has no disk pressure",
    "reason": "KubeletHasNoDiskPressure",
    "status": "False",
    "type": "DiskPressure"
  },
  {
    "lastHeartbeatTime": "2022-10-03T11:13:34Z",
    "lastTransitionTime": "2022-10-03T10:14:20Z",
    "message": "kubelet has sufficient PID available",
    "reason": "KubeletHasSufficientPID",
    "status": "False",
    "type": "PIDPressure"
  },
  {
    "lastHeartbeatTime": "2022-10-03T11:13:34Z",
    "lastTransitionTime": "2022-10-03T10:14:30Z",
    "message": "kubelet is posting ready status",
    "reason": "KubeletReady",
    "status": "True",
    "type": "Ready"
  }
]
```

## Mitigation

Check the migration configuration of the virtual machine to ensure that it is appropriate for the workload. 

You set a cluster-wide migration configuration by editing the `MigrationConfiguration` stanza of the `KubeVirt` custom resource.

<!--DS: You set a migration configuration for a specific scope by creating a migration policy.-->
<!--USstart-->
You set a migration configuration for a specific scope by creating a [migration policy](https://kubevirt.io/user-guide/operations/migration_policies/#migration-policies).
<!--USend-->

You can determine whether a VM is bound to a migration policy by viewing its `vm.Status.MigrationState.MigrationPolicyName` parameter.

<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->