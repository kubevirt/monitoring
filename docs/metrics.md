# KubeVirt components metrics

This document aims to help users that are not familiar with metrics exposed by all the KubeVirt components.
All metrics documented here are auto-generated in each component repository and gathered here.
They reflect and describe exactly what is being exposed.

## Operator Repositories

| Operator Name |
|---------------|
| [kubevirt](https://github.com/kubevirt/kubevirt/tree/main) |
| [cluster-network-addons-operator](https://github.com/kubevirt/cluster-network-addons-operator/tree/main) |
| [containerized-data-importer](https://github.com/kubevirt/containerized-data-importer/tree/main) |
| [hostpath-provisioner](https://github.com/kubevirt/hostpath-provisioner/tree/main) |
| [hostpath-provisioner-operator](https://github.com/kubevirt/hostpath-provisioner-operator/tree/main) |
| [hyperconverged-cluster-operator](https://github.com/kubevirt/hyperconverged-cluster-operator/tree/main) |
| [ssp-operator](https://github.com/kubevirt/ssp-operator/tree/main) |


## Metrics

The following table contains all metrics from operators listed above. Each row represents a single metric with its operator, name, type (Gauge, Counter, Histogram, etc.), and description.

| Operator Name | Metric Name | Type | Description |
|----------|------|------|-------------|
| kubevirt | `cluster:kubevirt_virt_controller_pods_running:count` | Gauge | The number of virt-controller pods that are running |
| kubevirt | `kubevirt_allocatable_nodes` | Gauge | The number of allocatable nodes in the cluster |
| kubevirt | `kubevirt_api_request_deprecated_total` | Counter | The total number of requests to deprecated KubeVirt APIs |
| kubevirt | `kubevirt_configuration_emulation_enabled` | Gauge | Indicates whether the Software Emulation is enabled in the configuration |
| kubevirt | `kubevirt_console_active_connections` | Gauge | Amount of active Console connections, broken down by namespace and vmi name |
| kubevirt | `kubevirt_info` | Gauge | Version information |
| kubevirt | `kubevirt_memory_delta_from_requested_bytes` | Gauge | The delta between the pod with highest memory working set or rss and its requested memory for each container, virt-controller, virt-handler, virt-api, virt-operator and compute(virt-launcher) |
| kubevirt | `kubevirt_node_deprecated_machine_types` | Gauge | List of deprecated machine types based on the capabilities of individual nodes, as detected by virt-handler |
| kubevirt | `kubevirt_nodes_with_kvm` | Gauge | The number of nodes in the cluster that have the devices.kubevirt.io/kvm resource available |
| kubevirt | `kubevirt_number_of_vms` | Gauge | The number of VMs in the cluster by namespace |
| kubevirt | `kubevirt_portforward_active_tunnels` | Gauge | Amount of active portforward tunnels, broken down by namespace and vmi name |
| kubevirt | `kubevirt_rest_client_rate_limiter_duration_seconds` | Histogram | Client side rate limiter latency in seconds. Broken down by verb and URL |
| kubevirt | `kubevirt_rest_client_request_latency_seconds` | Histogram | Request latency in seconds. Broken down by verb and URL |
| kubevirt | `kubevirt_rest_client_requests_total` | Counter | Number of HTTP requests, partitioned by status code, method, and host |
| kubevirt | `kubevirt_usbredir_active_connections` | Gauge | Amount of active USB redirection connections, broken down by namespace and vmi name |
| kubevirt | `kubevirt_virt_api_up` | Gauge | The number of virt-api pods that are up |
| kubevirt | `kubevirt_virt_controller_leading_status` | Gauge | Indication for an operating virt-controller |
| kubevirt | `kubevirt_virt_controller_ready` | Gauge | The number of virt-controller pods that are ready |
| kubevirt | `kubevirt_virt_controller_ready_status` | Gauge | Indication for a virt-controller that is ready to take the lead |
| kubevirt | `kubevirt_virt_controller_up` | Gauge | The number of virt-controller pods that are up |
| kubevirt | `kubevirt_virt_handler_up` | Gauge | The number of virt-handler pods that are up |
| kubevirt | `kubevirt_virt_operator_leading` | Gauge | The number of virt-operator pods that are leading |
| kubevirt | `kubevirt_virt_operator_leading_status` | Gauge | Indication for an operating virt-operator |
| kubevirt | `kubevirt_virt_operator_ready` | Gauge | The number of virt-operator pods that are ready |
| kubevirt | `kubevirt_virt_operator_ready_status` | Gauge | Indication for a virt-operator that is ready to take the lead |
| kubevirt | `kubevirt_virt_operator_up` | Gauge | The number of virt-operator pods that are up |
| kubevirt | `kubevirt_vm_container_memory_request_margin_based_on_rss_bytes` | Gauge | Difference between requested memory and rss for VM containers (request margin). Can be negative when usage exceeds request |
| kubevirt | `kubevirt_vm_container_memory_request_margin_based_on_working_set_bytes` | Gauge | Difference between requested memory and working set for VM containers (request margin). Can be negative when usage exceeds request |
| kubevirt | `kubevirt_vm_create_date_timestamp_seconds` | Gauge | Virtual Machine creation timestamp |
| kubevirt | `kubevirt_vm_created_by_pod_total` | Counter | The total number of VMs created by namespace and virt-api pod, since install |
| kubevirt | `kubevirt_vm_created_total` | Counter | The total number of VMs created by namespace, since install |
| kubevirt | `kubevirt_vm_disk_allocated_size_bytes` | Gauge | Allocated disk size of a Virtual Machine in bytes, based on its PersistentVolumeClaim. Includes persistentvolumeclaim (PVC name), volume_mode (disk presentation mode: Filesystem or Block), and device (disk name) |
| kubevirt | `kubevirt_vm_error_status_last_transition_timestamp_seconds` | Counter | Virtual Machine last transition timestamp to error status |
| kubevirt | `kubevirt_vm_info` | Gauge | Information about Virtual Machines |
| kubevirt | `kubevirt_vm_labels` | Gauge | The metric exposes the VM labels as Prometheus labels. Configure allowed and ignored labels via the 'kubevirt-vm-labels-config' ConfigMap |
| kubevirt | `kubevirt_vm_migrating_status_last_transition_timestamp_seconds` | Counter | Virtual Machine last transition timestamp to migrating status |
| kubevirt | `kubevirt_vm_non_running_status_last_transition_timestamp_seconds` | Counter | Virtual Machine last transition timestamp to paused/stopped status |
| kubevirt | `kubevirt_vm_resource_limits` | Gauge | Resources limits by Virtual Machine. Reports memory and CPU limits |
| kubevirt | `kubevirt_vm_resource_requests` | Gauge | Resources requested by Virtual Machine. Reports memory and CPU requests |
| kubevirt | `kubevirt_vm_running_status_last_transition_timestamp_seconds` | Counter | Virtual Machine last transition timestamp to running status |
| kubevirt | `kubevirt_vm_starting_status_last_transition_timestamp_seconds` | Counter | Virtual Machine last transition timestamp to starting status |
| kubevirt | `kubevirt_vm_vnic_info` | Gauge | Details of Virtual Machine (VM) vNIC interfaces, such as vNIC name, binding type, network name, and binding name for each vNIC defined in the VM's configuration |
| kubevirt | `kubevirt_vmi_cpu_system_usage_seconds_total` | Counter | Total CPU time spent in system mode |
| kubevirt | `kubevirt_vmi_cpu_usage_seconds_total` | Counter | Total CPU time spent in all modes (sum of both vcpu and hypervisor usage) |
| kubevirt | `kubevirt_vmi_cpu_user_usage_seconds_total` | Counter | Total CPU time spent in user mode |
| kubevirt | `kubevirt_vmi_dirty_rate_bytes_per_second` | Gauge | Guest dirty-rate in bytes per second |
| kubevirt | `kubevirt_vmi_filesystem_capacity_bytes` | Gauge | Total VM filesystem capacity in bytes |
| kubevirt | `kubevirt_vmi_filesystem_used_bytes` | Gauge | Used VM filesystem capacity in bytes |
| kubevirt | `kubevirt_vmi_guest_load_15m` | Gauge | Guest system load average over 15 minutes as reported by the guest agent. Load is defined as the number of processes in the runqueue or waiting for disk I/O. Requires qemu-guest-agent version 10.0.0 or above |
| kubevirt | `kubevirt_vmi_guest_load_1m` | Gauge | Guest system load average over 1 minute as reported by the guest agent. Load is defined as the number of processes in the runqueue or waiting for disk I/O. Requires qemu-guest-agent version 10.0.0 or above |
| kubevirt | `kubevirt_vmi_guest_load_5m` | Gauge | Guest system load average over 5 minutes as reported by the guest agent. Load is defined as the number of processes in the runqueue or waiting for disk I/O. Requires qemu-guest-agent version 10.0.0 or above |
| kubevirt | `kubevirt_vmi_guest_vcpu_queue` | Gauge | Guest queue length |
| kubevirt | `kubevirt_vmi_info` | Gauge | Information about VirtualMachineInstances |
| kubevirt | `kubevirt_vmi_last_api_connection_timestamp_seconds` | Gauge | Virtual Machine Instance last API connection timestamp. Including VNC, console, portforward, SSH and usbredir connections |
| kubevirt | `kubevirt_vmi_launcher_memory_overhead_bytes` | Gauge | Estimation of the memory amount required for virt-launcher's infrastructure components (e.g. libvirt, QEMU) |
| kubevirt | `kubevirt_vmi_memory_actual_balloon_bytes` | Gauge | Current balloon size in bytes |
| kubevirt | `kubevirt_vmi_memory_available_bytes` | Gauge | Amount of usable memory as seen by the domain. This value may not be accurate if a balloon driver is in use or if the guest OS does not initialize all assigned pages |
| kubevirt | `kubevirt_vmi_memory_cached_bytes` | Gauge | The amount of memory that is being used to cache I/O and is available to be reclaimed, corresponds to the sum of `Buffers` + `Cached` + `SwapCached` in `/proc/meminfo` |
| kubevirt | `kubevirt_vmi_memory_domain_bytes` | Gauge | The amount of memory in bytes allocated to the domain. The `memory` value in domain xml file |
| kubevirt | `kubevirt_vmi_memory_pgmajfault_total` | Counter | The number of page faults when disk IO was required. Page faults occur when a process makes a valid access to virtual memory that is not available. When servicing the page fault, if disk IO is required, it is considered as major fault |
| kubevirt | `kubevirt_vmi_memory_pgminfault_total` | Counter | The number of other page faults, when disk IO was not required. Page faults occur when a process makes a valid access to virtual memory that is not available. When servicing the page fault, if disk IO is NOT required, it is considered as minor fault |
| kubevirt | `kubevirt_vmi_memory_resident_bytes` | Gauge | Resident set size of the process running the domain |
| kubevirt | `kubevirt_vmi_memory_swap_in_traffic_bytes` | Gauge | The total amount of data read from swap space of the guest in bytes |
| kubevirt | `kubevirt_vmi_memory_swap_out_traffic_bytes` | Gauge | The total amount of memory written out to swap space of the guest in bytes |
| kubevirt | `kubevirt_vmi_memory_unused_bytes` | Gauge | The amount of memory left completely unused by the system. Memory that is available but used for reclaimable caches should NOT be reported as free |
| kubevirt | `kubevirt_vmi_memory_usable_bytes` | Gauge | The amount of memory which can be reclaimed by balloon without pushing the guest system to swap, corresponds to 'Available' in /proc/meminfo |
| kubevirt | `kubevirt_vmi_memory_used_bytes` | Gauge | Amount of `used` memory as seen by the domain |
| kubevirt | `kubevirt_vmi_migration_data_processed_bytes` | Gauge | The total Guest OS data processed and migrated to the new VM |
| kubevirt | `kubevirt_vmi_migration_data_remaining_bytes` | Gauge | The remaining guest OS data to be migrated to the new VM |
| kubevirt | `kubevirt_vmi_migration_data_total_bytes` | Counter | The total Guest OS data to be migrated to the new VM |
| kubevirt | `kubevirt_vmi_migration_dirty_memory_rate_bytes` | Gauge | The rate of memory being dirty in the Guest OS |
| kubevirt | `kubevirt_vmi_migration_end_time_seconds` | Gauge | The time at which the migration ended |
| kubevirt | `kubevirt_vmi_migration_failed` | Gauge | Indicates if the VMI migration failed |
| kubevirt | `kubevirt_vmi_migration_memory_transfer_rate_bytes` | Gauge | The rate at which the memory is being transferred |
| kubevirt | `kubevirt_vmi_migration_phase_transition_time_from_creation_seconds` | Histogram | Histogram of VM migration phase transitions duration from creation time in seconds |
| kubevirt | `kubevirt_vmi_migration_start_time_seconds` | Gauge | The time at which the migration started |
| kubevirt | `kubevirt_vmi_migration_succeeded` | Gauge | Indicates if the VMI migration succeeded |
| kubevirt | `kubevirt_vmi_migrations_in_pending_phase` | Gauge | Number of current pending migrations |
| kubevirt | `kubevirt_vmi_migrations_in_running_phase` | Gauge | Number of current running migrations |
| kubevirt | `kubevirt_vmi_migrations_in_scheduling_phase` | Gauge | Number of current scheduling migrations |
| kubevirt | `kubevirt_vmi_migrations_in_unset_phase` | Gauge | Number of current unset migrations. These are pending items the virt-controller hasn’t processed yet from the queue |
| kubevirt | `kubevirt_vmi_network_receive_bytes_total` | Counter | Total network traffic received in bytes |
| kubevirt | `kubevirt_vmi_network_receive_errors_total` | Counter | Total network received error packets |
| kubevirt | `kubevirt_vmi_network_receive_packets_dropped_total` | Counter | The total number of rx packets dropped on vNIC interfaces |
| kubevirt | `kubevirt_vmi_network_receive_packets_total` | Counter | Total network traffic received packets |
| kubevirt | `kubevirt_vmi_network_traffic_bytes_total` | Counter | [Deprecated] Total number of bytes sent and received |
| kubevirt | `kubevirt_vmi_network_transmit_bytes_total` | Counter | Total network traffic transmitted in bytes |
| kubevirt | `kubevirt_vmi_network_transmit_errors_total` | Counter | Total network transmitted error packets |
| kubevirt | `kubevirt_vmi_network_transmit_packets_dropped_total` | Counter | The total number of tx packets dropped on vNIC interfaces |
| kubevirt | `kubevirt_vmi_network_transmit_packets_total` | Counter | Total network traffic transmitted packets |
| kubevirt | `kubevirt_vmi_node_cpu_affinity` | Gauge | Number of VMI CPU affinities to node physical cores |
| kubevirt | `kubevirt_vmi_non_evictable` | Gauge | Indication for a VirtualMachine that its eviction strategy is set to Live Migration but is not migratable |
| kubevirt | `kubevirt_vmi_number_of_outdated` | Gauge | Indication for the total number of VirtualMachineInstance workloads that are not running within the most up-to-date version of the virt-launcher environment |
| kubevirt | `kubevirt_vmi_phase_count` | Gauge | Sum of VMIs per phase and node. `phase` can be one of the following: [`Pending`, `Scheduling`, `Scheduled`, `Running`, `Succeeded`, `Failed`, `Unknown`] |
| kubevirt | `kubevirt_vmi_phase_transition_time_from_creation_seconds` | Histogram | Histogram of VM phase transitions duration from creation time in seconds |
| kubevirt | `kubevirt_vmi_phase_transition_time_from_deletion_seconds` | Histogram | Histogram of VM phase transitions duration from deletion time in seconds |
| kubevirt | `kubevirt_vmi_phase_transition_time_seconds` | Histogram | Histogram of VM phase transitions duration between different phases in seconds |
| kubevirt | `kubevirt_vmi_status_addresses` | Gauge | The addresses of a VirtualMachineInstance. This metric provides the address of an available network interface associated with the VMI in the 'address' label, and about the type of address, such as internal IP, in the 'type' label |
| kubevirt | `kubevirt_vmi_storage_flush_requests_total` | Counter | Total storage flush requests |
| kubevirt | `kubevirt_vmi_storage_flush_times_seconds_total` | Counter | Total time spent on cache flushing |
| kubevirt | `kubevirt_vmi_storage_iops_read_total` | Counter | Total number of I/O read operations |
| kubevirt | `kubevirt_vmi_storage_iops_write_total` | Counter | Total number of I/O write operations |
| kubevirt | `kubevirt_vmi_storage_read_times_seconds_total` | Counter | Total time spent on read operations |
| kubevirt | `kubevirt_vmi_storage_read_traffic_bytes_total` | Counter | Total number of bytes read from storage |
| kubevirt | `kubevirt_vmi_storage_write_times_seconds_total` | Counter | Total time spent on write operations |
| kubevirt | `kubevirt_vmi_storage_write_traffic_bytes_total` | Counter | Total number of written bytes |
| kubevirt | `kubevirt_vmi_vcpu_delay_seconds_total` | Counter | Amount of time spent by each vcpu waiting in the queue instead of running |
| kubevirt | `kubevirt_vmi_vcpu_seconds_total` | Counter | Total amount of time spent in each state by each vcpu (cpu_time excluding hypervisor time). Where `id` is the vcpu identifier and `state` can be one of the following: [`OFFLINE`, `RUNNING`, `BLOCKED`] |
| kubevirt | `kubevirt_vmi_vcpu_wait_seconds_total` | Counter | Amount of time spent by each vcpu while waiting on I/O |
| kubevirt | `kubevirt_vmi_vnic_info` | Gauge | Details of VirtualMachineInstance (VMI) vNIC interfaces, such as vNIC name, binding type, network name, and binding name for each vNIC of a running instance |
| kubevirt | `kubevirt_vmsnapshot_disks_restored_from_source` | Gauge | Returns the total number of virtual machine disks restored from the source virtual machine |
| kubevirt | `kubevirt_vmsnapshot_disks_restored_from_source_bytes` | Gauge | Returns the amount of space in bytes restored from the source virtual machine |
| kubevirt | `kubevirt_vmsnapshot_persistentvolumeclaim_labels` | Gauge | Returns the labels of the persistent volume claims that are used for restoring virtual machines |
| kubevirt | `kubevirt_vmsnapshot_succeeded_timestamp_seconds` | Gauge | Returns the timestamp of successful virtual machine snapshot |
| kubevirt | `kubevirt_vnc_active_connections` | Gauge | Amount of active VNC connections, broken down by namespace and vmi name |
| kubevirt | `kubevirt_workqueue_adds_total` | Counter | Total number of adds handled by workqueue |
| kubevirt | `kubevirt_workqueue_depth` | Gauge | Current depth of workqueue |
| kubevirt | `kubevirt_workqueue_longest_running_processor_seconds` | Gauge | How many seconds has the longest running processor for workqueue been running |
| kubevirt | `kubevirt_workqueue_queue_duration_seconds` | Histogram | How long an item stays in workqueue before being requested |
| kubevirt | `kubevirt_workqueue_retries_total` | Counter | Total number of retries handled by workqueue |
| kubevirt | `kubevirt_workqueue_unfinished_work_seconds` | Gauge | How many seconds of work has done that is in progress and hasn't been observed by work_duration. Large values indicate stuck threads. One can deduce the number of stuck threads by observing the rate at which this increases |
| kubevirt | `kubevirt_workqueue_work_duration_seconds` | Histogram | How long in seconds processing an item from workqueue takes |
| kubevirt | `vmi:kubevirt_vmi_vcpu:count` | Gauge | The number of the VMI vCPUs |
| cluster-network-addons-operator | `kubevirt_cnao_cr_kubemacpool_aggregated` | Gauge | Total count of KubeMacPool manager pods deployed by CNAO CR |
| cluster-network-addons-operator | `kubevirt_cnao_cr_kubemacpool_deployed` | Gauge | KubeMacpool is deployed by CNAO CR |
| cluster-network-addons-operator | `kubevirt_cnao_cr_ready` | Gauge | CNAO CR Ready |
| cluster-network-addons-operator | `kubevirt_cnao_kubemacpool_duplicate_macs` | Gauge | Total count of duplicate KubeMacPool MAC addresses |
| cluster-network-addons-operator | `kubevirt_cnao_kubemacpool_manager_up` | Gauge | Total count of running KubeMacPool manager pods |
| cluster-network-addons-operator | `kubevirt_cnao_operator_up` | Gauge | Total count of running CNAO operators |
| containerized-data-importer | `kubevirt_cdi_clone_pods_high_restart` | Gauge | The number of CDI clone pods with high restart count |
| containerized-data-importer | `kubevirt_cdi_clone_progress_total` | Counter | The clone progress in percentage |
| containerized-data-importer | `kubevirt_cdi_cr_ready` | Gauge | CDI install ready |
| containerized-data-importer | `kubevirt_cdi_dataimportcron_outdated` | Gauge | DataImportCron has an outdated import |
| containerized-data-importer | `kubevirt_cdi_datavolume_pending` | Gauge | Number of DataVolumes pending for default storage class to be configured |
| containerized-data-importer | `kubevirt_cdi_import_pods_high_restart` | Gauge | The number of CDI import pods with high restart count |
| containerized-data-importer | `kubevirt_cdi_import_progress_total` | Counter | The import progress in percentage |
| containerized-data-importer | `kubevirt_cdi_openstack_populator_progress_total` | Counter | Progress of volume population |
| containerized-data-importer | `kubevirt_cdi_operator_up` | Gauge | CDI operator status |
| containerized-data-importer | `kubevirt_cdi_ovirt_progress_total` | Counter | Progress of volume population |
| containerized-data-importer | `kubevirt_cdi_storageprofile_info` | Gauge | `StorageProfiles` info labels: `storageclass`, `provisioner`, `complete` indicates if all storage profiles recommended PVC settings are complete, `default` indicates if it's the Kubernetes default storage class, `virtdefault` indicates if it's the default virtualization storage class, `rwx` indicates if the storage class supports `ReadWriteMany`, `smartclone` indicates if it supports snapshot or CSI based clone, `degraded` indicates it is not optimal for virtualization |
| containerized-data-importer | `kubevirt_cdi_upload_pods_high_restart` | Gauge | The number of CDI upload server pods with high restart count |
| hostpath-provisioner | `kubevirt_hpp_pool_path_shared_with_os` | Gauge | HPP pool path sharing a filesystem with OS, fix to prevent HPP PVs from causing disk pressure and affecting node operation |
| hostpath-provisioner-operator | `kubevirt_hpp_cr_ready` | Gauge | HPP CR Ready |
| hostpath-provisioner-operator | `kubevirt_hpp_operator_up` | Gauge | The number of running hostpath-provisioner-operator pods |
| hyperconverged-cluster-operator | `cluster:vmi_request_cpu_cores:sum` | Gauge | Sum of CPU core requests for all running virt-launcher VMIs across the entire Kubevirt cluster |
| hyperconverged-cluster-operator | `cnv_abnormal` | Gauge | Monitors resources for potential problems |
| hyperconverged-cluster-operator | `kubevirt_hco_dataimportcrontemplate_with_architecture_annotation` | Gauge | Indicates whether the DataImportCronTemplate has the ssp.kubevirt.io/dict.architectures annotation (0) or not (1) |
| hyperconverged-cluster-operator | `kubevirt_hco_dataimportcrontemplate_with_supported_architectures` | Gauge | Indicates whether the DataImportCronTemplate has supported architectures (0) or not (1) |
| hyperconverged-cluster-operator | `kubevirt_hco_hyperconverged_cr_exists` | Gauge | Indicates whether the HyperConverged custom resource exists (1) or not (0) |
| hyperconverged-cluster-operator | `kubevirt_hco_memory_overcommit_percentage` | Gauge | Indicates the cluster-wide configured VM memory overcommit percentage |
| hyperconverged-cluster-operator | `kubevirt_hco_misconfigured_descheduler` | Gauge | Indicates whether the optional descheduler is not properly configured (1) to work with KubeVirt or not (0) |
| hyperconverged-cluster-operator | `kubevirt_hco_out_of_band_modifications_total` | Counter | Count of out-of-band modifications overwritten by HCO |
| hyperconverged-cluster-operator | `kubevirt_hco_single_stack_ipv6` | Gauge | Indicates whether the underlying cluster is single stack IPv6 (1) or not (0) |
| hyperconverged-cluster-operator | `kubevirt_hco_system_health_status` | Gauge | Indicates whether the system health status is healthy (0), warning (1), or error (2), by aggregating the conditions of HCO and its secondary resources |
| hyperconverged-cluster-operator | `kubevirt_hco_unsafe_modifications` | Gauge | Count of unsafe modifications in the HyperConverged annotations |
| hyperconverged-cluster-operator | `kubevirt_hyperconverged_operator_health_status` | Gauge | Indicates whether HCO and its secondary resources health status is healthy (0), warning (1) or critical (2), based both on the firing alerts that impact the operator health, and on kubevirt_hco_system_health_status metric |
| ssp-operator | `cnv:vmi_status_running:count` | Gauge | The total number of running VMIs, labeled with node, instance type, preference and guest OS information |
| ssp-operator | `kubevirt_ssp_common_templates_restored_increase` | Gauge | The increase in the number of common templates restored by the operator back to their original state, over the last hour |
| ssp-operator | `kubevirt_ssp_common_templates_restored_total` | Counter | The total number of common templates restored by the operator back to their original state |
| ssp-operator | `kubevirt_ssp_operator_reconcile_succeeded` | Gauge | Set to 1 if the reconcile process of all operands completes with no errors, and to 0 otherwise |
| ssp-operator | `kubevirt_ssp_operator_reconcile_succeeded_aggregated` | Gauge | The total number of ssp-operator pods reconciling with no errors |
| ssp-operator | `kubevirt_ssp_operator_up` | Gauge | The total number of running ssp-operator pods |
| ssp-operator | `kubevirt_ssp_template_validator_rejected_increase` | Gauge | The increase in the number of rejected template validators, over the last hour |
| ssp-operator | `kubevirt_ssp_template_validator_rejected_total` | Counter | The total number of rejected template validators |
| ssp-operator | `kubevirt_ssp_template_validator_up` | Gauge | The total number of running virt-template-validator pods |
| ssp-operator | `kubevirt_ssp_vm_rbd_block_volume_without_rxbounce` | Gauge | [ALPHA] VM with RBD mounted Block volume (without rxbounce option set) |


