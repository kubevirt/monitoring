# KubeVirtVMStuckInErrorState

## Meaning

The `KubeVirtVMStuckInErrorState` alert means that a VirtualMachine has been in
an error state for more than 5 minutes. VirtualMachines are in error state when
they are in one of the following status:

1. CrashLoopBackOff
2. Unknown
3. Unschedulable
4. ErrImagePull
5. ImagePullBackOff
6. PvcNotFound
7. DataVolumeError

This alert can suggest an issue in the VirtualMachine configuration, e.g. a
missing PVC, or a problem in the cluster's underlying infrastructure, e.g.
network disruptions, node resource shortage, etc.

## Impact

There is no immediate impact. However, if there are multiple machines in an
error state, it might indicate that something is not working as planned, for
example, a script may be consistently creating incorrect VirtualMachines
configurations, or there might be a problem in the cluster's underlying
infrastructure.

## Diagnosis

Check the VirtualMachine's status and conditions, and VM logs and configuration
to find out what is causing the error state.

```bash
$ kubectl describe vmi testvmi-hxghp -n kubevirt-test-default1

Name:         testvmi-hxghp
Namespace:    kubevirt-test-default1
Labels:       name=testvmi-hxghp
Annotations:  kubevirt.io/latest-observed-api-version: v1
              kubevirt.io/storage-observed-api-version: v1alpha3
API Version:  kubevirt.io/v1
Kind:         VirtualMachineInstance
Metadata:
  ...
Spec:
  Domain:
    ...
    Resources:
      Requests:
        Cpu:     5000000Gi
        Memory:  5130000240Mi
  ...
Status:
  Active Pods:
    acbc8143-c1da-45e8-b498-3f0dafcd1383:  
  Conditions:
    Last Probe Time:       2022-10-03T11:11:07Z
    Last Transition Time:  2022-10-03T11:11:07Z
    Message:               Guest VM is not reported as running
    Reason:                GuestNotRunning
    Status:                False
    Type:                  Ready
    Last Probe Time:       <nil>
    Last Transition Time:  2022-10-03T11:11:07Z
    Message:               0/2 nodes are available: 2 Insufficient cpu, 2 Insufficient memory.
    Reason:                Unschedulable
    Status:                False
    Type:                  PodScheduled
  Guest OS Info:
  Phase:  Scheduling
  Phase Transition Timestamps:
    Phase:                        Pending
    Phase Transition Timestamp:   2022-10-03T11:11:07Z
    Phase:                        Scheduling
    Phase Transition Timestamp:   2022-10-03T11:11:07Z
  Qos Class:                      Burstable
  Runtime User:                   0
  Virtual Machine Revision Name:  revision-start-vm-3503e2dc-27c0-46ef-9167-7ae2e7d93e6e-1
Events:
  Type    Reason            Age   From                       Message
  ----    ------            ----  ----                       -------
  Normal  SuccessfulCreate  27s   virtualmachine-controller  Created virtual machine pod virt-launcher-testvmi-hxghp-xh9qn


```

Also, check the nodes statuses and conditions.

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

First, ensure that the VirtualMachine configuration is correct and all necessary
resources exist. For example, if a PVC is missing, it should be created. Also,
verify that the cluster's infrastructure is healthy and there are enough
resources to run the VirtualMachine.

This problem can be caused by several reasons. Therefore, we advise you to try
to identify and fix the root cause. If you cannot resolve this issue, please
open an issue and attach the artifacts gathered in the Diagnosis section.
