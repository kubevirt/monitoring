# CDINotReady
<!--Edited by davozeni, 17.10.2022-->

## Meaning

The CDINotReady alert indicates that a containerized data importer (CDI) installation is in a degraded state:

- Not progressing

- Not available to use

## Impact

CDI is not usable, so users cannot build virtual machine disks on persistent volume claims (PVCs) using CDI's data volumes.
CDI components are not ready and they stopped progressing towards a ready state.

## Diagnosis

1. Set the `CDI_NAMESPACE` environment variable:
	```bash
	$ export CDI_NAMESPACE="$(kubectl get deployment -A | grep cdi-operator | awk '{print $1}')"
	```

1. Check for CDI components that are currently not ready:
	```bash
	$ kubectl -n $CDI_NAMESPACE get deploy -l cdi.kubevirt.io
	```

1. Check the description of the failing pod:
	```bash
	$ kubectl -n $CDI_NAMESPACE describe pods <pod>
	```

1. Check the logs of the failing pod:
	```bash
	$ kubectl -n $CDI_NAMESPACE logs <pod>
	```

## Mitigation

<!--CNV: If you cannot resolve the issue, log in to the [Customer Portal](https://access.redhat.com) and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->

<!--KVstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--KVend-->
