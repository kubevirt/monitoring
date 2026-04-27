# VirtControllerHighMemory

## Meaning

This alert fires when the `virt-api` pod memory usage exceeds 90% of
its memory limit for more than 15 minutes.

## Impact

High memory pressure on `virt-controller` may cause OOM kills and disrupt
virtual machine lifecycle operations such as scheduling, migration, and
deletion.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ``` bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the memory usage of the `virt-controller` pod:

   ```bash
   $ kubectl -n $NAMESPACE top pod -l kubevirt.io=virt-controller
   ```

3. Check the resource limits and requests of the `virt-controller` pod:

   ```bash
   $ kubectl -n $NAMESPACE describe pod -l kubevirt.io=virt-controller
   ```

## Mitigation

Increase the memory limit of the `virt-controller` deployment or investigate
the root cause of the high memory usage by checking the logs:

```bash
$ kubectl -n $NAMESPACE logs -l kubevirt.io=virt-controller
```
