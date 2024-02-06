# Provisioning Machines

How can the CAPHV controller know which machines are free to use for a cluster?

You don't want existing machines in your account to get re-provisioned to cluster nodes :-)

When starting a cluster you define a worker-machine-type and a control-plane-machine-type. Both values can be equal, when you don't want to differentiate between both types.

Before you create your cluster you need to tag the devices accordingly.

The devices must have `caphv-use=allow` tag so that the controller can use them.

For example you set tags "caphvlabel:deviceType=hvControlPlane" on all machines which should become control planes, and "caphvlabel:deviceType=hvWorker" on all machines which should become worker nodes.

You can use the web-GUI of Hivelocity for this.

Then the CAPHV controller is able to select machines, and then provision them to become Kubernetes nodes.

The CAPHV controller uses [Cluster API bootstrap provider kubeadm](https://cluster-api.sigs.k8s.io/tasks/bootstrap/kubeadm-bootstrap.html) to provision the machines.

:warning: If you create a cluster with `make tilt-up` or other Makefile targets, then all machines having a
corresponding `caphvlabel:deviceType=` will get all their tags cleared. This means the machine is free to use,
and it is likely to become automatically provisioned. This means all data on this machine gets lost.

TODO: https://github.com/hivelocity/cluster-api-provider-hivelocity/issues/73
