# VirtPlatformAutopilotLoggingStorageNotConfigured

## Meaning

Logging is enabled, but its storage is not configured. Autopilot cannot find
a StorageClass for the LokiStack it manages: the HyperConverged resource does
not select one for logging, an existing LokiStack does not specify one, and
the cluster has no default StorageClass.

The alert fires when this condition holds for fifteen minutes:

```promql
kubevirt_autopilot_render_failed{asset="logging-lokistack", reason="configuration", code="NoDefaultStorageClass"} == 1
```

## Impact

Logging setup remains incomplete. The LokiStack is not created or updated, so
logs are not collected and saved: application, infrastructure, and, when
enabled, audit logs are discarded when their pods terminate. Other
Autopilot-managed assets continue to reconcile.

This alert reports only the missing storage-class selection. It does not
detect missing explicitly named storage classes, unsuitable storage classes,
object-storage credentials, or PVC provisioning failures; those problems are
reported by Loki and the storage provisioner. Internal rendering errors and
API read failures do not trigger this alert.

## Diagnosis

1. Set the HyperConverged name expected by Autopilot and discover its
   namespace:

   ```bash
   $ HCO="kubevirt-hyperconverged"
   $ HCO_NAMESPACE=$(kubectl get hyperconverged -A \
       --field-selector metadata.name="$HCO" \
       -o jsonpath='{.items[*].metadata.namespace}')
   $ printf 'HyperConverged: %s\nNamespace: %s\n' "$HCO" "$HCO_NAMESPACE"
   ```

   > **Note:** If no namespace or multiple namespaces are returned, inspect
   > `kubectl get hyperconverged -A` and identify the resource managed by the
   > affected Autopilot installation before continuing. Do not assume that
   > the Autopilot installation namespace is also the HyperConverged
   > namespace.

2. Read the warning events on HyperConverged for the configuration problem:

   ```bash
   $ kubectl -n "$HCO_NAMESPACE" get events \
       --field-selector involvedObject.kind=HyperConverged,involvedObject.name="$HCO",reason=RenderFailed \
       --sort-by=.lastTimestamp
   ```

   Each `RenderFailed` event names the asset that failed to render and the
   problem code. For this alert, the event reports the `logging-lokistack`
   asset with the `NoDefaultStorageClass` code. Events for other assets or
   with other codes indicate a different problem.

3. Check where a StorageClass can come from. Autopilot selects the logging
   storage class, in order, from the HyperConverged annotation, an existing
   LokiStack, and the cluster default:

   ```bash
   $ kubectl -n "$HCO_NAMESPACE" get hyperconverged "$HCO" \
       -o go-template='{{index .metadata.annotations "platform.kubevirt.io/logging-storage-class"}}{{"\n"}}'
   $ kubectl -n openshift-logging get lokistack logging-loki \
       -o jsonpath='{.spec.storageClassName}{"\n"}'
   $ kubectl get storageclass
   ```

   The first command prints the storage class selected for logging, if any.
   The second prints the storage class of an existing LokiStack; a new
   installation has no LokiStack yet, and a configured stack keeps its
   storage class even when the cluster default changes. The third lists all
   storage classes and marks the cluster default with `(default)`.

## Mitigation

Select a StorageClass that provisions persistent volumes for Loki's internal
filesystem storage, then set it explicitly on HyperConverged:

```bash
$ kubectl -n "$HCO_NAMESPACE" annotate hyperconverged "$HCO" \
    platform.kubevirt.io/logging-storage-class=<storage-class> --overwrite
```

The selected class must support the filesystem PVCs that Loki requests. This
is block storage for Loki's internal volumes, separate from the object
storage configured in the `logging-loki-storage` secret.
<!--DS: See
[Loki block storage](https://github.com/openshift-virtualization/virt-platform-autopilot/blob/main/docs/logging.md#loki-block-storage)
for the selection policy and requirements.-->

Alternatively, ask the storage administrator to configure a cluster default.
Changing the default affects other workloads that do not specify a storage
class; use the HyperConverged annotation when the choice should apply only
to logging:

```bash
$ kubectl annotate storageclass <storage-class> \
    storageclass.kubernetes.io/is-default-class=true --overwrite
```

Autopilot retries reconciliation automatically. The metric is a status gauge:
when rendering succeeds, or when logging is excluded or set to unmanaged, the
series returns to `0` (it persists rather than disappearing), and the alert
resolves. Check LokiStack and PVC readiness separately to confirm successful
setup.

Changing the selection does not migrate existing PVCs. Preserve an existing
stack's storage class unless performing a planned storage migration; do not
delete LokiStack resources or PVCs to clear this alert.

<!--DS: If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
