# KubeVirtDeprecatedAPIRequested
<!-- Edited by fmatouschek, May 2023-->

## Meaning

This alert fires when a deprecated KubeVirt API is requested.

## Impact

Usage of deprecated APIs is not recommended because they will be removed in a future release.

## Diagnosis

Check the `description` and `summary` alert annotations for more details on which API is being accessed, for example:
```
              description: "Detected requests to the deprecated virtualmachines.kubevirt.io/v1alpha3 API."
              summary: "2 requests were detected in the last 10 minutes."
```

## Mitigation

Make sure to only use a supported version when making requests to the API.

Some requests to deprecated APIs are made by KubeVirt components themselves (e.g VirtualMachineInstancePresets).
These alerts cannot be mitigated because the requests are still necessary to serve the deprecated API.
They are harmless and will be resolved when the deprecated API is removed in a future release of KubeVirt.

Alerts will resolve after 10 minutes if the deprecated API is not used again.

If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
