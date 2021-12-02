# VirtHandlerRESTErrorsHigh

## Meaning

More than 5% of the rest calls failed in virt-handler for the last hour

## Impact

Virt-handler has partially lost the connection to the apiserver. Node-related actions (where virt-handler is running on) like starting and migrating workloads will be delayed. Running workloads will not be affected, but reporting their current status may be delayed.

## Diagnosis

There exist two common types of errors:
- The apiserver is overloaded and we run into timeouts. Issues like this can be identified by checking the apiserver metrics and looking at its response times, overall calls, …
- The virt-handler pod can’t reach the apiserver. Common issues are: dns issues on the node, networking connectivity issues

## Mitigation
Check virt-handler logs to identify if it can’t connect to the apiserver at all. If so, delete the pod to force a restart. In this case the issue is normally related to DNS or CNI issues outside of the scope of kubevirt.
