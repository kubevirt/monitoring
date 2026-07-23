# KubeVirtVMStuckInStartingState [Deprecated]

This alert has been deprecated; it does not indicate a genuine issue. If
triggered, it may be safely ignored and silenced.

## Meaning

The `KubeVirtVMStuckInStartingState` alert means that a VirtualMachine has been
in a starting state for more than 5 minutes. This alert can suggest an issue in
the VirtualMachine configuration, e.g., a misconfigured priority class or a
missing network device.

## Impact

There is no immediate impact. However, if there are multiple machines in a
starting state, it might indicate that something is not working as planned, for
example, a script may be consistently creating incorrect Virtual Machines
configurations.

## Diagnosis

Check the VirtualMachine's status and conditions, and VM logs and configuration
to find out what is causing the starting state.

```bash
$ kubectl describe vmi testvmi-ldgrw -n kubevirt-test-default1

Name:         testvmi-ldgrw
Namespace:    kubevirt-test-default1
Labels:       name=testvmi-ldgrw
Annotations:  kubevirt.io/latest-observed-api-version: v1
              kubevirt.io/storage-observed-api-version: v1alpha3
API Version:  kubevirt.io/v1
Kind:         VirtualMachineInstance
Metadata:
  ...
Spec:
  ...
  Networks:
    Name:  default
    Pod:
  Priority Class Name:               non-preemtible
  Termination Grace Period Seconds:  0
Status:
  Conditions:
    Last Probe Time:       2022-10-03T11:08:30Z
    Last Transition Time:  2022-10-03T11:08:30Z
    Message:               virt-launcher pod has not yet been scheduled
    Reason:                PodNotExists
    Status:                False
    Type:                  Ready
    Last Probe Time:       <nil>
    Last Transition Time:  2022-10-03T11:08:30Z
    Message:               failed to create virtual machine pod: pods "virt-launcher-testvmi-ldgrw-" is forbidden: no PriorityClass with name non-preemtible was found
    Reason:                FailedCreate
    Status:                False
    Type:                  Synchronized
  Guest OS Info:
  Phase:  Pending
  Phase Transition Timestamps:
    Phase:                        Pending
    Phase Transition Timestamp:   2022-10-03T11:08:30Z
  Runtime User:                   0
  Virtual Machine Revision Name:  revision-start-vm-6f01a94b-3260-4c5a-bbe5-dc98d13e6bea-1
Events:
  Type     Reason        Age                From                       Message
  ----     ------        ----               ----                       -------
  Warning  FailedCreate  8s (x13 over 28s)  virtualmachine-controller  Error creating pod: pods "virt-launcher-testvmi-ldgrw-" is forbidden: no PriorityClass with name non-preemtible was found
```

## Mitigation

First, ensure that the VirtualMachine configuration is correct and all necessary
resources exist. For example, if a network device is missing, it should be
created.

If the state of the VirtualMachine is "Pending", it means that it wasn't
scheduled yet, which in turn rules out scheduling issues as the root cause. If
this is the case, possible causes include:

1. virt-launcher pod isn't scheduled
2. Topology hints for VMI aren't updated
3. DV is not provisioned/ready

This problem can be caused by several reasons. Therefore, we advise you to try
to identify and fix the root cause. If you cannot resolve this issue, please
open an issue and attach the artifacts gathered in the Diagnosis section.