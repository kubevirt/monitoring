<!-- Edited by Jiří Herrmann, 8 Nov 2022 -->

# NoReadyVirtController

## Meaning

This alert fires when no available `virt-controller` devices have been detected for 5 minutes.

A `virt-controller` handles monitoring the custom resource definitions (CRDs) of a virtual machine instance (VMI) and managing the associated pods. The device creates pods for VMIs and manages the lifecycle of the pods.

Therefore, `virt-controller` devices are critical for all cluster-wide virtualization functionality.

## Impact
Any actions related to VM lifecycle management fail. This notably includes launching a new VMI or shutting down an existing VMI.


## Diagnosis

1. Set the `NAMESPACE` environment variable:
```bash
$ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
```

2. Verify a `virt-controller` device is available:
```bash
$ kubectl get deployment -n $NAMESPACE virt-controller -o jsonpath='{.status.readyReplicas}'
```

3. Check the status of the `virt-controller` deployment to find out more information. The following commands provide the associated events and show if any problems occurred, such as crashing pods or failures to pull images:
```bash
$ kubectl -n $NAMESPACE get deploy virt-controller -o yaml
```
```bash
$ kubectl -n $NAMESPACE describe deploy virt-controller
```

4. Obtain the details of the `virt-controller` pods:
```bash
$ get pods -n $NAMESPACE | grep virt-controller
```

5. Inspect the logs for each `virt-controller`:
```bash
$ kubectl logs -n $NAMESPACE <virt-controller>
```

6. Check the nodes for problems, suchs as a `NotReady` state:
```bash
$ kubectl get nodes
```

## Mitigation

Based on the information obtained during Diagnosis, try to find and resolve the cause of the issue.

<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->

