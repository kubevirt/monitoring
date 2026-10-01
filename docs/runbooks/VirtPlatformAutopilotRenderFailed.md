# VirtPlatformAutopilotRenderFailed

## Meaning

This alert fires when virt-platform-autopilot cannot render an active asset's
desired configuration for at least five minutes.

The alert is driven by the `kubevirt_autopilot_render_failed` gauge, labelled
by `asset`:

- `1` = the last rendering attempt failed
- `0` = the last rendering attempt succeeded, including an intentionally empty
  template

The alert triggers on `kubevirt_autopilot_render_failed == 1` held for `5m`.
It detects failures before a desired resource is available, which are not
reliably covered by
[VirtPlatformAutopilotSyncFailed](VirtPlatformAutopilotSyncFailed.md).
The alert's `namespace` label identifies the operator installation namespace,
not the namespace of a resource that the asset creates.

## Impact

This is a `warning` alert. The affected asset is not created or updated to its
desired configuration. A dependent platform feature can remain unavailable
or retain its previous configuration. Other eligible assets continue to
reconcile, unless the failure affects the initial HyperConverged asset,
which blocks subsequent asset reconciliation.

## Diagnosis

1. Export the `NAMESPACE` environment variable using the `namespace` alert
   label, and set `ASSET` to the `asset` alert label:

   ```bash
   $ export NAMESPACE="<namespace>"
   $ export ASSET="<asset>"
   ```

2. Inspect the operator logs for the affected asset and its rendering error:

   ```bash
   $ kubectl -n "$NAMESPACE" logs --tail=500 \
       -l app=virt-platform-autopilot | grep -F "$ASSET"
   ```

   Look for `failed to render asset`, including the underlying error. A
   rendering failure can originate from loading a template, evaluating it,
   querying cluster configuration, or parsing the rendered YAML.

3. Inspect the HyperConverged configuration and recent events. Set `HCO` and
   `HCO_NAMESPACE` to the name and namespace shown by the first command:

   ```bash
   $ kubectl get hyperconverged -A
   $ export HCO="<hco-name>"
   $ export HCO_NAMESPACE="<hco-namespace>"
   $ kubectl -n "$HCO_NAMESPACE" get hyperconverged "$HCO" -o yaml
   $ kubectl -n "$HCO_NAMESPACE" get events \
       --field-selector involvedObject.kind=HyperConverged,involvedObject.name="$HCO" \
       --sort-by=.lastTimestamp
   ```

   Versions that emit `RenderFailed` events on the HyperConverged resource
   include the asset name and underlying error there. Otherwise, use the
   operator logs.

4. If the error refers to logging storage, inspect the available storage
   classes and the existing LokiStack:

   ```bash
   $ kubectl get storageclass
   $ kubectl -n openshift-logging get lokistack logging-loki \
       -o jsonpath='{.spec.storageClassName}{"\n"}'
   ```

   The S3 Secret is a separate requirement for Loki readiness. Do not print
   Secret data or include credentials in diagnostic artifacts.

## Mitigation

Resolve the underlying error identified in the logs:

- **Invalid logging storage-class override:** choose an existing storage
  class that supports Loki's block storage requirements, then correct the
  HyperConverged annotation. These commands apply to versions that support
  the `platform.kubevirt.io/logging-storage-class` annotation:

  ```bash
  $ kubectl -n "$HCO_NAMESPACE" annotate hyperconverged "$HCO" \
      platform.kubevirt.io/logging-storage-class=<storage-class> --overwrite
  ```

  To remove an unintended override and restore automatic selection:

  ```bash
  $ kubectl -n "$HCO_NAMESPACE" annotate hyperconverged "$HCO" \
      platform.kubevirt.io/logging-storage-class-
  ```

  Automatic selection preserves the existing LokiStack storage class, or
  uses the cluster default when no existing class is configured. If neither
  exists, select a suitable class with the annotation or ask the storage
  administrator to configure a cluster default. Changing a storage class
  does not migrate existing PVCs; do not delete LokiStack resources or PVCs
  to clear this alert.

- **API read denied (`forbidden`):** verify the exact denied verb and resource
  using the operator service account. For a cluster-scoped resource:

  ```bash
  $ kubectl auth can-i <verb> <resource-type> \
      --as="system:serviceaccount:$NAMESPACE:virt-platform-autopilot"
  ```

  Add `-n <resource-namespace>` for a namespaced resource. Restore the
  supported installation RBAC or complete the operator upgrade if the image
  and permissions are out of sync. Do not grant `cluster-admin` as a
  workaround. Impersonation checks require administrator permissions.

- **Temporary API lookup failure:** restore API availability or connectivity.
  Autopilot retries reconciliation; no resource deletion is necessary.

- **Template or rendered YAML error:** report the asset name, operator
  version, and sanitized error to the maintainers. Install a version with
  the corrected template. Do not edit managed operands to bypass rendering.

The alert clears on the next evaluation after rendering succeeds and the
gauge returns to `0`, or after the asset is excluded and its series is
removed. Rendering success does not guarantee that applying the resource
succeeds; investigate any remaining sync-failure alert separately.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
