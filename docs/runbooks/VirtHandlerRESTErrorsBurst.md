# VirtHandlerRESTErrorsBurst

## Meaning

More than 80% of REST calls failed in `virt-handler` in the last 5 minutes.

The `virt-handler` has likely fully lost connectivity to the apiserver. Running workloads on the affected node will not be impacted, but status updates will not be propagated and node-related actions, such as migrations, will fail.

## Diagnosis

This error is most frequently caused by one of the following problems:

- The apiserver is overloaded, which causes timeouts. To verify if this is the case, check the metrics of the apiserver, and view its response times and overall calls.

- The `virt-handler` pod cannot reach the apiserver. This is commonly caused by DNS issues on the node and networking connectivity issues.

Check `virt-handler` logs to identify if it can connect to the apiserver. (TBA how?)

## Mitigation
If the `virt-handler` cannot connect to the apiserver, delete the pod to force a restart. (TBA how?)

