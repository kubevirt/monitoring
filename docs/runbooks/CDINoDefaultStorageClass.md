# CDINoDefaultStorageClass

## Meaning

This alert fires when there is no default (Kubernetes or KubeVirt) storage
class, and a data volume is pending for one.

A default KubeVirt storage class has precedence over a default Kubernetes
storage class for creating a VirtualMachine disk image.

## Impact

If there is no default Kubernetes or KubeVirt storage class, a data volume that
does not have a specified storage class remains in a "pending" state.

## Diagnosis

1. Check for a default Kubernetes storage class by running the following
command:

```bash
$ kubectl get sc -o json | jq '.items[].metadata|select(.annotations."storageclass.kubernetes.io/is-default-class"=="true")|.name'
```

2. Check for a default KubeVirt storage class by running the following command:

```bash
$ kubectl get sc -o json | jq '.items[].metadata|select(.annotations."storageclass.kubevirt.io/is-default-virt-class"=="true")|.name'
```

## Mitigation

Create a default storage class for either Kubernetes or KubeVirt or for both.

A default KubeVirt storage class has precedence over a default Kubernetes
storage class for creating a virtual machine disk image.

* Create a default Kubernetes storage class by running the following command:

  ```bash
  $ kubectl patch storageclass <storage-class-name> -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
  ```

* Create a default KubeVirt storage class by running the following command:

  ```bash
  $ kubectl patch storageclass <storage-class-name> -p '{"metadata": {"annotations":{"storageclass.kubevirt.io/is-default-virt-class":"true"}}}'
  ```

<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->

<!--DS: If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
