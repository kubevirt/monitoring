# MultipathPathCritical

## Meaning

At least half the paths for a DM-multipath device, with at least one failed
path, have remained outside the active state for 2 minutes. This includes
total path loss and failure of a device with only one path. The metrics come
from the node-exporter `dmmultipath` collector.

## Impact

The device has severely reduced path redundancy or has already lost all
paths. Applications and VMs using it are at risk of I/O disruption or blocked
I/O. `MultipathPathDegraded` can fire alongside this alert.

## Diagnosis

Use the alert's `instance`, `device`, and `sysfs_name` labels to identify the
node and DM-multipath device. The instance label identifies the node-exporter
scrape target; confirm the corresponding node before inspecting it.

On the affected node, inspect the device and its underlying paths:

```bash
multipath -ll <device>
dmsetup info <device>
dmsetup status <device>
journalctl -u multipathd
```

Check storage fabric connectivity, HBA and switch ports, array controllers,
and iSCSI sessions where applicable. Correlate with storage maintenance and
recent configuration changes.

In Prometheus, compare the failed and total path metrics for the same instance,
device, and sysfs name:

```promql
node_dmmultipath_device_paths_failed{job="node-exporter",instance="<instance>",device="<device>",sysfs_name="<sysfs_name>"}
node_dmmultipath_device_paths{job="node-exporter",instance="<instance>",device="<device>",sysfs_name="<sysfs_name>"}
```

The path metrics classify SCSI `running` and NVMe `live` as active. They report
underlying device state, not multipathd path-checker or ALUA state. Missing
metrics do not establish that storage is healthy.

## Mitigation

Immediately work with the storage administrator to restore connectivity.
Prioritize devices with no usable paths. Investigate fabric, HBA, switch,
storage session, and array controller failures. Verify path recovery and
affected workload I/O before closing the incident. Do not remove a mapped
device while workloads are using it.
