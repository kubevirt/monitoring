# CDIDefaultStorageClassDegraded

## Meaning

This alert fires if the default storage class does not support smart cloning
(CSI or snapshot-based) or the ReadWriteMany access mode. The alert does not
fire if at least one default storage class supports these features.

A default virtualization storage class has precedence over a default Kubernetes
storage class for creating a VirtualMachine disk image.

<!--DS: In case of single-node OpenShift, the alert is suppressed if there is a default
storage class that supports smart cloning, but not ReadWriteMany.-->

## Impact

If the default storage class does not support smart cloning, the default cloning
method is host-assisted cloning, which is much less efficient.

If the default storage class does not support ReadWriteMany, virtual machines
(VMs) cannot be live migrated.

## Diagnosis

1. Get the default virtualization storage class by running the following
command:

   ```bash
   $ export CDI_DEFAULT_VIRT_SC="$(kubectl get sc -o jsonpath='{.items[?(.metadata.annotations.storageclass\.kubernetes\.io\/is-default-class=="true")].metadata.name}')"
   ```

2. If a default virtualization storage class exists, check that it supports
ReadWriteMany by running the following command:

   ```bash
   $ kubectl get storageprofile $CDI_DEFAULT_VIRT_SC -o jsonpath='{.status.claimPropertySets}' | grep ReadWriteMany
   ```

3. If there is no default virtualization storage class, get the default
Kubernetes storage class by running the following command:

   ```bash
   $ export CDI_DEFAULT_K8S_SC="$(kubectl get sc -o jsonpath='{.items[?(.metadata.annotations.storageclass\.kubernetes\.io\/is-default-class=="true")].metadata.name}')"
   ```

4. If a default Kubernetes storage class exists, check that it supports
ReadWriteMany by running the following command:

   ```bash
   $ kubectl get storageprofile $CDI_DEFAULT_VIRT_SC -o jsonpath='{.status.claimPropertySets}' | grep ReadWriteMany
   ```

<!--USstart-->
See [doc](https://github.com/kubevirt/containerized-data-importer/blob/main/doc/efficient-cloning.md)
for details about smart clone prerequisites.
<!--USend-->

## Mitigation

Ensure that you have a default (Kubernetes or virtualization) storage class, and
that the default storage class supports smart cloning and ReadWriteMany.

<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->

<!--DS: If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case, attaching
the artifacts gathered during the diagnosis procedure.-->
