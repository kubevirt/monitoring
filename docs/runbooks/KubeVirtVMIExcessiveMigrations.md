# KubeVirtVMIExcessiveMigrations

## Meaning
The `KubeVirtVMIExcessiveMigrations` alert means that a VirtualMachineInstance has been live-migrated more than 12 times over a period of 24 hours.  
This migration rate for a single VMI is significantly higher than a normal operation of KubeVirt, even if taking an upgrade into consideration.  
This alert can suggest an issue in the cluster's underlying infrastructure, e.g. network disruptions, node resource shortage, etc. that causes one VMI or more to be migrated every two hours or less in average in one day.

## Impact
A virtual machine that is being excessively migrated is experiencing a degraded performance, because of the memory page faults that are occurring during the transition from node to node.

## Diagnosis
Check nodes statuses and conditions, Kubelets, and VM logs and configuration to find out what is causing their excessive migrations.

```bash
$ kubectl get nodes -l node-role.kubernetes.io/worker= -o json | jq .items[].status.allocatable
{
  "cpu": "3500m",
  "devices.kubevirt.io/kvm": "1k",
  "devices.kubevirt.io/sev": "0",
  "devices.kubevirt.io/tun": "1k",
  "devices.kubevirt.io/vhost-net": "1k",
  "ephemeral-storage": "38161122446",
  "hugepages-1Gi": "0",
  "hugepages-2Mi": "0",
  "memory": "7000128Ki",
  "pods": "250"
}
...

$ kubectl get nodes -l node-role.kubernetes.io/worker= -o json | jq .items[].status.conditions
[
  {
    "lastHeartbeatTime": "2022-05-26T07:36:01Z",
    "lastTransitionTime": "2022-05-23T08:12:02Z",
    "message": "kubelet has sufficient memory available",
    "reason": "KubeletHasSufficientMemory",
    "status": "False",
    "type": "MemoryPressure"
  },
  {
    "lastHeartbeatTime": "2022-05-26T07:36:01Z",
    "lastTransitionTime": "2022-05-23T08:12:02Z",
    "message": "kubelet has no disk pressure",
    "reason": "KubeletHasNoDiskPressure",
    "status": "False",
    "type": "DiskPressure"
  },
  {
    "lastHeartbeatTime": "2022-05-26T07:36:01Z",
    "lastTransitionTime": "2022-05-23T08:12:02Z",
    "message": "kubelet has sufficient PID available",
    "reason": "KubeletHasSufficientPID",
    "status": "False",
    "type": "PIDPressure"
  },
  {
    "lastHeartbeatTime": "2022-05-26T07:36:01Z",
    "lastTransitionTime": "2022-05-23T08:24:15Z",
    "message": "kubelet is posting ready status",
    "reason": "KubeletReady",
    "status": "True",
    "type": "Ready"
  }
]
...
```

Check kubelet status by logging to the worker node:
```bash
$ systemctl status kubelet
● kubelet.service - Kubernetes Kubelet
   Loaded: loaded (/etc/systemd/system/kubelet.service; enabled; vendor preset: disabled)
  Drop-In: /etc/systemd/system/kubelet.service.d
           └─10-mco-default-madv.conf, 20-logging.conf, 20-nodenet.conf
   Active: active (running) since Sun 2022-03-13 09:25:50 UTC; 2 months 13 days ago
 Main PID: 3892 (kubelet)
    Tasks: 46 (limit: 3355442)
   Memory: 598.0M
      CPU: 2w 1d 15h 38.017s
   CGroup: /system.slice/kubelet.service
           └─3892 kubelet --config=/etc/kubernetes/kubelet.conf --bootstrap-kubeconfig=/etc/kubernetes/kubeconfig [...]
```

Check kubelet journal logs:
```bash
$ journalctl -r -u kubelet
```

## Mitigation
First, ensure that the worker nodes have enough available resources (CPU, memory, disk) to run the VM workloads without interruption.  
If the problem persists, try to identify the root cause and fix it, as it can be caused by several reasons. If you cannot fix it, please open an issue and attach the artifacts gathered in the Diagnosis section.
