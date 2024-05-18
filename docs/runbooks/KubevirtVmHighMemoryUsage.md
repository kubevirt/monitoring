# KubevirtVmHighMemoryUsage

## Meaning

This alert fires when a container hosting a virtual machine (VM) has less than
50 Mi free memory before reaching its requested memory.

## Impact

The virtual machine running inside the container is at risk of eviction by the
runtime if the container's memory request is exceeded.

## Diagnosis

1. Obtain the `virt-launcher` pod details:

   ```bash
   $ kubectl get pod <virt-launcher> -o yaml
   ```

2. Identify `compute` container processes with high memory usage in the
`virt-launcher` pod:

   ```bash
   $ kubectl exec -it <virt-launcher> -c compute -- top
   ```

## Mitigation

Increase the memory request in the `VirtualMachine` specification as in the
following example:

```yaml
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
