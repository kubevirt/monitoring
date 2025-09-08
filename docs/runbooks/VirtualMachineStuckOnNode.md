# VirtualMachineStuckOnNode

## Meaning

This alert triggers when a VirtualMachine (VM) with an associated
VirtualMachineInstance (VMI) has been stuck in an unhealthy state for more than
5 minutes on a specific node.

The alert indicates that the VM has progressed past initial scheduling and has
an active VMI, but is experiencing runtime issues on the assigned node. This
typically occurs after the VM has been scheduled to a node but encounters
problems during startup, operation, or shutdown phases.

**Affected States:**
- `Starting` - VMI exists but VM is failing to reach running state
- `Stopping` - VM is attempting to stop but the process is stuck
- `Terminating` - VM is being deleted but the termination process is hanging
- `Error` states - Runtime errors occurring on the node (ErrImagePull,
  ImagePullBackOff, etc.)

## Impact

- **Severity:** Warning
- **User Impact:** VMs are unresponsive or stuck in transition states
- **Business Impact:** Running workloads may be disrupted, which affects
application performance and availability
- **Node Impact:** Resources may be tied up by unresponsive VMs, which affects
other workloads on the same node

## Possible Causes

### Node-Level Issues
- **Node resource exhaustion** (CPU, memory, storage)
- **Container runtime problems** (issues with containerd or CRI-O)
- **Insufficient storage** on the node
- **Network connectivity issues** from the node
- **Node entering the NotReady state** or maintenance mode

### QEMU/KVM Issues
- **QEMU process failures** or hangs
- **KVM acceleration problems** on the node
- **Nested virtualization** configuration issues
- **Hardware compatibility** problems

### Image and Storage Issues
- **Container image pull failures** specific to the node
- **Local image cache corruption**
- **PVC mount failures** on the node
- **Storage backend connectivity** issues from the node
- **Volume attachment timeouts**

### virt-launcher Pod Issues
- **virt-launcher pod** stuck in non-ready state
- **Pod resource limits** being exceeded
- **Security policy violations** (SELinux, AppArmor)
- **Networking problems** within the pod

### libvirt/Domain Issues
- **libvirt daemon** problems on the node
- **Domain definition** conflicts or corruption
- **Migration failures** (if VM was being migrated)
- **Hot-plug operations** that failed and left VM in inconsistent
  state

## Diagnosis

### 1. Check VM and VMI Status
```bash
# Get VM details with node information
kubectl get vm <vm-name> -n <namespace> -o yaml

# Check VMI status and node assignment
kubectl get vmi <vm-name> -n <namespace> -o yaml
kubectl describe vmi <vm-name> -n <namespace>

# Look for related events
kubectl get events -n <namespace> \
  --field-selector involvedObject.name=<vm-name>
```

### 2. Examine virt-launcher Pod
```bash
# Find the virt-launcher pod for this VM
kubectl get pods -n <namespace> -l kubevirt.io/domain=<vm-name>

# Check pod status and events
kubectl describe pod <virt-launcher-pod> -n <namespace>

# Check pod logs for errors
kubectl logs <virt-launcher-pod> -n <namespace> -c compute
kubectl logs <virt-launcher-pod> -n <namespace> -c istio-proxy \
  # if using Istio

# Optional: Check resource usage for the virt-launcher pod
kubectl top pod <virt-launcher-pod> -n <namespace>
```

### 3. Investigate Node Health
```bash
# Check node status and conditions (may require admin
# permissions)
kubectl describe node <node-name>

# Discover the KubeVirt installation namespace
export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"

# Check virt-handler on the affected node
kubectl get pods -n "$NAMESPACE" -o wide | grep <node-name>
kubectl logs <virt-handler-pod> -n "$NAMESPACE"
```

### 4. Check Storage and Volumes
```bash
# Verify PVC status and mounting
kubectl get pvc -n <namespace>
kubectl describe pvc <pvc-name> -n <namespace>

# Check volume attachments on the node
kubectl get volumeattachment | grep <node-name>

# For DataVolumes, check their status
kubectl get dv -n <namespace>
kubectl describe dv <dv-name> -n <namespace>
```

### 5. Verify Image Accessibility from Node
```bash
# Verify image accessibility from the affected node
kubectl debug node/<node-name> -it --image=busybox

# Inside the debug pod, check which container runtime is used
ps aux | grep -E "(containerd|dockerd|crio)"

# For CRI-O/containerd clusters:
crictl pull <vm-disk-image>

# For Docker-based clusters (less common):
docker pull <vm-disk-image>

# Exit the debug session when done
exit
```

### 6. Exec into the virt‑launcher pod’s compute container and inspect domains
```bash
kubectl exec -it <virt-launcher-pod> -n <namespace> -c compute \
  -- virsh list --all | grep <vm-name>
kubectl exec -it <virt-launcher-pod> -n <namespace> -c compute \
  -- virsh dumpxml <domain-name>

```

## Mitigation

### Pod-Level Issues
1. **Restart the virt-launcher pod**:
   ```bash
   kubectl delete pod <virt-launcher-pod> -n <namespace>
   # The VMI controller will recreate it
   ```

2. **Check resource constraints**:
   ```bash
   kubectl describe pod <virt-launcher-pod> -n <namespace>
   # Look for resource limit violations
   ```

### Image Issues on Node
1. **Inspect and, if necessary, clear image cache** on the node:
   ```bash
   # SSH to the node or start a debug session on the node:
   kubectl debug node/<node-name> -it --image=busybox

   # Detect which container runtime is in use
   ps aux | grep -E "(containerd|dockerd|crio)"

   # List cached images first
   # For CRI-O/containerd clusters:
   crictl images
   # For Docker-based clusters:
   docker images

   # Remove only if a corrupted/stale image is suspected
   # For CRI-O/containerd clusters:
   crictl rmi <problematic-image>
   # For Docker-based clusters:
   docker rmi <problematic-image>

   exit
   ```

2. **Force image re-pull**:
   ```bash
   # Delete and recreate the virt-launcher pod
   kubectl delete pod <virt-launcher-pod> -n <namespace>
   ```

### Storage Issues
1. **Check PVC binding and mounting**:
   ```bash
   kubectl get pvc -n <namespace>
   # If PVC is stuck, check the storage provisioner
   ```

2. **Resolve volume attachment issues**:
   ```bash
   kubectl get volumeattachment
   # Delete stuck volume attachments if necessary
   kubectl delete volumeattachment <attachment-name>
   ```

### Node-Level Issues Resolution
1. **Drain and uncordon the node** if it is in a bad state:
   ```bash
   kubectl drain <node-name> --ignore-daemonsets \
     --delete-emptydir-data
   kubectl uncordon <node-name>
   ```

2. **Restart node-level components**:
   ```bash
   # Restart virt-handler on the node
   kubectl delete pod <virt-handler-pod> -n "$NAMESPACE"
   ```

### VM-Level Resolution
1. **Force‑delete the VMI (triggers creating a new VMI)**:
   ```bash
   kubectl delete vmi <vm-name> -n <namespace> --force \
     --grace-period=0
   ```

2. **Migrate the VM to a different node**:
   ```bash
   virtctl migrate <vm-name> -n <namespace>
   ```

### Emergency Actions
- **Live migrate** critical VMs away from the problematic node
- **Force delete** the unresponsive VMI if safe to do so:
  ```bash
  kubectl delete vmi <vm-name> -n <namespace> --force --grace-period=0
  ```
- **Cordon the node** to prevent new VM scheduling while investigating

## Prevention

1. **Node Health Monitoring:**
   - Monitor node resource utilization (CPU, memory,
     storage)
   - Set up alerts for node conditions and taints
   - Perform regular health checks on container runtime

2. **Resource Management:**
   - Set appropriate resource requests/limits on VMs
   - Monitor PVC and storage utilization
   - Plan for node capacity and VM density

3. **Image Management:**
   - Use image pull policies appropriately (Always,
     IfNotPresent)
   - Pre-pull critical images to nodes
   - Monitor image registry health and connectivity

4. **Networking:**
   - Ensure stable network connectivity between nodes and storage
   - Monitor DNS resolution and service discovery
   - Validate network policies do not block required traffic

5. **Regular Maintenance:**
   - Keep nodes and KubeVirt components updated

## Escalation

Escalate to the cluster administrator if:
- Multiple VMs are affected on the same node simultaneously
- VMs consistently fail to start or become unresponsive after
  following troubleshooting steps
- Node‑specific issues persist (Like kubelet or kernel panic) and
  require a node reboot
- A malfunction of KubeVirt components is suspected
- You are unable to access system logs for further diagnosis
- You do not have enough permissions to run the diagnosis and/or mitigation
  steps.

## Related Alerts

- `OrphanedVirtualMachineInstances` - May indicate virt-handler
  problems on nodes
- `VirtHandlerDown` - Related to virt-handler pod failures
- `VirtualMachineStuckInUnhealthyState` - For VMs that haven't
  progressed to having VMIs

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support
case, attaching the artifacts gathered during the diagnosis
procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](
  https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
