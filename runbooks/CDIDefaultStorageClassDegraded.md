# CDIDefaultStorageClassDegraded

## Meaning

This alert fires when there is no default storage class that supports smart
cloning (CSI or snapshot-based) or the ReadWriteMany access mode.

A default virtualization storage class has precedence over a default Kubernetes
storage class for creating a VirtualMachine disk image.

## Impact

If the default storage class does not support smart cloning, the default cloning
method is host-assisted cloning, which is much less efficient.

If the default storage class does not support ReadWriteMany, virtual machines
(VMs) cannot be live migrated.

<!--DS: Note: A default OpenShift Virtualization storage class has precedence
over a default OpenShift Container Platform storage class when creating a
VM disk.-->

## Diagnosis

1. Get the default KubeVirt storage class by running the following command:

   ```bash
   $ export CDI_DEFAULT_VIRT_SC="$(kubectl get sc -o json | jq -r '.items[].metadata|select(.annotations."storageclass.kubevirt.io/is-default-virt-class"=="true")|.name')"
   $ echo default_virt_sc=$CDI_DEFAULT_VIRT_SC
   ```

2. If a default KubeVirt storage class exists, check that it supports
ReadWriteMany by running the following command:

   ```bash
   $ kubectl get storageprofile $CDI_DEFAULT_VIRT_SC -o json | jq '.status.claimPropertySets'| grep ReadWriteMany
   ```

3. If there is no default KubeVirt storage class, get the default Kubernetes
storage class by running the following command:

   ```bash
   $ export CDI_DEFAULT_K8S_SC="$(kubectl get sc -o json | jq -r '.items[].metadata|select(.annotations."storageclass.kubernetes.io/is-default-class"=="true")|.name')"
   $ echo default_k8s_sc=$CDI_DEFAULT_K8S_SC
   ```

4. If a default Kubernetes storage class exists, check that it supports
ReadWriteMany by running the following command:

   ```bash
   $ kubectl get storageprofile $CDI_DEFAULT_K8S_SC -o json | jq '.status.claimPropertySets'| grep ReadWriteMany
   ```

<!--USstart-->
See [doc](https://github.com/kubevirt/containerized-data-importer/blob/main/doc/efficient-cloning.md)
for details about smart clone prerequisites.
<!--USend-->

## Mitigation

Ensure that you have a default storage class, either Kubernetes or KubeVirt, and
that the default storage class supports smart cloning and ReadWriteMany.

<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->

<!--DS: If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case, attaching
the artifacts gathered during the diagnosis procedure.-->
