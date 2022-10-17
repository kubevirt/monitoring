# KubeVirtMultipleSchedulingVMI

## Meaning

The `KubeVirtMultipleSchedulingVMI` alert means that more than 15 Virtual
Machine Instances are in scheduling state for more than 5 minutes. This alert
can suggest an issue in the VirtualMachine's configuration, e.g., a
misconfigured node selector or a problem in the cluster's underlying
infrastructure, e.g. network disruptions,node resource shortage, etc.

## Impact

There is no immediate impact. However, if there are multiple machines in
scheduling state, it might indicate that something is not working as planned,
for example, a script may be consistently creating incorrect VirtualMachines
configurations, or there might be a problem in the cluster's underlying
infrastructure.

## Diagnosis

Set the environment variable NAMESPACE

```bash
export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
```

Check the VirtualMachine's status and conditions, and VM logs and configuration
to find out what is causing the scheduling state.

```bash
kubectl describe vmi <vmi name> -n $NAMESPACE
```

## Mitigation

First, ensure that the VirtualMachine configuration is correct and that there
are no disruptions or resource shortages in the cluster.

This problem can be caused by several reasons. Therefore, we advise you to try
to identify and fix the root cause. If you cannot resolve this issue, please
open an issue and attach the artifacts gathered in the Diagnosis section.
