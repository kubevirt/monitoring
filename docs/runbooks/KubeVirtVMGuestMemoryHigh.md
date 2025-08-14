# KubeVirtVMGuestMemoryHigh

## Meaning

This alert fires when a Virtual Machine has been consuming more than 85% of
its available memory for more than 5 minutes, as seen from the guest OS
perspective.

## Impact

High memory usage in a VM can lead to:
- **Performance degradation** - Applications may become slow due to memory pressure
- **Memory swapping** - Guest OS starts using swap space, severely impacting performance
- **Application crashes** - Out-of-memory (OOM) conditions may kill processes
- **System instability** - Critical system processes may fail due to memory exhaustion
- **Cascading failures** - Memory pressure can affect the entire application stack

## Diagnosis

### 1. Identify the Affected VM

The alert provides the following information: VM Name, Namespace, Node and
Memory Usage.

### 2. Check VM Memory Allocation

```bash
# Get VM configuration
kubectl get vm <vm-name> -n <namespace> -o yaml

# Check current VMI status
kubectl get vmi <vm-name> -n <namespace> -o yaml
```

Look for:
- `spec.domain.memory.guest` - Guest memory allocation
- `spec.domain.resources.requests.memory` - Memory resource requests
- `spec.domain.resources.limits.memory` - Memory resource limits

### 3. Monitor Current Memory Usage

```bash
# Check VM resource consumption
virtctl top vmi <vm-name> -n <namespace>

# Get detailed memory metrics
kubectl exec -n <namespace> \
  $(kubectl get pods -n <namespace> -l kubevirt.io/vmi=<vm-name> \
    -o jsonpath='{.items[0].metadata.name}') \
  -c compute -- cat /proc/meminfo
```

### 4. Check Memory Usage Inside the VM

Connect to the VM to identify memory consumption:

```bash
# Console access
virtctl console <vm-name> -n <namespace>

# Or SSH if available
virtctl ssh <vm-name> -n <namespace>
```

Inside the VM, run:
```bash
# Overall memory usage
free -h

# Top processes by memory usage
top -o %MEM

# Detailed memory information
cat /proc/meminfo

# Check for memory-hungry processes
ps aux --sort=-%mem | head -20

# Check for memory leaks
pmap -x $(pgrep -f <process-name>)
```

### 5. Check for Memory Pressure Indicators

```bash
# Check for OOM events in VM
dmesg | grep -i "killed process"
journalctl | grep -i "out of memory"

# Check swap usage
swapon --show
free -h | grep Swap

# Check memory pressure stall information (if available)
cat /proc/pressure/memory
```

### 6. Check Node-Level Metrics

```bash
# Check node memory usage
kubectl top node <node>

# Check other VMs on the same node
kubectl get vmi --all-namespaces --field-selector spec.nodeName=<node>

# Check node conditions and pressure
kubectl describe node <node>
```

### 7. Historical Analysis

Use Prometheus/Grafana to analyze:

```promql
# VM memory usage over time
kubevirt_vmi_memory_used_bytes{name="<vm-name>", namespace="<namespace>"}

# VM memory usage percentage
(kubevirt_vmi_memory_used_bytes{name="<vm-name>", namespace="<namespace>"} / \
 kubevirt_vmi_memory_available_bytes{name="<vm-name>", namespace="<namespace>"}) * 100

# Memory allocation vs usage
kubevirt_vmi_memory_domain_bytes{name="<vm-name>", namespace="<namespace>"}

# Compare with other VMs
kubevirt_vmi_memory_used_bytes / kubevirt_vmi_memory_available_bytes
```

## Mitigation

### Immediate Actions

1. **Scale VM Memory** (Hot-plug memory - no downtime):
   ```bash
   # Edit VM to increase memory allocation
   kubectl edit vm <vm-name> -n <namespace>
   ```
   - If using instancetypes: Change `spec.instancetype` to a larger
     instancetype with more memory
   - If not using instancetypes: Modify `spec.template.spec.domain.memory.guest`
     or `spec.template.spec.domain.resources`
   
   ```bash
   # Verify the change was applied
   kubectl get vmi <vm-name> -n <namespace> \
     -o jsonpath='{.spec.domain.memory.guest}'
   ```

   **Note**: Memory hot-plug is supported in KubeVirt and allows scaling
   without VM restart. The VM guest OS should recognize the new memory
   automatically.

   **Verify hot-plug inside the VM**:
   ```bash
   # Connect to the VM and check memory
   virtctl console <vm-name> -n <namespace>
   
   # Inside the VM, verify new memory is available
   free -h                 # Show current memory usage
   cat /proc/meminfo       # Detailed memory information
   ```

2. **Free Memory Inside VM**:
   - Restart memory-intensive processes or services
   - Clear caches: `sync && echo 3 > /proc/sys/vm/drop_caches` (as root)
   - Identify processes with unusually high memory usage
   - Investigate if these processes are legitimate or problematic
   - Only terminate processes after confirming they are safe to kill and
     not critical system processes

### Long-term Solutions

1. **Right-size VM Memory**:
   - Analyze historical memory usage patterns
   - Adjust memory requests and limits based on actual needs
   - Consider memory overcommitment ratios

2. **Application Optimization**:
   - Review application memory usage and optimize
   - Implement proper garbage collection for applications
   - Configure application memory limits appropriately
   - Consider memory-efficient alternatives

3. **Memory Management**:
   - Configure swap appropriately (or disable if not needed)
   - Tune VM memory parameters (balloon driver, huge pages)
   - Implement memory monitoring within applications

4. **Monitoring and Alerting**:
   - Set up additional memory-specific alerts (swap usage, OOM events)
   - Monitor memory growth trends over time

5. **Capacity Planning**:
   - Review cluster memory capacity
   - Plan for node scaling or hardware upgrades
   - Use Kubernetes node affinity/anti-affinity rules to distribute
     memory-intensive VMs across different nodes
   - Implement resource quotas and limits

## Prevention

- **Proper memory allocation**: Set appropriate memory requests and limits
  based on application requirements
- **Regular monitoring**: Implement comprehensive memory monitoring and
  capacity planning
- **Load testing**: Test applications under expected memory load before
  production deployment
- **Resource quotas**: Use namespace resource quotas to prevent memory
  over-allocation
- **Memory profiling**: Regularly profile applications to identify memory
  leaks and optimization opportunities
- **Gradual scaling**: Implement proper scaling policies based on memory usage trends

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->

If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](
  https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
