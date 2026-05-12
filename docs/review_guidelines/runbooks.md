# Runbook Review Guidelines

Applies to: `docs/runbooks/**`

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
- `kubevirt/application-aware-quota`
- `k8snetworkplumbingwg/kubemacpool`

## Alert May Only Exist in an Open PR

Often the alert is not yet merged in the source repository — it only
exists in an open pull request. This is the normal workflow: the alert
PR and the runbook PR are created in parallel.

When reviewing a runbook:
1. First search the main branch of the source repo for the alert name
2. If not found, search **open pull requests** in the source repo
3. The PR description should link to the corresponding alert PR —
   follow that link to find the alert definition
4. If no alert definition can be found anywhere (neither merged nor in
   an open PR), **post a comment asking the author to add the alert PR
   link to the PR description** (or provide it in a comment reply).
   The runbook cannot be verified without it.
5. If the alert PR link is provided but the PR is not yet merged,
   note this in the review — the runbook should not be merged before
   the alert PR is merged or at least approved.

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

- Must reflect the alert's severity level:
  - `critical` alerts indicate complete unavailability or data loss
    risk — the Impact section must reflect urgency
  - `warning` alerts indicate degraded performance or reduced
    redundancy — the Impact section should reflect this lower severity
- Flag any mismatch between the alert's `severity` label and the
  urgency conveyed in the Impact section
- Based on where the metrics are set in source code, what functionality
  is actually affected when this condition is true?
- Does the runbook capture the real user-facing impact, or is it
  vague/generic?

#### Cross-check: Meaning vs Impact

- The Meaning and Impact sections must not contradict each other. For
  example, if Meaning says "no ready pods exist", Impact must not say
  "reduced redundancy" — it should say "complete unavailability".
  Flag contradictions.

#### `## Diagnosis`

- `kubectl`/`oc` commands must target the correct resources and labels
  for the component the alert monitors
- Verify that diagnosis commands use the correct Kubernetes resource
  type for the component. Known mappings:
  - `virt-handler` → DaemonSet (NOT Deployment)
  - `virt-api` → Deployment
  - `virt-controller` → Deployment
  - `virt-operator` → Deployment
  Flag any command that uses the wrong resource type (e.g.,
  `kubectl get deployment virt-handler` is wrong because virt-handler
  is a DaemonSet).
- Use label-based queries, not hardcoded pod names:
  - Good: `kubectl -n $NAMESPACE get pod -l kubevirt.io=virt-api`
  - Bad: `kubectl -n $NAMESPACE get pod virt-api-xyz`
- When the alert condition involves pods in CrashLoopBackOff or
  failure states, diagnosis MUST include a log-checking step:
  `kubectl -n $NAMESPACE logs -l kubevirt.io=<component>`
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

## Writing Style

- Documentation must be written in present tense. Flag future tense
  ("will be removed", "will happen") unless it refers to a genuinely
  future event. Examples:
  - Correct: "Ephemeral volumes are lost on restart."
  - Incorrect: "Ephemeral volumes will be lost on restart."
- When a PR reformats or rewords existing runbook text, verify that
  the technical meaning is preserved. Changing "if" to "when" or "are"
  to "will" can alter semantics. Flag any meaning change disguised as
  a style edit.

## Formatting Rules

- One runbook per file: `<AlertName>.md`
- H1 heading must match filename: `VirtAPIDown.md` → `# VirtAPIDown`
- Required H2 sections in order: `## Meaning`, `## Impact`,
  `## Diagnosis`, `## Mitigation`
- Code fences: ```` ```bash ```` (no space before language tag)
- Line length: 80 characters for prose (code blocks exempt).
  URLs inside markdown links are exempt from the 80-character limit.
- Prefix example commands with `$ `
- `<!--DS: ... -->` and `<!--USstart-->...<!--USend-->` are intentional
  downstream/upstream content markers — do not flag them
- Deprecated runbooks are moved to `docs/deprecated_runbooks/`, not
  deleted.
