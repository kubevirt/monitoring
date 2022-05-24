## Monitoring Guidelines

- [Monitoring Guidelines](#monitoring-guidelines)
  - [KubeVirt Metrics](#kubevirt-metrics)
    - [Naming a New KubeVirt Metrics](#naming-a-new-kubevirt-metrics)
  - [KubeVirt Recording Rules](#kubevirt-recording-rules)
    - [Naming a New KubeVirt Recording Rule](#naming-a-new-kubevirt-recording-rule)
  - [KubeVirt Alerts Rules](#kubevirt-alerts-rules)
 
### KubeVirt Metrics
#### Naming a New KubeVirt Metrics

The KubeVirt metrics should align with the Kubernetes metrics names.

The KubeVirt Users should have the same experience when searching for a node, container, pod and virtual machine metrics.

**Naming requirements:**
1. Check if a similar Kubernetes metric, for node, container or pod, exists and try to align to it.
2. KubeVirt metrics prefixes:
   1.  Running VM metrics should have a `kubevirt_vmi_` prefix
   2.  HCO operator metrics should have a `kubevirt_hco_` prefix
   3.  Network operator metrics should have a `kubevirt_network_` prefix
   4.  Storage operator metrics should have a `kubevirt_cdi_` prefix
   5.  SSP operator metrics should have a `kubevirt_ssp_` prefix
   6.  HPP Operator metrics should have a `kubevirt_hpp_` prefix


    For Example, see the following Kubernetes network metrics:
    - **node**_network_receive_packets_total
    - **node**_network_transmit_packets_total
    - **container**_network_receive_packets_total
    - **container**_network_transmit_packets_total

    The KubeVirt metrics for vmi should be:
    - **kubevirt_vmi**_network_receive_packets_total
    - **Kubevirt_vmi**_network_transmit_packets_total


3. Metric `Help` message MUST be verbose, since it is being propagated to the [metrics.md](https://github.com/kubevirt/kubevirt/blob/main/docs/metrics.md) file, when running `make-generate`.

### KubeVirt Recording Rules

#### Naming a New KubeVirt Recording Rule

Use [recording rules](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/#recording-rules) when doing calculations or when using the same query for other alerts or dashboards, instead of repeating the same query in many places.

The Prometheus recording rules appear in Prometheus as metrics.

In order to easily identify the KubeVirt recording rules, they should follow the same naming conventions as the metrics.

### KubeVirt Alerts Rules

When creating a KubeVirt alert rule, please follow the [OpenShift Alerting Consistency Guide](https://github.com/openshift/enhancements/blob/master/enhancements/monitoring/alerting-consistency.md#alerting-consistency).

In addition to the [OpeShift Style Guide](https://github.com/openshift/enhancements/blob/master/enhancements/monitoring/alerting-consistency.md#style-guide) the KubeVirt alerts MUST include:
1. `kubernetes_operator_part_of` label indicating the operator name. Value should be set to `kubevirt`.
2. `kubernetes_operator_component` label indicating the value of the sub operator name.

Optional labels:
1. `priority` label indicating the alert's level of importance and the order in which it should be fixed.
  * Valid priorities are: `high`, `medium`, or `low`.
    The higher the priority, the sooner the alert should be resolved.
  * If the alert doesn't include a `priority` label, we can assume it is a `medium` priority alert.
2. `infra_alert` label indicating alerts that are related to the infrastructure of the operator. Boolean.

**Note:**
KubeVirt alert runbooks are saved in [kubevirt/monitoring repository](https://github.com/kubevirt/monitoring/tree/main/docs/runbooks).