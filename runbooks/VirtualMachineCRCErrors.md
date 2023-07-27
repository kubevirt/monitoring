# VirtualMachineCRCErrors
<!-- Edited by machadovilaca, july 2023-->

## Meaning 

Storage Class is incorrectly configured and the usage of a system-wide and
shared dummy page may cause CRC errors when data is written and read across
different processes or threads.

## Impact

Cluster may report a huge number of CRC errors and the cluster might experience
major service outages.

## Diagnosis

1) Get the volume name from the VM:

```bash
$ kubectl get vm <vm-name> -o jsonpath='{.spec.template.spec.volumes}'
```

2) Get the storage class name from the volume:

```bash
$ kubectl get pvc <volume-name> -o jsonpath='{.spec.storageClassName}'
```

3) Get the storage class configuration:

```bash
$ kubectl get sc <storage-class-name> -o yaml
```

4) Check if the storage class has the "krbd:rxbounce" map option:

## Mitigation

Add the "krbd:rxbounce" map option to your storage class configuration, to use
a bounce buffer when receiving data. The default behavior is to read directly
into the destination buffer. A bounce buffer is needed if the destination buffer
isn't guaranteed to be stable.

```bash
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: vm-sc
parameters:
  ...
  mounter: rbd
  mapOptions: "krbd:rxbounce"
provisioner: openshift-storage.rbd.csi.ceph.com
...
```

If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
