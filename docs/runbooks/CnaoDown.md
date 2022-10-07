# CnaoDown

## Meaning

The `Cluster-network-addons-operator` (CNAO) deploys additional networking components on top of the cluster.
This alert fires when the CNAO is down.

## Impact

If the CNAO is not running, the cluster cannot reconcile changes to virtual machine components. As a result, the changes might fail to take effect.

## Diagnosis

1. Obtain the namespace of the CNAO pod:
	```
	$ export NAMESPACE="$(kubectl get deployment -A | grep cluster-network-addons-operator | awk '{print $1}')"
	```

1. Check the status of the CNAO pod:
	```
	$ kubectl -n $NAMESPACE get pods -l name=cluster-network-addons-operator
	```
 
1. Check the CNAO pod logs:
    ```
    $ kubectl -n $NAMESPACE logs -l name=cluster-network-addons-operator
	```

1. Generate the CNAO pod description:
	```
	$ kubectl -n $NAMESPACE describe pods -l name=cluster-network-addons-operator
	```

## Mitigation

Open an issue (TBA where?) and attach the artifacts gathered in the Diagnosis section.
