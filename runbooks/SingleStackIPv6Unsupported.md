# SingleStackIPv6Unsupported

## Meaning

This alert fires when user tries to install KubeVirt Hyperconverged on a single
stack IPv6 cluster.

KubeVirt Hyperconverged is not yet supported on an OpenShift cluster configured
with single stack IPv6. It's progress is being tracked on [this issue](https://issues.redhat.com/browse/CNV-28924).

## Impact

KubeVirt Hyperconverged Operator can't be installed on a single stack IPv6
cluster, and hence creation virtual machines on top of such a cluster is not
possible.

## Diagnosis

1. Obtain the network details for the cluster using:
   ```shell
    $ kubectl describe network cluster
    ```

2. It contains only an IPv6 CIDR under `Cluster Network`.

## Mitigation

It is recommended to use single stack IPv4 or a dual stack IPv4/IPv6 networking
to use KubeVirt Hyperconverged. Refer the [documentation](https://docs.openshift.com/container-platform/latest/networking/ovn_kubernetes_network_provider/converting-to-dual-stack.html).
