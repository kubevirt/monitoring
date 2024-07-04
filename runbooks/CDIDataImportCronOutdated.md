# CDIDataImportCronOutdated

## Meaning

This alert fires when `DataImportCron` cannot poll or import the latest disk
image versions.

`DataImportCron` polls disk images, checking for the latest versions, and
imports the images into persistent volume claims (PVCs) or VolumeSnapshots. This
process ensures that these sources are updated to the latest version so that
they can be used as reliable clone sources or golden images for virtual machines
(VMs).

For golden images, _latest_ refers to the latest operating system of the
distribution. For other disk images, _latest_ refers to the latest hash of the
image that is available.

**Note:** If the status of a `DataImportCron` PVC is `Pending` because there is no
default storage class, the `CDIDataImportCronOutdated` alert is suppressed and the
`CDINoDefaultStorageClass` alert is triggered.

## Impact

VMs might be created from outdated disk images.

VMs might fail to start because no boot source is available for cloning.

## Diagnosis

1. Check the cluster for a default Kubernetes storage class:
   ```bash
   $ kubectl get sc -o jsonpath='{.items[?(.metadata.annotations.storageclass\.kubernetes\.io\/is-default-class=="true")].metadata.name}'
   ```

   Check the cluster for a default virtualization storage class:
   ```bash
   $ kubectl get sc -o jsonpath='{.items[?(.metadata.annotations.storageclass\.kubevirt\.io\/is-default-virt-class=="true")].metadata.name}'
   ```

   The output displays the default (Kubernetes and/or virtualization) storage
   class. You must either set a default storage class on the cluster, or ask for
   a specific storage class in the `DataImportCron` specification, in order for
   the `DataImportCron` to poll and import golden images. If the default
   storage class does not exist, the created import DataVolume and PVC will be
   in `Pending` phase.

2. List the `DataImportCron` objects that are not up-to-date:

   ```bash
   $ kubectl get dataimportcron -A -o jsonpath='{range .items[*]}{.status.conditions[?(@.type=="UpToDate")].status}{"\t"}{.metadata.namespace}{"/"}{.metadata.name}{"\n"}{end}' | grep False
   ```

3. If a default storage class is not defined on the cluster, check the
`DataImportCron` specification for a `DataVolume` template storage class:

   ```bash
   $ kubectl -n <namespace> get dataimportcron <dataimportcron> -o jsonpath='{.spec.template.spec.storage.storageClassName}{"\n"}'
   ```

4. Obtain the name of the `DataVolume` associated with the `DataImportCron`
object:

   ```bash
   $ kubectl -n <namespace> get dataimportcron <dataimportcron> -o jsonpath='{.status.lastImportedPVC.name}{"\n"}'
   ```

5. Check the `DataVolume` status:

   ```bash
   $ kubectl -n <namespace> get dv <datavolume> -o jsonpath-as-json='{.status}'
   ```

6. Set the `CDI_NAMESPACE` environment variable:

   ```bash
   $ export CDI_NAMESPACE="$(kubectl get deployment -A -o jsonpath='{.items[?(.metadata.name=="cdi-operator")].metadata.namespace}')"
   ```

7. Check the `cdi-deployment` log for error messages:

   ```bash
   $ kubectl logs -n $CDI_NAMESPACE deployment/cdi-deployment
   ```

## Mitigation

1. Set a default storage class, either on the cluster or in the `DataImportCron`
specification, to poll and import golden images. The updated Containerized Data
Importer (CDI) should resolve the issue within a few seconds.

2. If the issue does not resolve itself, or, if you have changed the default
storage class in the cluster, you must delete the existing boot sources
(data volumes or volume snapshots) in the cluster namespace that are configured
with the previous default storage class. The CDI will recreate the data volumes
with the newly configured default storage class.

3. If your cluster is installed in a restricted network environment, disable the
`enableCommonBootImageImport` feature gate in order to opt out of automatic
updates:

   ```bash
   $ kubectl patch hco kubevirt-hyperconverged -n $CDI_NAMESPACE --type json -p '[{"op": "replace", "path": "/spec/featureGates/enableCommonBootImageImport", "value": false}]'
   ```

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
See the [HCO cluster configuration documentation](https://github.com/kubevirt/hyperconverged-cluster-operator/blob/main/docs/cluster-configuration.md#enablecommonbootimageimport-feature-gate)
for more information.

If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
