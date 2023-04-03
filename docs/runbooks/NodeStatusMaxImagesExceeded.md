# NodeStatusMaxImagesExceeded

## Meaning

This alert fires when a node exceeds the configured maximum number of images
reported in the node status (default: 50). Once this limit is reached, the
scheduler might not be able to accurately schedule new pods because the reported
image count is wrong. [BZ#1984442](https://bugzilla.redhat.com/1984442)

## Impact

Pod scheduling might be imbalanced across nodes.

## Diagnosis

1. Obtain the image count for your nodes:

   ```bash
   $ kubectl get nodes -o json |  jq -r '[.items[] | { name: .metadata.name, images: .status.images | length }]'
   [
      {
         "name": "crc-chf7b-master-0",
         "images": 50
      },
      {
         "name": "crc-chf7b-master-1",
         "images": 20
      },
      {
         "name": "crc-chf7b-master-2",
         "images": 50
      },
   ]
   ```

2. If the image count is equal to the default (50) or to a configured value,
you must resolve the issue.

## Mitigation

You can either increase the `nodeStatusMaxImages` value of the `KubeletConfig`
object or set its value to `-1` to disable the maximum image limit.

<!--DS: If you cannot resolve the issue, log in to the link:https://access.redhat.com[Customer Portal] and open a support case, attaching the artifacts gathered during the Diagnosis procedure.-->
<!--USstart-->
If you cannot resolve the issue, see the following resources:

- [OKD Help](https://www.okd.io/help/)
- [#virtualization Slack channel](https://kubernetes.slack.com/channels/virtualization)
<!--USend-->
