# KubevirtHyperconvergedClusterOperatorInstallationNotCompletedAlert

## Meaning
The `KubevirtHyperconvergedClusterOperatorInstallationNotCompletedAlert` alert means that the HypeConverged cluster operator (HCO) is running, but the `HyperConverged` custom resource does not exist for more than an hour.

There are two reasons for that to happen: First, after installation of HCO, it is required to create the `HyperConverged` custom resource in order to complete the installation. The purpose of this alert is to guide the user to complete the installation by creating the `HyperConverged` custom resource. 

The other reason for this alert is during uninstallation of HCO. In this case, the user is instructed to remove the `HyperConverged` custom resource before uninstalling the operator. If the `HyperConverged` custom resource was indeed removed, but the HCO is still installed and running for more than an hour, the alert will be triggered. In this case, the alert can usually be safely ignored.

## Mitigation
In case of installation, just create the `HyperConverged` custom resource to complete the installation; for example, to create the `HyperConverged` custom resource with default values, do:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: hco-operatorgroup
  namespace: kubevirt-hyperconverged
spec: {}
EOF
```

In case of uninstallation, if HCO is still installed and running for more than an hour after removing the `HyperConverged` custom resource, complete the process by uninstalling HCO. If for some reason, HCO is still running (uninstallation is stuck), this issue should be solved in order to cancel the alert.
