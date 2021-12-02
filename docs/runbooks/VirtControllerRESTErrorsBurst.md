# VirtControllerRESTErrorsBurst

## Meaning

More than 80% of the rest calls failed in virt-controller for the last 5 minutes

## Impact

virt-controller has potentially fully lost connectivity to the apiserver. Running workloads will not be impacted but status updates will not be propagated and actions like migrations can not take place.

## Diagnosis

There exist two common types of errors:
- The apiserver is overloaded and we run into timeouts.
  Issues like this can be identified by checking the apiserver metrics and looking at its response times, overall calls, …
- The virt-controller pod can’t reach the apiserver. Common issues are: dns issues on the node, networking connectivity issues

## Mitigation

Check virt-controller logs to identify if it can’t connect to the apiserver at all. If so, delete the pod to force a restart. In this case the issue is normally related to DNS or CNI issues outside of the scope of kubevirt.