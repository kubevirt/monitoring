# VirtControllerRESTErrorsHigh

## Meaning

More than 5% of the rest calls failed in virt-controller for the last hour

## Impact

Virt-controller has partially lost the connection to the apiserver.
Cluster-level related actions like starting, migrating and scheduling VMs, will be delayed.
Running workloads will not be affected, but reporting their current status may be delayed.

## Diagnosis

There exist two common types of errors:
1. The apiserver is overloaded and we run into timeouts. Issues like this can be identified by checking the apiserver metrics and looking at its response times, overall calls, …
2. The virt-controller pod can’t reach the apiserver. Common issues are: dns issues on the node, networking connectivity issues

## Mitigation
Check virt-controller logs to identify if it can’t connect to the apiserver at all. If so, delete the pod to force a restart.
In this case the issue is normally related to DNS or CNI issues outside of the scope of kubevirt.
