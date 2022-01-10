# KubevirtVmHighMemoryUsageBasedOnRequests

## Meaning

The container, that is hosting the Virtual Machine, has less than 20 MB free memory and is close to its requested memory.

## Impact

The scheduler considers the memory request when scheduling a container to a node, then fences off the requested memory on the chosen node for the use of the container.

If a node’s memory is exhausted, the system prioritizes evicting its containers whose memory usage most exceeds their memory request. In serious cases of memory exhaustion, the node OOM killer may select and terminate the Virtual Machine in the container based on a similar metric.

## Diagnosis

- Check compute container memory resource request and limit:
```
kubectl get pod <virt launcher pod name> -o yaml
```
Look for container name: compute

- Try to identify processes with high memory usage:
```
kubectl exec -it <virt launcher pod name> -c compute -- top
```

## Mitigation

Consider changing container memory requested.

Memory resource request and limit are set in VirtualMachine object spec.

Example:
```
spec:
  running: false
  template:
    metadata:
      labels:
        kubevirt.io/vm: vm-name
    spec:
      domain:
        resources:
          limits:
            memory: 200Mi
          requests:
            memory: 128Mi
```
