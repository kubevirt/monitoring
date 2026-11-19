# OutdatedGuestDeviceDrivers

## Meaning

This alert fires when outdated VirtIO drivers are detected in a Windows virtual
machine (VM) for more than 1 hour.

The QEMU guest agent reports the name and version of the driver of each VirtIO
device in the guest. A driver is outdated when its version is not the latest
version provided by the cluster for that device and Windows version.

The alert fires once per VM and device.

## Impact

The VM continues to run, but its drivers might lack performance, stability and
security fixes that are included in the latest drivers provided by the cluster.

The cluster and other VMs are not affected.

## Diagnosis

Run the following PromQL queries, replacing the placeholders with the values of
the alert labels.

<!--USstart-->
**Note:** You can use the Prometheus metrics explorer to run the queries.
<!--USend-->

<!--DS: **Note:** You can run the queries in the web console by navigating to
**Observe** -> **Metrics**.-->

1. List all outdated drivers of the VM:

   ```promql
   vmi:kubevirt_vmi_guest_device_driver_timestamp_seconds:max{namespace="<namespace>", name="<name>"}
   ```

   Each result is one outdated driver. The `driver_name` and `driver_version`
   labels show the installed driver.

2. Obtain the latest driver versions provided by the cluster for a device and
   Windows version:

   ```promql
   kubevirt_guest_device_driver_latest_version_info{device_id="<device_id>", guest_os_version_id="<guest_os_version_id>"}
   ```

   The `driver_version` label of each result is a latest version.

If the alert fires for many VMs at the same time, the latest driver versions
provided by the cluster have most likely changed, for example after a cluster
upgrade. In this case, all affected VMs need a driver update.

## Mitigation

The VM owner must update all VirtIO drivers in the VM to the latest drivers
provided by the cluster. The cluster provides the drivers as a container disk.

<!--DS: Update the drivers by following
[Update VirtIO drivers and the guest agent on a Windows VM](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/virtualization/managing-vms#virt-updating-virtio-drivers-guest-agent-windows_virt-update-virtio-drivers).-->
<!--USstart-->
1. Set the environment variables from the alert labels:

   ```bash
   $ export NAMESPACE="<namespace label>"
   $ export VM_NAME="<name label>"
   ```

2. Obtain the VirtIO container disk image provided by the cluster:

   - On clusters with HCO, read the image from the `virtio-win` ConfigMap in
     the namespace of the `HyperConverged` resource:

     ```bash
     $ export HCO_NAMESPACE="$(kubectl get hyperconverged -A --no-headers | awk '{print $1}')"
     $ kubectl get configmap virtio-win -n $HCO_NAMESPACE \
       -o jsonpath='{.data.virtio-win-image}'
     ```

     If you cannot list `HyperConverged` resources, ask your cluster
     administrator for the namespace. The default namespace is
     `kubevirt-hyperconverged`.

   - On clusters without HCO, use the `quay.io/kubevirt/virtio-container-disk`
     image with the tag matching your KubeVirt version, for example `v1.9.1`.
     Do not use the `latest` tag, it is outdated.

3. Attach the container disk to the VM as a SATA CD drive by editing the VM:

   ```bash
   $ kubectl edit vm $VM_NAME -n $NAMESPACE
   ```

   Add the following disk and volume, replacing `<virtio_win_image>` with the
   image from the previous step:

   ```yaml
   spec:
     template:
       spec:
         domain:
           devices:
             disks:
               - name: virtiocontainerdisk
                 cdrom:
                   bus: sata
         volumes:
           - name: virtiocontainerdisk
             containerDisk:
               image: <virtio_win_image>
   ```

4. Restart the VM to apply the change:

   ```bash
   $ virtctl restart $VM_NAME -n $NAMESPACE
   ```

5. In the Windows guest, open **File Explorer**, navigate to the `virtio-win`
   CD drive and run the `virtio-win-gt-x64` installer. The default options
   update all VirtIO drivers and the QEMU guest agent.

6. Reboot the VM after the installation is complete.

For more information, see
[Windows virtio drivers](https://kubevirt.io/user-guide/user_workloads/windows_virtio_drivers/).
<!--USend-->

The alert resolves after the guest agent reports the updated driver versions.

<!--DS: If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
