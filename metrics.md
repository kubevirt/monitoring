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
- [hostpath-provisioner](#hostpath-provisioner)
- [hyperconverged-cluster-operator](#hyperconverged-cluster-operator)

<div id='kubevirt'></div>

## [kubevirt](https://github.com/kubevirt/kubevirt/tree/main)

### kubevirt_info
Version information.

### kubevirt_allocatable_nodes
The number of allocatable nodes in the cluster. Type: Gauge.

### kubevirt_api_request_deprecated_total
The total number of requests to deprecated KubeVirt APIs. Type: Counter.

### kubevirt_configuration_emulation_enabled
Indicates whether the Software Emulation is enabled in the configuration. Type: Gauge.

### kubevirt_nodes_with_kvm
The number of nodes in the cluster that have the devices.kubevirt.io/kvm resource available. Type: Gauge.

### kubevirt_number_of_vms
The number of VMs in the cluster by namespace. Type: Gauge.

### kubevirt_virt_api_up
The number of virt-api pods that are up. Type: Gauge.

### kubevirt_virt_controller_leading_status
Indication for an operating virt-controller. Type: Gauge.

### kubevirt_virt_controller_ready
The number of virt-controller pods that are ready. Type: Gauge.

### kubevirt_virt_controller_ready_status
Indication for a virt-controller that is ready to take the lead. Type: Gauge.

### kubevirt_virt_controller_up
The number of virt-controller pods that are up. Type: Gauge.

### kubevirt_virt_handler_up
The number of virt-handler pods that are up. Type: Gauge.

### kubevirt_virt_operator_leading
The number of virt-operator pods that are leading. Type: Gauge.

### kubevirt_virt_operator_leading_status
Indication for an operating virt-operator. Type: Gauge.

### kubevirt_virt_operator_ready
The number of virt-operator pods that are ready. Type: Gauge.

### kubevirt_virt_operator_ready_status
Indication for a virt-operator that is ready to take the lead. Type: Gauge.

### kubevirt_virt_operator_up
The number of virt-operator pods that are up. Type: Gauge.

### kubevirt_vm_container_free_memory_bytes_based_on_rss
The current available memory of the VM containers based on the rss. Type: Gauge.

### kubevirt_vm_container_free_memory_bytes_based_on_working_set_bytes
The current available memory of the VM containers based on the working set. Type: Gauge.

### kubevirt_vm_created_total
Amount of VMs created, broken down by namespace, since install. Type: Counter.

### kubevirt_vm_error_status_last_transition_timestamp_seconds
Virtual Machine last transition timestamp to error status. Type: Counter.

### kubevirt_vm_migrating_status_last_transition_timestamp_seconds
Virtual Machine last transition timestamp to migrating status. Type: Counter.

### kubevirt_vm_non_running_status_last_transition_timestamp_seconds
Virtual Machine last transition timestamp to paused/stopped status. Type: Counter.

### kubevirt_vm_running_status_last_transition_timestamp_seconds
Virtual Machine last transition timestamp to running status. Type: Counter.

### kubevirt_vm_starting_status_last_transition_timestamp_seconds
Virtual Machine last transition timestamp to starting status. Type: Counter.

### kubevirt_vmi_cpu_system_usage_seconds_total
Total CPU time spent in system mode. Type: Counter.

### kubevirt_vmi_cpu_usage_seconds_total
Total CPU time spent in all modes (sum of both vcpu and hypervisor usage). Type: Counter.

### kubevirt_vmi_cpu_user_usage_seconds_total
Total CPU time spent in user mode. Type: Counter.

### kubevirt_vmi_filesystem_capacity_bytes
Total VM filesystem capacity in bytes. Type: Gauge.

### kubevirt_vmi_filesystem_used_bytes
Used VM filesystem capacity in bytes. Type: Gauge.

### kubevirt_vmi_memory_actual_balloon_bytes
Current balloon size in bytes. Type: Gauge.

### kubevirt_vmi_memory_available_bytes
Amount of usable memory as seen by the domain. This value may not be accurate if a balloon driver is in use or if the guest OS does not initialize all assigned pages Type: Gauge.

### kubevirt_vmi_memory_cached_bytes
The amount of memory that is being used to cache I/O and is available to be reclaimed, corresponds to the sum of `Buffers` + `Cached` + `SwapCached` in `/proc/meminfo`. Type: Gauge.

### kubevirt_vmi_memory_domain_bytes
The amount of memory in bytes allocated to the domain. The `memory` value in domain xml file. Type: Gauge.

### kubevirt_vmi_memory_pgmajfault_total
The number of page faults when disk IO was required. Page faults occur when a process makes a valid access to virtual memory that is not available. When servicing the page fault, if disk IO is required, it is considered as major fault. Type: Counter.

### kubevirt_vmi_memory_pgminfault_total
The number of other page faults, when disk IO was not required. Page faults occur when a process makes a valid access to virtual memory that is not available. When servicing the page fault, if disk IO is NOT required, it is considered as minor fault. Type: Counter.

### kubevirt_vmi_memory_resident_bytes
Resident set size of the process running the domain. Type: Gauge.

### kubevirt_vmi_memory_swap_in_traffic_bytes
The total amount of data read from swap space of the guest in bytes. Type: Gauge.

### kubevirt_vmi_memory_swap_out_traffic_bytes
The total amount of memory written out to swap space of the guest in bytes. Type: Gauge.

### kubevirt_vmi_memory_unused_bytes
The amount of memory left completely unused by the system. Memory that is available but used for reclaimable caches should NOT be reported as free. Type: Gauge.

### kubevirt_vmi_memory_usable_bytes
The amount of memory which can be reclaimed by balloon without pushing the guest system to swap, corresponds to 'Available' in /proc/meminfo Type: Gauge.

### kubevirt_vmi_memory_used_bytes
Amount of `used` memory as seen by the domain. Type: Gauge.

### kubevirt_vmi_migration_data_processed_bytes
The total Guest OS data processed and migrated to the new VM. Type: Gauge.

### kubevirt_vmi_migration_data_remaining_bytes
The remaining guest OS data to be migrated to the new VM. Type: Gauge.

### kubevirt_vmi_migration_dirty_memory_rate_bytes
The rate of memory being dirty in the Guest OS. Type: Gauge.

### kubevirt_vmi_migration_disk_transfer_rate_bytes
The rate at which the memory is being transferred. Type: Gauge.

### kubevirt_vmi_migration_failed
Indicates if the VMI migration failed. Type: Gauge.

### kubevirt_vmi_migration_phase_transition_time_from_creation_seconds
Histogram of VM migration phase transitions duration from creation time in seconds. Type: Histogram.

### kubevirt_vmi_migration_succeeded
Indicates if the VMI migration succeeded. Type: Gauge.

### kubevirt_vmi_migrations_in_pending_phase
Number of current pending migrations. Type: Gauge.

### kubevirt_vmi_migrations_in_running_phase
Number of current running migrations. Type: Gauge.

### kubevirt_vmi_migrations_in_scheduling_phase
Number of current scheduling migrations. Type: Gauge.

### kubevirt_vmi_network_receive_bytes_total
Total network traffic received in bytes. Type: Counter.

### kubevirt_vmi_network_receive_errors_total
Total network received error packets. Type: Counter.

### kubevirt_vmi_network_receive_packets_dropped_total
The total number of rx packets dropped on vNIC interfaces. Type: Counter.

### kubevirt_vmi_network_receive_packets_total
Total network traffic received packets. Type: Counter.

### kubevirt_vmi_network_traffic_bytes_total
Deprecated. Type: Counter.

### kubevirt_vmi_network_transmit_bytes_total
Total network traffic transmitted in bytes. Type: Counter.

### kubevirt_vmi_network_transmit_errors_total
Total network transmitted error packets. Type: Counter.

### kubevirt_vmi_network_transmit_packets_dropped_total
The total number of tx packets dropped on vNIC interfaces. Type: Counter.

### kubevirt_vmi_network_transmit_packets_total
Total network traffic transmitted packets. Type: Counter.

### kubevirt_vmi_node_cpu_affinity
Number of VMI CPU affinities to node physical cores. Type: Gauge.

### kubevirt_vmi_non_evictable
Indication for a VirtualMachine that its eviction strategy is set to Live Migration but is not migratable. Type: Gauge.

### kubevirt_vmi_number_of_outdated
Indication for the total number of VirtualMachineInstance workloads that are not running within the most up-to-date version of the virt-launcher environment. Type: Gauge.

### kubevirt_vmi_phase_count
Sum of VMIs per phase and node. `phase` can be one of the following: [`Pending`, `Scheduling`, `Scheduled`, `Running`, `Succeeded`, `Failed`, `Unknown`]. Type: Gauge.

### kubevirt_vmi_phase_transition_time_from_creation_seconds
Histogram of VM phase transitions duration from creation time in seconds. Type: Histogram.

### kubevirt_vmi_phase_transition_time_from_deletion_seconds
Histogram of VM phase transitions duration from deletion time in seconds. Type: Histogram.

### kubevirt_vmi_phase_transition_time_seconds
Histogram of VM phase transitions duration between different phases in seconds. Type: Histogram.

### kubevirt_vmi_storage_flush_requests_total
Total storage flush requests. Type: Counter.

### kubevirt_vmi_storage_flush_times_seconds_total
Total time spent on cache flushing. Type: Counter.

### kubevirt_vmi_storage_iops_read_total
Total number of I/O read operations. Type: Counter.

### kubevirt_vmi_storage_iops_write_total
Total number of I/O write operations. Type: Counter.

### kubevirt_vmi_storage_read_times_seconds_total
Total time spent on read operations. Type: Counter.

### kubevirt_vmi_storage_read_traffic_bytes_total
Total number of bytes read from storage. Type: Counter.

### kubevirt_vmi_storage_write_times_seconds_total
Total time spent on write operations. Type: Counter.

### kubevirt_vmi_storage_write_traffic_bytes_total
Total number of written bytes. Type: Counter.

### kubevirt_vmi_vcpu_delay_seconds_total
Amount of time spent by each vcpu waiting in the queue instead of running. Type: Counter.

### kubevirt_vmi_vcpu_seconds_total
Total amount of time spent in each state by each vcpu (cpu_time excluding hypervisor time). Where `id` is the vcpu identifier and `state` can be one of the following: [`OFFLINE`, `RUNNING`, `BLOCKED`]. Type: Counter.

### kubevirt_vmi_vcpu_wait_seconds_total
Amount of time spent by each vcpu while waiting on I/O. Type: Counter.

### kubevirt_vmsnapshot_disks_restored_from_source
Returns the total number of virtual machine disks restored from the source virtual machine. Type: Gauge.

### kubevirt_vmsnapshot_disks_restored_from_source_bytes
Returns the amount of space in bytes restored from the source virtual machine. Type: Gauge.

### kubevirt_vmsnapshot_persistentvolumeclaim_labels
Returns the labels of the persistent volume claims that are used for restoring virtual machines. Type: Gauge.

<div id='containerized-data-importer'></div>

## [containerized-data-importer](https://github.com/kubevirt/containerized-data-importer/tree/main)

### kubevirt_cdi_cr_ready
CDI install ready. Type: Gauge.

### kubevirt_cdi_dataimportcron_outdated
DataImportCron has an outdated import. Type: Gauge.

### kubevirt_cdi_datavolume_pending
Number of DataVolumes pending for default storage class to be configured. Type: Gauge.

### kubevirt_cdi_storageprofile_info
`StorageProfiles` info labels: `storageclass`, `provisioner`, `complete` indicates if all storage profiles recommended PVC settings are complete, `default` indicates if it's the Kubernetes default storage class, `virtdefault` indicates if it's the default virtualization storage class, `rwx` indicates if the storage class supports `ReadWriteMany`, `smartclone` indicates if it supports snapshot or CSI based clone. Type: Gauge.

### kubevirt_cdi_clone_pods_high_restart
The number of CDI clone pods with high restart count. Type: Gauge.

### kubevirt_cdi_import_pods_high_restart
The number of CDI import pods with high restart count. Type: Gauge.

### kubevirt_cdi_operator_up
CDI operator status. Type: Gauge.

### kubevirt_cdi_upload_pods_high_restart
The number of CDI upload server pods with high restart count. Type: Gauge.

<div id='cluster-network-addons-operator'></div>

## [cluster-network-addons-operator](https://github.com/kubevirt/cluster-network-addons-operator/tree/main)

### kubevirt_cnao_cr_kubemacpool_aggregated
Total count of KubeMacPool manager pods deployed by CNAO CR. Type: Gauge.

### kubevirt_cnao_cr_kubemacpool_deployed
KubeMacpool is deployed by CNAO CR. Type: Gauge.

### kubevirt_cnao_cr_ready
CNAO CR Ready. Type: Gauge.

### kubevirt_cnao_kubemacpool_duplicate_macs
Total count of duplicate KubeMacPool MAC addresses. Type: Gauge.

### kubevirt_cnao_kubemacpool_manager_up
Total count of running KubeMacPool manager pods. Type: Gauge.

### kubevirt_cnao_operator_up
Total count of running CNAO operators. Type: Gauge.

<div id='ssp-operator'></div>

## [ssp-operator](https://github.com/kubevirt/ssp-operator/tree/main)

### kubevirt_ssp_common_templates_restored_increase
The increase in the number of common templates restored by the operator back to their original state, over the last hour. Type: Gauge.
### kubevirt_ssp_common_templates_restored_total
The total number of common templates restored by the operator back to their original state. Type: Counter.
### kubevirt_ssp_operator_reconcile_succeeded
Set to 1 if the reconcile process of all operands completes with no errors, and to 0 otherwise. Type: Gauge.
### kubevirt_ssp_operator_reconcile_succeeded_aggregated
The total number of ssp-operator pods reconciling with no errors. Type: Gauge.
### kubevirt_ssp_operator_up
The total number of running ssp-operator pods. Type: Gauge.
### kubevirt_ssp_template_validator_rejected_increase
The increase in the number of rejected template validators, over the last hour. Type: Gauge.
### kubevirt_ssp_template_validator_rejected_total
The total number of rejected template validators. Type: Counter.
### kubevirt_ssp_template_validator_up
The total number of running virt-template-validator pods. Type: Gauge.
### kubevirt_ssp_vm_rbd_block_volume_without_rxbounce
VM with RBD mounted Block volume (without rxbounce option set). Type: Gauge.
<div id='hostpath-provisioner-operator'></div>

## [hostpath-provisioner-operator](https://github.com/kubevirt/hostpath-provisioner-operator/tree/main)

### kubevirt_hpp_cr_ready
HPP CR Ready. Type: Gauge.

### kubevirt_hpp_operator_up
The number of running hostpath-provisioner-operator pods. Type: Gauge.

<div id='hostpath-provisioner'></div>

## [hostpath-provisioner](https://github.com/kubevirt/hostpath-provisioner/tree/main)

### kubevirt_hpp_pool_path_shared_with_os
HPP pool path sharing a filesystem with OS, fix to prevent HPP PVs from causing disk pressure and affecting node operation. Type: Gauge.

<div id='hyperconverged-cluster-operator'></div>

## [hyperconverged-cluster-operator](https://github.com/kubevirt/hyperconverged-cluster-operator/tree/main)

### kubevirt_hco_hyperconverged_cr_exists
Indicates whether the HyperConverged custom resource exists (1) or not (0). Type: Gauge.

### kubevirt_hco_out_of_band_modifications_total
Count of out-of-band modifications overwritten by HCO. Type: Counter.

### kubevirt_hco_single_stack_ipv6
Indicates whether the underlying cluster is single stack IPv6 (1) or not (0). Type: Gauge.

### kubevirt_hco_system_health_status
Indicates whether the system health status is healthy (0), warning (1), or error (2), by aggregating the conditions of HCO and its secondary resources. Type: Gauge.

### kubevirt_hco_unsafe_modifications
Count of unsafe modifications in the HyperConverged annotations. Type: Gauge.

### kubevirt_hyperconverged_operator_health_status
Indicates whether HCO and its secondary resources health status is healthy (0), warning (1) or critical (2), based both on the firing alerts that impact the operator health, and on kubevirt_hco_system_health_status metric. Type: Gauge.

