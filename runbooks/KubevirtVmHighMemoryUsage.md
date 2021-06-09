# KubevirtVmHighMemoryUsage

## Meaning

Container that hosting Virtual Machine has free memory less than 20 MB and it is close to memory limit.

## Impact

Virtual Machine is at risk of being terminated by the runtime when container memory usage will exceed the memory limit.

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
