# Provisioning Machines

How can the CAPHV controller know which machines are free to use for a cluster?

You don't want existing machines in your account to get re-provisioned to cluster node :-)

When starting a cluster you define a worker-machine-type and a control-plane-machine-type. Both values can be equal, when you don't want to differentiate between both types.

Before you create your cluster you need to label the devices accordingly.

For example you set label "caphv-device-type=hvControlPlane" on all machines which should become control planes, and "caphv-device-type=hvWorker" on all machines which should become worker nodes.

You can use the web-GUI of Hivelocity for this.

Then the CAPHV controller is able to select machines, and then provision them to become Kubernetes nodes.

The CAPHV controller uses [Cluster API bootstrap provider kubeadm](https://cluster-api.sigs.k8s.io/tasks/bootstrap/kubeadm-bootstrap.html) to provision the machines.
