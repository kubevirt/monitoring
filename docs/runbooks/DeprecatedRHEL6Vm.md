# DeprecatedRHEL6Vm

## Meaning

This alert fires when a running virtual machine is based on RHEL6 template, which will not be supported in the next release.

## Impact
RHEL6 virtual machines will not be supported in the next release in case of issues.

## Diagnosis

1. Check which running virtual machines has the rhel6 template label.

   ```bash
   $ kubectl get -A vm -ojson | jq  '.items | map(select((.metadata.labels."vm.kubevirt.io/template" | test("rhel6-.*")?) and .status.printableStatus == "Running")) | {kind: "List", items: .}'  | kubectl get -A -f -
   ```

## Mitigation

Upgrade all the RHEL6 virtual machines to a newer version or stop them.\
Since the warning looks at the value of the `vm.kubevirt.io/template` label and is triggered if it starts with `rhel6`,
in case of an update, you should also manually change its value to a correct value other than `rhel6-*`.

If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)

