# Runbook Review Guidelines

These guidelines define how to review KubeVirt alert runbooks in
`docs/runbooks/`. The key principle: **every claim in a runbook must be
verified against the actual alert source code** in the kubevirt
organization repositories.

## Alert Source Code Location

Alert definitions do NOT live in this repo. They are defined in operator
repositories across the kubevirt GitHub organization. Alert names do NOT
always contain the operator name as a prefix — you cannot guess the repo
from the alert name. Search ALL of these repos for the exact alert name:

- `kubevirt/kubevirt`
- `kubevirt/containerized-data-importer`
- `kubevirt/hyperconverged-cluster-operator`
- `kubevirt/ssp-operator`
- `kubevirt/cluster-network-addons-operator`
- `kubevirt/hostpath-provisioner-operator`
- `kubevirt/hostpath-provisioner`
- `k8snetworkplumbingwg/kubemacpool`
- `kubevirt/application-aware-quota`

## Alert May Only Exist in an Open PR

Often the alert is not yet merged in the source repository — it only
exists in an open pull request. This is the normal workflow: the alert
PR and the runbook PR are created in parallel.

When reviewing a runbook:
1. First search the main branch of the source repo for the alert name
2. If not found, search **open pull requests** in the source repo
3. The runbook PR description should link to the corresponding alert
   PR — follow that link to find the alert definition
4. If no alert definition can be found anywhere (neither merged nor in
   an open PR), flag this clearly: the runbook cannot be verified

## What to Verify

### 1. Find the alert definition

For each runbook, locate the alert by its name in the source repo
(main branch or open PR). An alert definition includes:
- **`expr`**: The PromQL expression that triggers the alert
- **`for`**: How long the condition must hold before the alert fires
- **`severity`**: critical, warning, or info
- **`operator_health_impact`**: critical, warning, or none
- **`summary`** and **`description`** annotations

### 2. Trace the PromQL expression end-to-end

Break down the `expr` to understand what it actually evaluates:

- **Metrics**: Find where each metric in the expression is defined in
  the Go source code. Understand what it measures — is it a counter,
  gauge, histogram? What events cause it to increment or change?
- **Recording rules**: If the expression references recording rules
  (metrics that don't exist as raw instrumentation), find the recording
  rule definition and trace what it computes.
- **Thresholds and conditions**: What numeric thresholds or boolean
  conditions trigger the alert? Over what time window?

The goal is to fully understand: what system state causes this alert
to fire?

### 3. Verify each runbook section

#### `## Meaning`

- Must accurately describe what the PromQL expression evaluates
- Must state the correct firing condition (threshold, duration)
- Must reference the correct component — if the expression queries
  `virt-controller` pods, the runbook must not say `virt-api`
- The `for` duration should match the alert definition

#### `## Impact`

- Based on where the metrics are set in source code, what functionality
  is actually affected when this condition is true?
- Does the runbook capture the real user-facing impact, or is it
  vague/generic?

#### `## Diagnosis`

- `kubectl`/`oc` commands must target the correct resources and labels
  for the component the alert monitors
- Label selectors must match what the source code uses (e.g.,
  `kubevirt.io=virt-controller`)
- Commands should help verify the specific condition the PromQL
  expression checks
- Commands must be syntactically valid
- Use `$NAMESPACE` variable pattern:
  ```bash
  $ export NAMESPACE="$(kubectl get kubevirt -A \
      -o custom-columns="":.metadata.namespace)"
  ```

#### `## Mitigation`

- Remediation steps should address the root causes that the source
  code reveals
- Should be concrete and actionable, not just "investigate the issue"

## Formatting Rules

- One runbook per file: `<AlertName>.md`
- H1 heading must match filename: `VirtAPIDown.md` → `# VirtAPIDown`
- Required H2 sections in order: `## Meaning`, `## Impact`,
  `## Diagnosis`, `## Mitigation`
- Code fences: ```` ```bash ```` (no space before language tag)
- Line length: 80 characters for prose (code blocks exempt)
- Prefix example commands with `$ `
- `<!--DS: ... -->` and `<!--USstart-->...<!--USend-->` are intentional
  downstream/upstream content markers — do not flag them
