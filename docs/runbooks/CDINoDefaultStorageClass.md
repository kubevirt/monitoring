# CDINoDefaultStorageClass
<!-- Edited by agilboa, 6 Dec 2023 -->

## Meaning

This alert fires when there is no default (Kubernetes or virtualization) storage class, and a data volume is pending for one.

A default virtualization storage class has precedence over a default Kubernetes storage class for creating a VirtualMachine disk image.

## Impact

If there is no default (k8s or virt) storage class, a data volume that requests a default storage class (storage class not explicitly specified) will be pending for one.

## Diagnosis

Get the default Kubernetes storage class:
```bash
$ kubectl get sc -o json | jq '.items[].metadata|select(.annotations."storageclass.kubernetes.io/is-default-class"=="true")|.name'
```

Get the default virtualization storage class:
```bash
$ kubectl get sc -o json | jq '.items[].metadata|select(.annotations."storageclass.kubevirt.io/is-default-virt-class"=="true")|.name'
```

To set the default Kubernetes storage class if needed:
```bash
$ kubectl patch storageclass <storage-class-name> -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
```

To set the default virtualization storage class if needed:
```bash
$ kubectl patch storageclass <storage-class-name> -p '{"metadata": {"annotations":{"storageclass.kubevirt.io/is-default-virt-class":"true"}}}'
```

## Mitigation

Ensure that there is one storage class that has the default (k8s or virt) storage class annotation.

<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->