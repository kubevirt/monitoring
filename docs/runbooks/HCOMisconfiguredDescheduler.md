# HCOMisconfiguredDescheduler

## Meaning

A descheduler is a Kubernetes application that causes the control plane to
re-arrange the workloads in a more optimized way.

The descheduler uses the Kubernetes eviction API to evict pods. However,
KubeVirt VMs, which run as `virt-launcher` pods, require special handling
so they can be properly migrated by KubeVirt through live migration.

In earlier releases, a specific `profileCustomizations` setting
(`devEnableEvictionsInBackground`) was required to configure the descheduler
correctly for KubeVirt. More recent releases include dedicated descheduler
profiles that already handle this requirement.

## Impact

Using the descheduler operator for KubeVirt VMs without the correct
configuration may result in unstable or unpredictable behavior,
which can negatively affect cluster stability.

## Diagnosis

1. Check the `KubeDescheduler` custom resource (CR):

   ```bash
   kubectl get -n openshift-kube-descheduler-operator KubeDescheduler cluster -o yaml
   ```

2. Verify that the CR contains one of the following profile configurations:

   ```yaml
   spec:
     profiles:
     - KubeVirtRelieveAndMigrate
   ```

   or

   ```yaml
   spec:
     profiles:
     - DevKubeVirtRelieveAndMigrate
   ```

   or

   ```yaml
   spec:
     profiles:
     - LongLifecycle
     profileCustomizations:
       devEnableEvictionsInBackground: true
   ```

If none of these are present, the `Kube Descheduler Operator` is not correctly
configured for KubeVirt.

## Mitigation

Configure the descheduler with an appropriate profile for your cluster.
The order of preference is:

1. `KubeVirtRelieveAndMigrate`  
2. `DevKubeVirtRelieveAndMigrate`  
3. `LongLifecycle` with the `devEnableEvictionsInBackground` customization  

To check which profiles are supported in your descheduler release,
use the following command:

```bash
oc get crd kubedeschedulers.operator.openshift.io -o json | jq '.spec.versions[] | select(.name=="v1").schema.openAPIV3Schema.properties.spec.properties.profiles.items.enum[]'
```

Depending on the available profiles, apply one of the following configurations
to the `KubeDescheduler` CR (kind: `kubedeschedulers`, namespace:
`openshift-kube-descheduler-operator`, name: `cluster`):

* **Configure the `KubeVirtRelieveAndMigrate` profile**  
   ```yaml
   spec:
     profiles:
     - KubeVirtRelieveAndMigrate
   ```

* **Configure the `DevKubeVirtRelieveAndMigrate` profile**  
   ```yaml
   spec:
     profiles:
     - DevKubeVirtRelieveAndMigrate
   ```

* **Configure the `LongLifecycle` profile with the
`devEnableEvictionsInBackground` profile customization**  
   ```yaml
   spec:
     profiles:
     - LongLifecycle
     profileCustomizations:
       devEnableEvictionsInBackground: true
   ```

* If none of the listed steps are possible,
**remove the `Kube Descheduler Operator`**.

<!--DS: If you cannot resolve the issue, log in to the
link:https://access.redhat.com[Customer Portal] and open a support case,
attaching the artifacts gathered during the diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://okd.io/docs/community/help)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
