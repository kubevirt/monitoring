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

### kubevirt_virt_controller_leading
Indication for an operating virt-controller.

### kubevirt_virt_controller_ready
Indication for a virt-controller that is ready to take the lead.

### kubevirt_vmi_cpu_affinity
The vcpu affinity details.

### kubevirt_vmi_memory_actual_balloon_bytes
Current balloon bytes.

### kubevirt_vmi_memory_available_bytes
Amount of `usable` memory as seen by the domain.

### kubevirt_vmi_memory_domain_total_bytes
The amount of memory in bytes allocated to the domain. The `memory` value in domain xml file.

### kubevirt_vmi_memory_pgmajfault
The number of page faults when disk IO was required.

### kubevirt_vmi_memory_pgminfault
The number of other page faults, when disk IO was not required.

### kubevirt_vmi_memory_resident_bytes
Resident set size of the process running the domain.

### kubevirt_vmi_memory_swap_in_traffic_bytes_total
Swap in memory traffic in bytes.

### kubevirt_vmi_memory_swap_out_traffic_bytes_total
Swap out memory traffic in bytes.

### kubevirt_vmi_memory_unused_bytes
Amount of `unused` memory as seen by the domain.

### kubevirt_vmi_memory_usable_bytes
The amount of memory which can be reclaimed by balloon without causing host swapping in bytes.

### kubevirt_vmi_memory_used_bytes
Amount of `used` memory as seen by the domain.

### kubevirt_vmi_network_receive_bytes_total
Network traffic receive in bytes.

### kubevirt_vmi_network_receive_errors_total
Network receive error packets.

### kubevirt_vmi_network_receive_packets_dropped_total
The number of rx packets dropped on vNIC interfaces.

### kubevirt_vmi_network_receive_packets_total
Network traffic receive packets.

### kubevirt_vmi_network_traffic_bytes_total
Deprecated.

### kubevirt_vmi_network_transmit_bytes_total
Network traffic transmit in bytes.

### kubevirt_vmi_network_transmit_errors_total
Network transmit error packets.

### kubevirt_vmi_network_transmit_packets_dropped_total
The number of tx packets dropped on vNIC interfaces.

### kubevirt_vmi_network_transmit_packets_total
Network traffic transmit packets.

### kubevirt_vmi_non_evictable
Indication for a VirtualMachine that its eviction strategy is set to Live Migration but is not migratable.

### kubevirt_vmi_outdated_count
Indication for the number of VirtualMachineInstance workloads that are not running within the most up-to-date version of the virt-launcher environment.

### kubevirt_vmi_phase_count
Sum of VMIs per phase and node.

`phase` can be one of the following: [`Pending`, `Scheduling`, `Scheduled`, `Running`, `Succeeded`, `Failed`, `Unknown`]

### kubevirt_vmi_storage_flush_requests_total
Storage flush requests.

### kubevirt_vmi_storage_flush_times_ms_total
Total time (ms) spent on cache flushing.

### kubevirt_vmi_storage_iops_read_total
I/O read operations.

### kubevirt_vmi_storage_iops_write_total
I/O write operations.

### kubevirt_vmi_storage_read_times_ms_total
Storage read operation time.

### kubevirt_vmi_storage_read_traffic_bytes_total
Storage read traffic in bytes.

### kubevirt_vmi_storage_write_times_ms_total
Storage write operation time.

### kubevirt_vmi_storage_write_traffic_bytes_total
Storage write traffic in bytes.

### kubevirt_vmi_vcpu_seconds
Amount of time spent in each state by each vcpu. Where `id` is the vcpu identifier and `state` can be one of the following: [`OFFLINE`, `RUNNING`, `BLOCKED`].

### kubevirt_vmi_vcpu_wait_seconds
Amount of time spent by each vcpu while waiting on I/O.

<div id='containerized-data-importer'></div>

## [containerized-data-importer](https://github.com/kubevirt/containerized-data-importer/tree/main)

### kubevirt_cdi_clone_dv_unusual_restartcount_total
Total restart count in CDI Data Volume cloner pod
### kubevirt_cdi_dataimportcron_outdated_total
Total count of outdated DataImportCron imports
### kubevirt_cdi_import_dv_unusual_restartcount_total
Total restart count in CDI Data Volume importer pod
### kubevirt_cdi_operator_up_total
CDI operator status
### kubevirt_cdi_upload_dv_unusual_restartcount_total
Total restart count in CDI Data Volume upload server pod
<div id='cluster-network-addons-operator'></div>

## [cluster-network-addons-operator](https://github.com/kubevirt/cluster-network-addons-operator/tree/main)

### kubevirt_cnao_cr_kubemacpool_deployed
KubeMacpool is deployed by CNAO CR
### kubevirt_cnao_cr_ready
CNAO CR Ready
<div id='ssp-operator'></div>

## [ssp-operator](https://github.com/kubevirt/ssp-operator/tree/master)

### kubevirt_ssp_num_of_operator_reconciling_properly
The total number of ssp-operator pods reconciling with no errors
### kubevirt_ssp_operator_up_total
The total number of running ssp-operator pods
### kubevirt_ssp_rejected_vms_total
The total number of vms rejected by virt-template-validator
### kubevirt_ssp_template_validator_up_total
The total number of running virt-template-validator pods
### kubevirt_ssp_total_restored_common_templates
The total number of common templates restored by the operator back to their original state
<div id='hostpath-provisioner-operator'></div>

## [hostpath-provisioner-operator](https://github.com/kubevirt/hostpath-provisioner-operator/tree/main)

### kubevirt_hpp_operator_up_total
The total number of running hostpath-provisioner-operator pods
<div id='hyperconverged-cluster-operator'></div>

## [hyperconverged-cluster-operator](https://github.com/kubevirt/hyperconverged-cluster-operator/tree/main)

### kubevirt_hco_out_of_band_modifications_count
Count of out-of-band modifications overwritten by HCO
### kubevirt_hco_unsafe_modification_count
Count of unsafe modifications in the HyperConverged annotations
