# LowVirtAPICount

## Meaning

This alert is fired when for the whole 60 minutes there is only 1 virt-api pod available, although there are at least 2 nodes available for scheduling.

## Impact

Virt-api pod becomes a single point of failure, can lead to API calls outage in case of eviction.

## Diagnosis

- Set the environment variable NAMESPACE

```
export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
```

- Run
```
kubectl get deployment -n $NAMESPACE virt-api -o jsonpath='{.status.readyReplicas}'
```

## Mitigation

- Check the status of the virt-api deployment to find out more information. The following commands will provide the associated events and show if there are any issues with pulling an image, crashing pod, etc. 
```
kubectl -n $NAMESPACE get deploy virt-api -o yaml
```
```
kubectl -n $NAMESPACE describe deploy virt-api
```
- Check if there are issues with the nodes. For example, if they are in a NotReady state.
```
kubectl get nodes
```