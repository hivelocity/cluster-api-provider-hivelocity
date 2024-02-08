# Getting Started

## Terminology

Before we begin we need to answer some questions.

## What is Hivelocity?

Hivelocity provides Dedicated Servers, Colocation and Cloud Hosting services to customers from over 130 countries since 2002. Hivelocity operates over 70,000 sq ft of data center space offering services in Tampa FL, Miami FL, Atlanta GA, New York NY, and Los Angeles CA. Each of Hivelocity's data centers are HIPAA, PCI, ISAE-3402, SSAE 16 SOC1 & SOC2 certified.

## What is Kubernetes?

Kubernetes is an open-source container orchestration system for automating software deployment, scaling, and management. Originally designed by Google, the project is now maintained by the Cloud Native Computing Foundation.

Kubernetes defines a set of building blocks that collectively provide mechanisms that deploy, maintain, and scale applications based on CPU, memory or custom metrics.

Source: [Wikpedia](https://en.wikipedia.org/wiki/Kubernetes)

## What is Cluster API?

Cluster API is a Kubernetes sub-project focused on providing declarative APIs and tooling to simplify provisioning, upgrading, and operating multiple Kubernetes clusters.

The Cluster API project uses Kubernetes-style APIs and patterns to automate cluster lifecycle management for platform operators. The supporting infrastructure, like virtual machines, networks, load balancers, and VPCs, as well as the Kubernetes cluster configuration are all defined in the same way that application developers operate deploying and managing their workloads. This enables consistent and repeatable cluster deployments across a wide variety of infrastructure environments.

Source: [cluster-api.sigs.k8s.io](https://cluster-api.sigs.k8s.io/)

Cluster API uses Kubernetes Controllers: In a **management-cluster** runs a controller which reconciles the state of **workload-clusters** until the state reaches the desired state.

The desire state gets specified in yaml manifests.

![cluster-api: management-cluster and workload-clusters](./cluster-api.png)

After the workload-cluster was created successfully, you can **move** the management-cluster into the workload-cluster. This gets done with `clusterclt move`. See [Cluster API Docs "pivot"](https://cluster-api.sigs.k8s.io/clusterctl/commands/move.html#pivot)


## What is Cluster API Provider Hivelocity?

Cluster API Hivelocity adds the infrastructure provider Hivelocity to the list of supported providers. Other providers supported by Cluster API are: AWS, Azure, Google Cloud Platform, OpenStack ... (See [complete list](https://cluster-api.sigs.k8s.io/reference/providers.html#infrastructure))

## Setup

At this moment we only support cluster management with Tilt. So follow below instructions to create management and workload cluster.

### Create a management cluster

Please run below command and this will use `tilt-provider.yaml` to create the provider and
`.envrc` to get all the environment variables.
```shell
# Please run the command from root of this repository
make tilt-up
```

### Create workload cluster

There is a button in top right corner of the Tilt console to create the workload cluster.
![Screenshot of Tilt](./create_workload.jpg)

### Tear down resources

There is a button in top right corner of the Tilt console to create the workload cluster.
![Screenshot of Tilt](./delete_workload.jpg)

Once done delete management cluster by -
```shell
make delete-mgt-cluster
```

## Current Limitations

### Limitation: Missing Loadbalancers

Up to now Loadbalancers are not supported yet. But we are working on it.

See [issue #55](https://github.com/hivelocity/cluster-api-provider-hivelocity/issues/55)

### Limitation: Missing VPC Networking

Up to now the machines have public IPs.

For security nodes should not be accessible from the public internet.

See [issue #78](https://github.com/hivelocity/cluster-api-provider-hivelocity/issues/78)

## Current Known Issues

### Known Issue: Broken Machine State

Sometimes machines get stuck. The hang in state "Reloading" forever. Then the support of Hivelocity
need to reset the machine.

Related issue at Github [#59](https://github.com/hivelocity/cluster-api-provider-hivelocity/issues/59).

## Current State: alpha

Up to now CAPHV is not in the [official list of infrastructure providers](https://cluster-api.sigs.k8s.io/reference/providers.html#infrastructure).

But we are working on it.

Please have a look at the [Developer Guide](../developer/index.md), if you want to setup a cluster.

## Navigation in the docs

On the left and right side of the documentation you see angle brackets. You can use them to switch to the next/previous page.
