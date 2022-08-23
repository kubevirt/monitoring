# KubeVirt components metrics

This document aims to help users that are not familiar with metrics exposed by all the KubeVirt components.
All metrics documented here are auto-generated in each component repository and gathered here.
They reflect and describe exactly what is being exposed.

## Table of Contents

- [kubevirt](#kubevirt)
- [containerized-data-importer](#containerized-data-importer)
- [cluster-network-addons-operator](#cluster-network-addons-operator)
- [ssp-operator](#ssp-operator)
- [hostpath-provisioner-operator](#hostpath-provisioner-operator)
- [hyperconverged-cluster-operator](#hyperconverged-cluster-operator)

<div id='kubevirt'></div>

## [kubevirt](https://github.com/kubevirt/kubevirt/tree/main)

### kubevirt_info
Version information.

### kubevirt_migrate_vmi_data_processed_bytes
The total Guest OS data processed and migrated to the new VM. Type: Gauge.

### kubevirt_migrate_vmi_data_remaining_bytes
The remaining guest OS data to be migrated to the new VM. Type: Gauge.

### kubevirt_migrate_vmi_dirty_memory_rate_bytes
The rate of memory being dirty in the Guest OS. Type: Gauge.

### kubevirt_migrate_vmi_disk_transfer_rate_bytes
The rate at which the disk is being transferred. Type: Gauge.

### kubevirt_migrate_vmi_memory_transfer_rate_bytes
The rate at which the memory is being transferred. Type: Gauge.

### kubevirt_virt_controller_leading
Indication for an operating virt-controller. Type: Gauge.

### kubevirt_virt_controller_ready
Indication for a virt-controller that is ready to take the lead. Type: Gauge.

### kubevirt_vmi_cpu_affinity
The vcpu affinity details. Type: Counter.

### kubevirt_vmi_filesystem_total_bytes
Total VM filesystem capacity in bytes. Type: Gauge.

### kubevirt_vmi_filesystem_used_bytes
Used VM filesystem capacity in bytes. Type: Gauge.

### kubevirt_vmi_memory_actual_balloon_bytes
Current balloon bytes. Type: Gauge.

### kubevirt_vmi_memory_available_bytes
Amount of `usable` memory as seen by the domain. Type: Gauge.

### kubevirt_vmi_memory_domain_total_bytes
The amount of memory in bytes allocated to the domain. The `memory` value in domain xml file. Type: Gauge.

### kubevirt_vmi_memory_pgmajfault
The number of page faults when disk IO was required. Type: Counter.

### kubevirt_vmi_memory_pgminfault
The number of other page faults, when disk IO was not required. Type: Counter.

### kubevirt_vmi_memory_resident_bytes
Resident set size of the process running the domain. Type: Gauge.

### kubevirt_vmi_memory_swap_in_traffic_bytes_total
Swap in memory traffic in bytes. Type: Gauge.

### kubevirt_vmi_memory_swap_out_traffic_bytes_total
Swap out memory traffic in bytes. Type: Gauge.

### kubevirt_vmi_memory_unused_bytes
Amount of `unused` memory as seen by the domain. Type: Gauge.

### kubevirt_vmi_memory_usable_bytes
The amount of memory which can be reclaimed by balloon without causing host swapping in bytes. Type: Gauge.

### kubevirt_vmi_memory_used_bytes
Amount of `used` memory as seen by the domain. Type: Gauge.

### kubevirt_vmi_network_receive_bytes_total
Network traffic receive in bytes. Type: Counter.

### kubevirt_vmi_network_receive_errors_total
Network receive error packets. Type: Counter.

### kubevirt_vmi_network_receive_packets_dropped_total
The number of rx packets dropped on vNIC interfaces. Type: Counter.

### kubevirt_vmi_network_receive_packets_total
Network traffic receive packets. Type: Counter.

### kubevirt_vmi_network_traffic_bytes_total
Deprecated. Type: Counter.

### kubevirt_vmi_network_transmit_bytes_total
Network traffic transmit in bytes. Type: Counter.

### kubevirt_vmi_network_transmit_errors_total
Network transmit error packets. Type: Counter.

### kubevirt_vmi_network_transmit_packets_dropped_total
The number of tx packets dropped on vNIC interfaces. Type: Counter.

### kubevirt_vmi_network_transmit_packets_total
Network traffic transmit packets. Type: Counter.

### kubevirt_vmi_non_evictable
Indication for a VirtualMachine that its eviction strategy is set to Live Migration but is not migratable. Type: Gauge.

### kubevirt_vmi_outdated_count
Indication for the number of VirtualMachineInstance workloads that are not running within the most up-to-date version of the virt-launcher environment. Type: Gauge.

### kubevirt_vmi_phase_count
Sum of VMIs per phase and node.

`phase` can be one of the following: [`Pending`, `Scheduling`, `Scheduled`, `Running`, `Succeeded`, `Failed`, `Unknown`]. Type: Gauge.

### kubevirt_vmi_storage_flush_requests_total
Storage flush requests. Type: Counter.

### kubevirt_vmi_storage_flush_times_ms_total
Total time (ms) spent on cache flushing. Type: Counter.

### kubevirt_vmi_storage_iops_read_total
I/O read operations. Type: Counter.

### kubevirt_vmi_storage_iops_write_total
I/O write operations. Type: Counter.

### kubevirt_vmi_storage_read_times_ms_total
Storage read operation time. Type: Counter.

### kubevirt_vmi_storage_read_traffic_bytes_total
Storage read traffic in bytes. Type: Counter.

### kubevirt_vmi_storage_write_times_ms_total
Storage write operation time. Type: Counter.

### kubevirt_vmi_storage_write_traffic_bytes_total
Storage write traffic in bytes. Type: Counter.

### kubevirt_vmi_vcpu_seconds
Amount of time spent in each state by each vcpu. Where `id` is the vcpu identifier and `state` can be one of the following: [`OFFLINE`, `RUNNING`, `BLOCKED`]. Type: Counter.

### kubevirt_vmi_vcpu_wait_seconds
Amount of time spent by each vcpu while waiting on I/O. Type: Counter.

<div id='containerized-data-importer'></div>

## [containerized-data-importer](https://github.com/kubevirt/containerized-data-importer/tree/main)

### clone_progress
The clone progress in percentage. Type: Counter.
### kubevirt_cdi_clone_dv_unusual_restartcount_total
Total restart count in CDI Data Volume cloner pod. Type: Counter.
### kubevirt_cdi_cr_ready
CDI CR Ready. Type: Gauge.
### kubevirt_cdi_dataimportcron_outdated
DataImportCron has an outdated import. Type: Gauge.
### kubevirt_cdi_dataimportcron_outdated_total
Total count of outdated DataImportCron imports. Type: Counter.
### kubevirt_cdi_import_dv_unusual_restartcount_total
Total restart count in CDI Data Volume importer pod. Type: Counter.
### kubevirt_cdi_incomplete_storageprofiles_total
Total number of incomplete and hence unusable StorageProfile. Type: Gauge.
### kubevirt_cdi_operator_up_total
CDI operator status. Type: Gauge.
### kubevirt_cdi_upload_dv_unusual_restartcount_total
Total restart count in CDI Data Volume upload server pod. Type: Counter.
<div id='cluster-network-addons-operator'></div>

## [cluster-network-addons-operator](https://github.com/kubevirt/cluster-network-addons-operator/tree/main)

### kubevirt_cnao_cr_kubemacpool_deployed
KubeMacpool is deployed by CNAO CR. Type: Gauge.
### kubevirt_cnao_cr_kubemacpool_deployed_total
Total count of KubeMacPool manager pods deployed by CNAO CR. Type: Gauge.
### kubevirt_cnao_cr_ready
CNAO CR Ready. Type: Gauge.
### kubevirt_cnao_kubemacpool_manager_num_up_pods_total
Total count of running KubeMacPool manager pods. Type: Gauge.
### kubevirt_cnao_num_up_operators
Total count of running CNAO operators. Type: Gauge.
### kubevirt_kubemacpool_duplicate_macs_total
Total count of duplicate KubeMacPool MAC addresses. Type: Gauge.
<div id='ssp-operator'></div>

## [ssp-operator](https://github.com/kubevirt/ssp-operator/tree/master)

### kubevirt_ssp_num_of_operator_reconciling_properly
The total number of ssp-operator pods reconciling with no errors. Type: Gauge.
### kubevirt_ssp_operator_up_total
The total number of running ssp-operator pods. Type: Gauge.
### kubevirt_ssp_rejected_vms_total
The total number of vms rejected by virt-template-validator. Type: Counter.
### kubevirt_ssp_template_validator_up_total
The total number of running virt-template-validator pods. Type: Gauge.
### kubevirt_ssp_total_restored_common_templates
The total number of common templates restored by the operator back to their original state. Type: Counter.
<div id='hostpath-provisioner-operator'></div>

## [hostpath-provisioner-operator](https://github.com/kubevirt/hostpath-provisioner-operator/tree/main)

### kubevirt_hpp_operator_up_total
The total number of running hostpath-provisioner-operator pods. Type: Gauge.
<div id='hyperconverged-cluster-operator'></div>

## [hyperconverged-cluster-operator](https://github.com/kubevirt/hyperconverged-cluster-operator/tree/main)

### kubevirt_hco_hyperconverged_cr_exists
Indicates whether the HyperConverged custom resource exists (1) or not (0). Type: Gauge.
### kubevirt_hco_out_of_band_modifications_count
Count of out-of-band modifications overwritten by HCO. Type: Counter.
### kubevirt_hco_unsafe_modification_count
Count of unsafe modifications in the HyperConverged annotations. Type: Gauge.
