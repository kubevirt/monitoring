# MultipathPathDegraded

## Meaning

One or more paths for a DM-multipath device have remained outside the active
state for 5 minutes. The metrics come from the node-exporter `dmmultipath`
collector. When that collector is unavailable or disabled, or no multipath
devices exist, this alert does not fire.

## Impact

Path redundancy is degraded. I/O can continue through usable paths, but
further failures may disrupt access to the device. If all paths are already
failed, applications and VMs using the device may have blocked I/O.

`MultipathPathCritical` can fire alongside this alert.

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

Restore the failed storage paths with the storage administrator. Check cable
and transceiver health, zoning, HBA ports, storage sessions, and array ports.
Verify that the failed path count returns to zero and the alert resolves.
Do not remove the mapped device while workloads are using it.
