# VirtOperatorDown

## Meaning

This alert fires when no `virt-operator` pod in the `Running` state has been
detected for 10 minutes.

The `virt-operator` is the first Operator to start in a cluster. Its primary
responsibilities include the following:

- Installing, live-updating, and live-upgrading a cluster
- Monitoring the life cycle of top-level controllers, such as `virt-controller`,
`virt-handler`, `virt-launcher`, and managing their reconciliation
- Certain cluster-wide tasks, such as certificate rotation and infrastructure
management

The `virt-operator` deployment has a default replica of 2 pods.

## Impact

This alert indicates a failure at the level of the cluster. Critical
cluster-wide management functionalities, such as certification rotation,
upgrade, and reconciliation of controllers, might not be available.

The `virt-operator` is not directly responsible for virtual machines (VMs) in
the cluster. Therefore, its temporary unavailability does not significantly
affect VM workloads.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-operator` deployment:

   ```bash
   $ kubectl -n $NAMESPACE get deploy virt-operator -o yaml
   ```

3. Obtain the details of the `virt-operator` deployment:

   ```bash
   $ kubectl -n $NAMESPACE describe deploy virt-operator
   ```

4. Check the status of the `virt-operator` pods:

   ```bash
   $ kubectl get pods -n $NAMESPACE -l=kubevirt.io=virt-operator
   ```

5. Check for node issues, such as a `NotReady` state:

   ```bash
   $ kubectl get nodes
   ```

## Mitigation

Based on the information obtained during the diagnosis procedure, try to
identify the root cause and resolve the issue.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
