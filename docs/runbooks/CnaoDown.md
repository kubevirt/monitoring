# CnaoDown

## Meaning

The `Cluster-network-addons-operator` (CNAO) deploys additional networking components on top of the cluster.
This alert fires when CNAO is down.

## Impact

With CNAO down, the cluster cannot reconcile changes to virtual machine components. As a result, the changes might fail to take effect.

## Diagnosis

- Check CNAO's operator pod namespace:
	```
	$ export NAMESPACE="$(kubectl get deployment -A | grep cluster-network-addons-operator | awk '{print $1}')"
	```

- Check to see if CNAO's operator pod is down:
	```
	$ kubectl -n $NAMESPACE get pods -l name=cluster-network-addons-operator
	```
 
- Check CNAO's operator pod logs and description:
    ```
	$ kubectl -n $NAMESPACE describe pods -l name=cluster-network-addons-operator
    $ kubectl -n $NAMESPACE logs -l name=cluster-network-addons-operator
	```

## Mitigation

Open an issue (TBA where?) and attach the artifacts gathered in the Diagnosis section.
