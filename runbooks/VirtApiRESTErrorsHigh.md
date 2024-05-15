# VirtApiRESTErrorsHigh

## Meaning

More than 5% of REST calls have failed in the `virt-api` pods in the last 60
minutes.

## Impact

A high rate of failed REST calls to `virt-api` might lead to slow response and
execution of API calls.

However, currently running virtual machine workloads are not likely to be
affected.

## Diagnosis

1. Set the `NAMESPACE` environment variable as follows:

   ```bash
   $ export NAMESPACE="$(kubectl get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-api` pods:

   ```bash
   $ kubectl -n $NAMESPACE get pods -l kubevirt.io=virt-api
   ```

3. Check the `virt-api` logs:

   ```bash
   $ kubectl logs -n  $NAMESPACE <virt-api>
   ```

4. Obtain the details of the `virt-api` pods:

   ```bash
   $ kubectl describe -n $NAMESPACE <virt-api>
   ```

5. Check if any problems occurred with the nodes. For example, they might be in
a `NotReady` state:

   ```bash
   $ kubectl get nodes
   ```

6. Check the status of the `virt-api` deployment:

   ```bash
   $ kubectl -n $NAMESPACE get deploy virt-api -o yaml
   ```

7. Obtain the details of the `virt-api` deployment:

   ```bash
   $ kubectl -n $NAMESPACE describe deploy virt-api
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
