# KubevirtVmHighMemoryUsage

## Meaning

The container, that is hosting the Virtual Machine, has less than 20 MB free memory and is close to its memory limit.

## Impact

When the container memory usage will exceed the memory limit, the Virtual Machine will be terminated by the runtime.

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

Consider changing container memory limit.

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
