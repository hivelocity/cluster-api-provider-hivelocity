# Development with Tilt

```
make tilt-up
```

This will:

* create a local container registry with `ctlptl`
* create a management-cluster in [Kind](https://kind.sigs.k8s.io/)
* start [Tilt](https://tilt.dev/)

Set the device Tag of a device to `caphv-device-type=hvCustom`. Web-GUI API: https://developers.hivelocity.net/reference/put_device_tag_id_resource

Upload a SSH public key with the name "ssh-key-hivelocity-pub" (see
templates/cluster-templates/bases/hivelocity-hivelocityCluster.yaml) with
[API Add SSH Key](https://developers.hivelocity.net/reference/post_ssh_key_resource).

When Tilt is started (everything is green), you can create a Hivelocity cluster with the corresponding button (at the top) in Tilt.

You can use `make watch` to see the output of relevant `kubectl` commands and the tail of the controller.

More: [Logging](logging.md)
