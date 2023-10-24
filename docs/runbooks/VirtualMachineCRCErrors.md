# VirtualMachineCRCErrors
<!-- Edited by machadovilaca, october 2023-->

## Meaning 

Storage Class is incorrectly configured and the usage of a system-wide and
shared dummy page may cause CRC errors when data is written and read across
different processes or threads.

## Impact

Cluster may report a huge number of CRC errors and the cluster might experience
major service outages.

## Diagnosis

Obtain a list of VirtualMachines with an incorrectly configured storage class by
running the following PromQL query:

<!--DS: You can use the Openshift metrics explorer available at 'https://{OPENSHIFT_BASE_URL}/monitoring/query-browser'.-->

```promql
$ kubevirt_ssp_vm_rbd_volume{rxbounce_enabled="false", volume_mode="Block"} == 1

kubevirt_ssp_vm_rbd_volume{name="testvmi-gwgdqp22k7", namespace="test_ns", pv_name="testvmi-gwgdqp22k7", rxbounce_enabled="false", volume_mode="Block"} 1
```

The output displays a list of VirtualMachines that use a storage class without
`rxbounce_enabled`.

Obtain the storage class name by running the following command:

```bash
$ kubectl get pv ${PV_NAME} -o=jsonpath='{.spec.storageClassName}'
```

## Mitigation

Add the "krbd:rxbounce" map option to your storage class configuration, to use
a bounce buffer when receiving data. The default behavior is to read directly
into the destination buffer. A bounce buffer is required if the stability of the
destination buffer cannot be guaranteed.

```bash
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: vm-sc
parameters:
  # ...
  mounter: rbd
  mapOptions: "krbd:rxbounce"
provisioner: openshift-storage.rbd.csi.ceph.com
# ...
```

If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
