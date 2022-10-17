# CDINotReady
<!--Edited by davozeni, 17.10.2022-->

## Meaning

The CDINotReady alert indicates that a containerized data importer (CDI) installation is in a degraded state:

- Not progressing

- Not available to use

## Impact

CDI is not usable, so users cannot build Virtual Machine Disks on PVCs using CDI's DataVolumes. 
CDI components are not ready and they stopped progressing towards a ready state.

## Diagnosis

1. Check the pod namespace of `cdi-operator`:
	```bash
	$ export CDI_NAMESPACE="$(kubectl get deployment -A | grep cdi-operator | awk '{print $1}')"
	```

1. Check for CDI components that are currently not ready:
	```bash
	$ kubectl -n $CDI_NAMESPACE get deploy -l cdi.kubevirt.io
	```

1. Check the description of the failing pod:
	```bash
	$ kubectl -n $CDI_NAMESPACE describe pods <failing-pod>
	```

1. Check the logs of the failing pod:
	```bash
	$ kubectl -n $CDI_NAMESPACE logs <failing-pod>
	```

## Mitigation

Open an issue and attach the artifacts gathered in the Diagnosis section.
