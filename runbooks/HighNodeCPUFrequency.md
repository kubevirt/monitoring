# HighNodeCPUFrequency

## Meaning

This alert fires when a CPU frequency on a node exceeds 80% of the maximum
frequency for more than 5 minutes.

## Impact

High CPU frequency can indicate:
- High CPU utilization pushing the processor to higher frequencies
- Potential thermal issues causing frequency scaling
- Power management concerns affecting system stability
- Reduced CPU lifespan due to sustained high-frequency operation

## Diagnosis

1. Identify the affected node and CPU:
   ```bash
   kubectl get nodes
   ```

2. Check current CPU frequency on the node:
   ```bash
   kubectl debug node/<node-name> -it --image=registry.redhat.io/ubi8/ubi
   ```

   Then run inside the debug pod:
   ```bash
   cat /proc/cpuinfo | grep -i "cpu mhz"
   ```

3. Monitor CPU utilization and temperature:
   ```bash
   kubectl top nodes
   ```

   ```bash
   kubectl top pods --all-namespaces --sort-by=cpu
   ```

   Check system temperature (if available):
   ```bash
   sensors
   ```

4. Review node resource allocation:
   ```bash
   kubectl describe node <node-name>
   ```

5. Check for CPU-intensive workloads:
   ```bash
   ps aux --sort=-%cpu | head -20
   ```

## Mitigation

1. Immediate actions:
   - Monitor the CPU temperature to ensure it's within safe limits
   - Check if the high frequency is due to legitimate high CPU demand
   - Verify CPU settings if needed

2. If caused by high CPU utilization:
   - Identify and analyze CPU-intensive pods
   - Consider redistributing workloads across nodes
   - Evaluate pod resource limits and requests
   - Scale horizontally if needed

3. If thermal throttling is suspected:
   - Check system cooling and ventilation
   - Monitor ambient temperature
   - Consider reducing workload temporarily
   - Verify proper thermal management settings

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
