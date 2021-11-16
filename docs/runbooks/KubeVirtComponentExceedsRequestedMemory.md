# KubeVirtComponentExceedsRequestedMemory

## Meaning

This alert fires when the amount of memory that is being used by a container is more than what was requested.

## Impact

If this alert consistently fires this could mean that the node's memory resources are not being optimally used and could be overloaded.  

## Diagnosis

- Set the environment variable `NAMESPACE`
	```
	export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
	```

- Check to see what the cpu resource limit is.
	```
	kubectl -n $NAMESPACE get deployment <name-of-resource-firing-alert> -o yaml | grep requests: -A 2
	```

- Check actual resource usage using promQL
  ```  
  container_memory_usage_bytes{namespace="$NAMESPACE",container="<name-of-resource-firing-alert>"}
  ```

## Mitigation

After checking the actual resource usage, determine what a better resource request is for the resource and update it using the `customizeComponents` option on the KubeVirt CR.

```
spec:
  customizeComponents:
    patches:
    - type:
      resourceName: < name-of-resource-firing-alert >
      resourceType: < Deployment|DaemonSet >
      type: strategic
      patch: '{"spec":{"template":{"spec":{"containers":[{"name":"< name-of-resource-firing-alert >","resources":{"requests":{"memory":" < new-memory-request > "}}}]}}}}'
```
