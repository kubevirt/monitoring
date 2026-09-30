# HCOMultiArchGoldenImagesDisabled

## Meaning

DataImportCron objects (DIC; also known as golden images) are used to create
boot images for virtual machines (VMs). The images are preloaded in the cluster,
and then used to create VM boot disks with a specific operating system.
By default, the preloaded images use the architecture of the cluster node that
was used to create the image.

If the multi-architecture feature is enabled, multiple preloaded
images are created for each DataImportCronTemplate (DICT), one for each
architecture supported by the cluster and by the original image.

This allows the VMs to be scheduled on nodes with the same architecture as the
preloaded image.

This alert is triggered when running on a heterogeneous cluster (a cluster with
nodes of different architectures) while the multi-architecture feature is disabled.

## Impact

When running on a heterogeneous cluster, if the preloaded image uses different
architecture than the architecture of the node that the VM is scheduled on,
the VM fails to start.

## Diagnosis

HCO checks the workload node architectures in the cluster. By default, HCO
considers the worker nodes as the workload nodes. If the workload node
placement is configured, HCO considers the nodes that match the node selector
in this field as the workload nodes.

> **Note**:
> * In version <!--USstart-->`v1.19.0`<!--USend--><!--DS: v4.23.0-->
  or above, the workload configuration is under `spec.deployment.nodePlacements.workload`
> * In versions earlier than <!--USstart-->`v1.19.0`<!--USend--><!--DS: v4.23.0-->,
    the workload configuration is under `spec.workloads.nodePlacement`

HCO publishes the list of the workload node architectures in the
`status.nodeInfo.workloadsArchitectures` field in the HyperConverged CR.

Read the HyperConverged CR:

```bash
$ kubectl get hyperconverged -n kubevirt-hyperconverged kubevirt-hyperconverged -o yaml
```

* In version <!--USstart-->`v1.19.0`<!--USend--><!--DS: v4.23.0-->
  or above, the result looks similar to this:
  ```yaml
  apiVersion: hco.kubevirt.io/v1
  kind: HyperConverged
  spec:
    ...
    deployment:
      nodePlacements:
        workload: # check if the spec.deployment.nodePlacements.workload field is populated
    ...
  status:
    ...
    nodeInfo:
      workloadsArchitectures:
        - amd64
        - arm64
  ...
  ```

* In versions earlier than <!--USstart-->`v1.19.0`<!--USend--><!--DS: v4.23.0-->, the result
  looks similar to this:
  ```yaml
  apiVersion: hco.kubevirt.io/v1beta1
  kind: HyperConverged
  spec:
    ...
    workloads: # check if the spec.workloads.nodePlacement field is populated
      nodePlacement:
    ...
  status:
    ...
    nodeInfo:
      workloadsArchitectures:
        - amd64
        - arm64
  ...
  ```

## Mitigation

To address this issue, you can either enable the multi-architecture boot image
feature, or modify the workloads node placement in the HyperConverged CR to
include only nodes with a single architecture.

### Enable the multi-architecture boot image feature

The multi-architecture boot image feature is not enabled by default. Enabling
this feature results in the creation of multiple preloaded images for each
DataImportCronTemplate (DICT), one for each architecture supported by the
cluster, and by the original image. However, this feature is not generally
available, and it is not fully supported.

Enablement of the multi-architecture boot image feature:

* Enable the multi-architecture boot image feature
* If the HyperConverged CR contains the `dataImportCronTemplates` field,
and this field is not empty, then you might need to add the
`ssp.kubevirt.io/dict.architectures` annotation to each DICT object in this
field. See the [HCOGoldenImageWithNoArchitectureAnnotation](HCOGoldenImageWithNoArchitectureAnnotation.md)
runbook for more details.

1. open the editor to edit the `HyperConverged` CR:
   ```shell
   $ NAMESPACE="$(kubectl get hyperconverged -A --no-headers | awk '{print $1}')"
   $ kubectl edit hyperconverged -n "${NAMESPACE}" kubevirt-hyperconverged -o yaml
   ```

   The editor opens with the HyperConverged CR YAML.
2. Edit the CR:
   * In version <!--USstart-->`v1.19.0`<!--USend--><!--DS: v4.23.0--> or above, use the `v1` API
     version:
      ```yaml
      apiVersion: hco.kubevirt.io/v1
      kind: HyperConverged
      spec:
        ...
        workloadSources:
          enableMultiArchBootImageImport: true
          dataImportCronTemplates:
            ...
        ...
      ```
   * In versions earlier than <!--USstart-->`v1.19.0`<!--USend--><!--DS: v4.23.0-->, use the
     `v1beta1` API version:
       ```yaml
       apiVersion: hco.kubevirt.io/v1beta1
       kind: HyperConverged
       spec:
         dataImportCronTemplates:
           ...
         ...
         featureGates:
           ...
           enableMultiArchBootImageImport: true
           ...
       ```
3. Save the changes and exit the editor, to apply the changes.

### Modify the Workloads Node Placement

If you do not want to enable the multi-architecture boot image feature, you can
modify the workloads node placement in the HyperConverged CR to include only
nodes with a single architecture.

Below is an example of how to modify the workloads node placement to include
only nodes with the `amd64` architecture, using node affinity:

1. Edit the HyperConverged CR:
    ```bash
    $ kubectl edit hyperconverged -n kubevirt-hyperconverged kubevirt-hyperconverged -o yaml
    ```

    The editor opens with the HyperConverged CR YAML.

2. Edit the CR:
    * In version <!--USstart-->`v1.19.0`<!--USend--><!--DS: v4.23.0--> or above, use the `v1` API
      version:
    ```yaml
    apiVersion: hco.kubevirt.io/v1
    kind: HyperConverged
    spec:
      ...
      deployment:
        nodePlacements:
          workload:
            affinity:
              nodeAffinity:
                requiredDuringSchedulingIgnoredDuringExecution:
                  nodeSelectorTerms:
                    - matchExpressions:
                        - key: kubernetes.io/arch
                          operator: In
                          values:
                            - amd64
      ...
    ```
   * In versions earlier than <!--USstart-->`v1.19.0`<!--USend--><!--DS: v4.23.0-->, use the
     `v1beta1` API version:
    ```yaml
    apiVersion: hco.kubevirt.io/v1beta1
    kind: HyperConverged
    spec:
      ...
      workloads:
        nodePlacement:
          affinity:
            nodeAffinity:
              requiredDuringSchedulingIgnoredDuringExecution:
                nodeSelectorTerms:
                  - matchExpressions:
                      - key: kubernetes.io/arch
                        operator: In
                        values:
                          - amd64
      ...
    ```

3. Save the changes and exit the editor, to apply the changes.

<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->

<!--DS: If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
