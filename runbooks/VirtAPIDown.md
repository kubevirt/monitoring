# VirtAPIDown

## Meaning

This alert fires when all KubeVirt API servers are down.

## Impact

Without KubeVirt API servers, no API call around KubeVirt entities can be made anymore.

## Diagnosis

Study the virt-api deployment and ensure that the minimum number of pods is running.

## Mitigation

There can be several reasons for virt-api pods to be down, identify the root cause and fix it.

For example:
- …
- …
