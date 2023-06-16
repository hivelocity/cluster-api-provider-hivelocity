# Hivelocity API

Hivelocity provides a rich API which you can access and test with your web browser:

[developers.hivelocity.net/reference](https://developers.hivelocity.net/reference/)

## APIs

Here are some API which gets used by CAPHV:

### [get_bare_metal_device_resource](https://developers.hivelocity.net/reference/get_bare_metal_device_resource)

Get all devices. CAPHV uses this API to search for devices which are free to get provisioned.

### [get_bare_metal_device_id_resource](https://developers.hivelocity.net/reference/get_bare_metal_device_id_resource)

Get a single device. CAPHV uses this API to read the tags of a single device.

### [put_bare_metal_device_id_resource](https://developers.hivelocity.net/reference/put_bare_metal_device_id_resource)

Update/reload instant device. CAPHV uses this API to provision a device.

Provisioning gets done via [cloud-init](https://cloudinit.readthedocs.io/en/latest/)

### [post_ssh_key_resource](https://developers.hivelocity.net/reference/post_ssh_key_resource)

Add public ssh key

### [put_device_tag_id_resource](https://developers.hivelocity.net/reference/put_device_tag_id_resource)

Update device tags. CAPHV uses this API to ensure that device is part of exactly one cluster.

## Client Go

CAPHV uses [hivelocity-client-go](https://github.com/hivelocity/hivelocity-client-go) to access the API from the programming language Golang.

This client gets automatically created from the `swagger.json` file provided by Hivelocity.

