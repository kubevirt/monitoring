# KubeVirtVMStuckInMigratingState [Deprecated]

This alert has been deprecated; it does not indicate a genuine issue. If
triggered, it may be safely ignored and silenced.

## Meaning

The `KubeVirtVMStuckInMigratingState` alert means that a VirtualMachine has been
in a migrating state for more than 5 minutes. This alert can suggest a problem
in the cluster's underlying infrastructure, e.g. network disruptions, node
resource shortage, etc.

## Impact

There is no immediate impact. However, if there are multiple machines in a
migrating state, it might indicate there might be a problem in the cluster's
underlying infrastructure.

## Diagnosis

Check the nodes statuses and conditions for potential issues.

```bash
$ kubectl get nodes -l node-role.kubernetes.io/worker= -o json | jq '.items | .[].status.allocatable'

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

```bash
$ kubectl get nodes -l node-role.kubernetes.io/worker= -o json | jq '.items | .[].status.conditions'

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

Ensure you applied the appropriate migration configuration to the VirtualMachine
according to the nature of the workload.

Migration configurations can either be set globally via Kubevirt CR's
`MigrationConfiguration` struct or to a specific scope with
[Migration Policies](https://kubevirt.io/user-guide/operations/migration_policies/#migration-policies).
To check whether a VirtualMachine is bound to a migration policy, please refer
to its `vm.Status.MigrationState.MigrationPolicyName`.

This problem can be caused by several reasons. Therefore, we advise you to try
to identify and fix the root cause. If you cannot resolve this issue, please
open an issue and attach the artifacts gathered in the Diagnosis section.