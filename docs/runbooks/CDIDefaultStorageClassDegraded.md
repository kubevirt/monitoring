# CDIDefaultStorageClassDegraded
<!-- Edited by agilboa, 6 Dec 2023 -->

## Meaning

This alert fires when the default (Kubernetes or virtualization) storage class supports smart clone (either CSI or snapshot based) and ReadWriteMany.

A default virtualization storage class has precedence over a default Kubernetes storage class for creating a VirtualMachine disk image.

## Impact

If the default storage class does not support smart clone, we fallback to host-assisted cloning, which is the least efficient method of cloning.

If the default storage class does not suppprt ReadWriteMany, a virtual machine using it is not live-migratable.

## Diagnosis

Get the default virtualization storage class:
```bash
$ export CDI_DEFAULT_VIRT_SC="$(kubectl get sc -o json | jq -r '.items[].metadata|select(.annotations."storageclass.kubevirt.io/is-default-virt-class"=="true")|.name')"
$ echo default_virt_sc=$CDI_DEFAULT_VIRT_SC
```

If the default virtualization storage class is set, check if it supports ReadWriteMany
```bash
$ kubectl get storageprofile $CDI_DEFAULT_VIRT_SC -o json | jq '.status.claimPropertySets'| grep ReadWriteMany
```

Otherwise, if the default virtualization storage class is not set, get the default Kubernetes storage class:
```bash
$ export CDI_DEFAULT_K8S_SC="$(kubectl get sc -o json | jq -r '.items[].metadata|select(.annotations."storageclass.kubernetes.io/is-default-class"=="true")|.name')"
$ echo default_k8s_sc=$CDI_DEFAULT_K8S_SC
```

If the default Kubernetes storage class is set, check if it supports ReadWriteMany:
```bash
$ kubectl get storageprofile $CDI_DEFAULT_K8S_SC -o json | jq '.status.claimPropertySets'| grep ReadWriteMany
```

See [doc](https://github.com/kubevirt/containerized-data-importer/blob/main/doc/efficient-cloning.md) for details about smart clone prerequisites.

## Mitigation

Ensure that the default storage class supports smart clone and ReadWriteMany.

<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->