# KubeVirtVMHighCPUUsage

## Meaning

This alert fires when a Virtual Machine (VM) has been consuming more than
80% of its allocated CPU resources for more than 5 minutes.

## Impact

High CPU usage in a VM can indicate:
- **Performance degradation** - Applications running inside the VM may
  become slow or unresponsive
- **Potential application issues** - Runaway processes, infinite loops, or
  inefficient code
- **Capacity planning needs** - The VM may need additional CPU resources
- **Workload growth** - Legitimate increase in application demand
  requiring more CPU

## Diagnosis

### 1. Identify the Affected VM

The alert provides the VM Name, Namespace, Node and CPU Usage.

### 2. Check VM CPU Allocation

```bash
# Get VM configuration
kubectl get vm <vm-name> -n <namespace> -o yaml

# Check current VMI status
kubectl get vmi <vm-name> -n <namespace> -o yaml
```

Look for:
- `spec.domain.cpu.cores` - Number of allocated vCPUs
- `spec.domain.resources.requests.cpu` - CPU resource requests
- `spec.domain.resources.limits.cpu` - CPU resource limits

### 3. Monitor Current CPU Usage and Processes

Connect to the VM to identify high CPU processes:

```bash
# Console access
virtctl console <vm-name> -n <namespace>

# Or SSH if available
virtctl ssh vm/<vm-name> -n <namespace>
```

Inside the VM, run:
```bash
# Top processes by CPU usage
top -c

# Process tree
htop

# Detailed process information
ps aux --sort=-%cpu | head -20
```

### 4. Check Node-Level Metrics

```bash
# Check node resource usage
kubectl top node <node>

# Check other VMs on the same node
kubectl get vmi --all-namespaces \
  --field-selector spec.nodeName=<node>

# Check node conditions
kubectl describe node <node>
```

## Mitigation

### Immediate Actions

1. **Scale VM Resources** (requires VM restart):
   ```bash
   # Edit VM to increase CPU allocation
   kubectl edit vm <vm-name> -n <namespace>
   ```
   - If using instancetypes: Change `spec.instancetype` to a larger
     instancetype with more CPU cores
   - If not using instancetypes: Modify
     `spec.template.spec.domain.cpu.cores` or
     `spec.template.spec.domain.resources`

   **Note**: Changing CPU allocation requires stopping and restarting
   the VM.

2. **Identify and Address High-CPU Processes**:
   - Connect to the VM and identify runaway processes or processes with
     unusually high CPU load
   - Investigate if these processes are legitimate or problematic
     (infinite loops, stuck processes, etc.)
   - Only terminate processes after confirming they are safe to kill and
     not critical system processes

### Long-term Solutions

1. **Right-size VM Resources**:
   - Analyze historical usage patterns
   - Adjust CPU requests and limits appropriately

2. **Application Optimization**:
   - Review application performance and efficiency
   - Consider migrating CPU-intensive workloads to dedicated nodes

3. **Monitoring and Alerting**:
   - Set up additional monitoring for application-specific metrics

4. **Capacity Planning**:
   - Review cluster CPU capacity
   - Plan for node scaling or hardware upgrades
   - Consider CPU affinity and anti-affinity rules

## Prevention

- **Proper resource allocation**: Set appropriate CPU requests and limits
  based on workload requirements
- **Regular monitoring**: Implement comprehensive monitoring and capacity
  planning
- **Performance testing**: Test applications under expected load before
  production deployment
- **Resource quotas**: Use namespace resource quotas to prevent resource
  over-allocation
- **Strategic VM placement**: Use Kubernetes node affinity/anti-affinity
  rules to distribute CPU-intensive VMs across different nodes

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->

If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](
  https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
