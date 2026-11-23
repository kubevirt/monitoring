# VirtPlatformAutopilotMachineConfigUpdateStaged

## Meaning

> **Note:** This is an informational alert with no operator health impact.

Autopilot has prepared a node configuration update (a MachineConfig) and holds
it instead of applying it. The alert fires when an update has been held for at
least five minutes while Autopilot waits for a matching MachineConfigPool
rollout; it selects `kubevirt_autopilot_machineconfig_update_staged == 1` for
`5m`. The `machineconfig` and `pool` labels identify the staged update and its
matching pool.

Autopilot stages the update to coalesce configuration changes with an
existing rollout and reduce node reboots: the live MachineConfig keeps its
previous configuration while the update is held. Coalescing is best-effort,
not a guarantee of a single rollout: a matching pool permits an apply attempt
when it reports `Updating=True`, which can occur after the rollout is already
underway, or `RenderDegraded=True`, which means the Machine Config Operator
cannot render the pool's configuration at all, so holding the update has no
benefit. Nodes that the rollout has already updated may need to converge again
once the staged configuration is applied.

## Impact

Nodes keep their previous configuration until the staged update is released
and applied. The intended configuration change is delayed, not lost, and no
administrator action is required while the update is held.

## Diagnosis

Inspect the named MachineConfigPool and MachineConfig:

```bash
$ kubectl get machineconfigpool <pool> -o yaml
$ kubectl get machineconfig <machineconfig> -o yaml
```

Check whether the pool is updating or paused. The informational alert alone
does not mean that the Machine Config Operator is unhealthy.

## Mitigation

No administrator action is required for normal staging. Autopilot releases
the update automatically once the matching pool permits the apply. Follow the
existing maintenance plan for a deliberately paused pool; do not start a
rollout solely to clear this informational alert.

If the pool is degraded or a rollout is stuck, investigate the Machine Config
Operator's conditions and alerts. The staged-update alert clears when
Autopilot applies or discards the staged update and removes its metric series.

<!--DS: If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
