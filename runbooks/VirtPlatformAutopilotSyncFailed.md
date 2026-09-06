# VirtPlatformAutopilotSyncFailed

## Meaning

This alert fires when virt-platform-autopilot fails to apply its desired
configuration (the "Opinionated State") to a managed resource for at least 15
minutes.

The alert is driven by the `kubevirt_autopilot_compliance_status` gauge, which
is set per managed resource (labelled by `kind`, `name`, and `namespace`):

- `1` = synced (the live resource matches the Golden State)
- `0` = drifted or the last apply failed

The alert triggers on `kubevirt_autopilot_compliance_status == 0` held for
`15m`. The delay tolerates transient API errors and slow rollouts; if the
condition persists longer, the automation is genuinely stuck and needs
attention.

## Impact

This is a `critical` alert. The affected resource is not in its desired state
and virt-platform-autopilot cannot correct the drift on its own. Platform
features that depend on the resource may be degraded or unavailable until the
condition clears.

## Diagnosis

1. Export the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get deployment virt-platform-autopilot -A \
       -o custom-columns="":.metadata.namespace | tr -d '[:space:]')"
   ```

2. Read the `kind`, `name`, and `namespace` of the affected resource from the
   alert labels.

3. Inspect the operator logs for errors relating to that resource:

   ```bash
   $ kubectl -n $NAMESPACE logs --tail=-1 -l app=virt-platform-autopilot | \
       grep -iE "forbidden|admission|webhook|conflict|not found|failed to apply"
   ```

4. Inspect the affected resource and its recent events. For namespaced resources:

   ```bash
   $ kubectl get <kind> <name> -n <namespace> -o yaml
   $ kubectl get events -n <namespace> \
       --field-selector involvedObject.kind=<kind>,involvedObject.name=<name>
   ```

   For cluster-scoped resources like `KubeDescheduler` or `MachineConfig`, omit
   `-n <namespace>` and use `-A` for events:

   ```bash
   $ kubectl get <kind> <name> -o yaml
   $ kubectl get events -A \
       --field-selector involvedObject.kind=<kind>,involvedObject.name=<name>
   ```

## Mitigation

Resolve the root cause revealed by the logs and events:

- **Admission webhook rejection:** fix the policy, or exempt the operator
  service account `system:serviceaccount:$NAMESPACE:virt-platform-autopilot`.
- **RBAC denial:** this indicates a bug in the virt-platform-autopilot
  installation. The operator should have the required permissions by default.
  Report this issue with the output of:

  ```bash
  $ kubectl auth can-i update <resource-type> -n <namespace> \
      --as system:serviceaccount:$NAMESPACE:virt-platform-autopilot
  $ kubectl get clusterrolebinding virt-platform-autopilot-rolebinding
  ```

- **Conflicting controller or user (edit war):** see
  [VirtPlatformAutopilotThrashingDetected](VirtPlatformAutopilotThrashingDetected.md).
- **Missing CRD:** see
  [VirtPlatformAutopilotDependencyMissing](VirtPlatformAutopilotDependencyMissing.md).

If the drift is intentional and the resource must be managed outside of
virt-platform-autopilot, stop the automation for it:

```bash
$ kubectl annotate <kind> <name> -n <namespace> \
    platform.kubevirt.io/mode=unmanaged --overwrite
```

If the operator appears stuck, restart it to force a fresh reconcile:

```bash
$ kubectl -n $NAMESPACE rollout restart deployment virt-platform-autopilot
```

Verify the rollout completes:

```bash
$ kubectl -n $NAMESPACE rollout status deployment virt-platform-autopilot
```

The alert resolves when `kubevirt_autopilot_compliance_status` is no longer `0`
for the affected resource: either it returns to `1` (synced), or the compliance
series is removed entirely (the resource is now unmanaged or excluded from
autopilot control).

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
