# Development with Tilt

```
make tilt-up
```

This will:

* create a local container registry with `ctlptl`
* create a management-cluster in [Kind](https://kind.sigs.k8s.io/)
* start [Tilt](https://tilt.dev/)

Set the device Tag of a device to `caphv-device-type=hvCustom`. Web-GUI API: https://developers.hivelocity.net/reference/put_device_tag_id_resource


When Tilt is started (everything is green), you can create a Hivelocity cluster with the corresponding button (at the top) in Tilt.

More: [Logging](logging.md)
