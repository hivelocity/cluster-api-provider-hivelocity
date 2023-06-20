# make watch

`make watch` is an ad-hoc monitoring tool that checks whether the management cluster and the workload cluster 
are operational.

It's main use-case is to monitor the bootstrapping of new clusters.

It is not a replacement for a professional monitoring solution.

If there is an error, then `make watch` will show it (❌). 

Please create a [Github Issue](https://github.com/hivelocity/cluster-api-provider-hivelocity/issues) if `make watch` does not detect an error.

When the cluster starts some errors are common.

For example it will take up to 22 minutes until the ssh port of the first control plane is reachable.

Here a screenshot, six minutes after the Cluster and MachineDeployment manifests got applied to the management-cluster.

![make watch](./make-watch.jpg)


