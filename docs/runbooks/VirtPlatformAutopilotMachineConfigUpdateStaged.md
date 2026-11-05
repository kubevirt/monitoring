# VirtPlatformAutopilotMachineConfigUpdateStaged

## Meaning

Autopilot holds a desired MachineConfig update for at least five minutes while
waiting for a matching MachineConfigPool rollout. The alert selects
`kubevirt_autopilot_machineconfig_update_staged == 1` for `5m`.
The `machineconfig` and `pool` labels identify the staged update and matching
pool. This is an informational alert with no operator health impact.

## Impact

The live MachineConfig retains its previous configuration until the staged
update is released. This is intentional: Autopilot coalesces changes with an
existing rollout to avoid starting an additional rollout solely for its update.

## Diagnosis

Inspect the named MachineConfigPool and MachineConfig:

```bash
$ kubectl get machineconfigpool <pool> -o yaml
$ kubectl get machineconfig <machineconfig> -o yaml
```

Check whether the pool is updating or paused. The informational alert alone
does not mean that the Machine Config Operator is unhealthy.

## Mitigation

No administrator action is required for normal staging. Autopilot releases the
update when a matching pool starts a rollout. Follow the existing maintenance
plan for a deliberately paused pool; do not start a rollout solely to clear
this informational alert.

If the pool is degraded or a rollout is stuck, investigate the Machine Config
Operator's conditions and alerts. The staged-update alert clears when Autopilot
applies or discards the staged update and removes its metric series.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
