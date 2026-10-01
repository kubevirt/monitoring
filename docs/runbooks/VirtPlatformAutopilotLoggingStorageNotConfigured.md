# VirtPlatformAutopilotLoggingStorageNotConfigured

## Meaning

Autopilot cannot select a StorageClass for the `logging-lokistack` asset for
at least fifteen minutes. Logging is enabled, but no explicit storage class,
existing LokiStack storage class, or cluster default StorageClass is available.

The warning alert selects
`kubevirt_autopilot_render_failed{asset="logging-lokistack", reason="configuration", code="NoDefaultStorageClass"} == 1`
for `15m`. The `namespace` label identifies the operator installation namespace.
The fifteen-minute delay allows storage setup to complete before notifying an
administrator. Internal template errors and API read failures do not trigger
this alert.

## Impact

Autopilot cannot create or update the LokiStack desired configuration.
Logging setup remains incomplete. Other eligible assets continue to reconcile.

This alert does not detect missing explicitly named storage classes, unsuitable
storage classes, object-storage credentials, or PVC provisioning failures.
Those problems are reported by Loki and the storage provisioner.

## Diagnosis

1. Set the HyperConverged name expected by Autopilot and discover its namespace:

   ```bash
   $ HCO="kubevirt-hyperconverged"
   $ HCO_NAMESPACE=$(kubectl get hyperconverged -A \
       --field-selector metadata.name="$HCO" \
       -o jsonpath='{.items[*].metadata.namespace}')
   $ printf 'HyperConverged: %s\nNamespace: %s\n' "$HCO" "$HCO_NAMESPACE"
   ```

   Continue only if the lookup succeeds and returns exactly one namespace.
   If it returns no namespace or multiple namespaces, inspect
   `kubectl get hyperconverged -A` and identify the resource managed by the
   affected Autopilot installation before continuing. Do not assume that the
   alert's operator namespace is also the HyperConverged namespace.
   Run the remaining commands in the same shell; `export` is unnecessary.

2. Read the warning events for the exact configuration problem:

   ```bash
   $ kubectl -n "$HCO_NAMESPACE" get events \
       --field-selector involvedObject.kind=HyperConverged,involvedObject.name="$HCO",reason=RenderFailed \
       --sort-by=.lastTimestamp
   ```

   The event identifies `logging-lokistack` and `NoDefaultStorageClass`.

3. Inspect storage classes, the HCO override, and any existing LokiStack:

   ```bash
   $ kubectl get storageclass
   $ kubectl -n "$HCO_NAMESPACE" get hyperconverged "$HCO" \
       -o go-template='{{index .metadata.annotations "platform.kubevirt.io/logging-storage-class"}}{{"\n"}}'
   $ kubectl -n openshift-logging get lokistack logging-loki \
       -o jsonpath='{.spec.storageClassName}{"\n"}'
   ```

   A new installation can have no LokiStack yet. A configured existing stack
   keeps its storage class even when the cluster default changes.

## Mitigation

Select an available StorageClass suitable for Loki's filesystem PVCs, then set
it explicitly on HyperConverged:

```bash
$ kubectl -n "$HCO_NAMESPACE" annotate hyperconverged "$HCO" \
    platform.kubevirt.io/logging-storage-class=<storage-class> --overwrite
```

Alternatively, ask the storage administrator to configure a cluster default:

```bash
$ kubectl annotate storageclass <storage-class> \
    storageclass.kubernetes.io/is-default-class=true --overwrite
```

Changing the default affects other workloads that do not specify a storage
class. Use the HCO annotation when the choice should apply only to logging.

Autopilot retries reconciliation automatically. The alert clears after rendering
succeeds and the failure series is removed, or after logging is excluded.
Check LokiStack and PVC readiness separately to confirm successful setup.

Changing the selection does not migrate existing PVCs. Preserve an existing
stack's storage class unless performing a planned storage migration; do not
delete LokiStack resources or PVCs to clear this alert.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
