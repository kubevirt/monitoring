# CnaoNotReady

## Meaning

Cluster-network-addons-operator (CNAO) deploys additional networking components on top of the cluster.
This alert fires when CNAO CR `NetworkAddonsConfig` is not ready, which means that one of the deployed components is not ready.

## Impact

Without CNAO components deployed and ready, not all the kubevirt networking scenarios can be executed.

## Diagnosis

- Check for messages on the `networkaddonsconfig` CR status:
	```
	kubectl get networkaddonsconfig -o custom-columns="":.status.conditions[*].message
	```
- Identify which component is not ready, gather more information from it.
    - Extract the component deployment/daemonSet and namespace from the message. 
      For example in the following message: `DaemonSet "cluster-network-addons/macvtap-cni" update is being processed...`
      the namespace is `cluster-network-addons` and DaemonSet name is `macvtap-cni`
    - Check the components's pod yaml
        - `kubectl -n <component-namespace> get daemonset <daemonset-name> -o yaml`
    - Check the component's pod logs and describe.
        - `kubectl -n <component-namespace> logs <pod-name>`
        - `kubectl -n <component-namespace> describe pod <pod-name>`

## Mitigation

Please open an issue and attach the artifacts gathered in the Diagnosis section.
