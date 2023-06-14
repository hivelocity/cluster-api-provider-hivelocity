# Development with Tilt

We use [Tilt](https://tilt.dev/) to start a development cluster.

First you need to label machines, so that the CAPHV controller can select and provision them.

See [[Provisioning Machines](../topics/provisioning-machines.md) for more about this.

You can use `go run ./test/claim-devices-or-fail` to claim all devices
with a particular label. But be careful, since all machines with this label will get all existing labels removed. Don't accidentally "free" running machines this way.

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

#############

# Developing Cluster API Provider Hivelocity

Developing our provider is quite easy. First, you need to install some base requirements. Second, you need to follow the quickstart documents to set up everything related to Hivelocity. Third, you need to configure your tilt set-up. After having done those three steps, you can start developing with the local Kind cluster and the Tilt UI to create one of the different workload clusters that are already pre-configured.

## Install Base requirements

In order to develop with Tilt, there are a few requirements. You can use the following command to check whether the versions of the tools are up-to-date and to install ones that are missing (for both mac & linux): `make install-dev-prerequisites`

This ensures the following:
- clusterctl
- ctlptl (required)
- go (required)
- helm (required)
- kind (required)
- kubectl (required)
- tilt (required)


## Setting Tilt up

You need to create a ```tilt-settings.yaml``` file and specify the values you need. Here is an example:

```yaml
kustomize_substitutions:
  HIVELOCITY_API_KEY: dummy-key
  HIVELOCITY_SSH_KEY: test
  HIVELOCITY_REGION: LAX2
  CONTROL_PLANE_MACHINE_COUNT: "3"
  WORKER_MACHINE_COUNT: "3"
  KUBERNETES_VERSION: v1.25.2
  HIVELOCITY_IMAGE_NAME: 1.25.2-ubuntu-20.04-containerd
  HIVELOCITY_CONTROL_PLANE_MACHINE_TYPE: todo
  HIVELOCITY_WORKER_MACHINE_TYPE: todo
```

## Developing with Tilt

<p align="center">
<img alt="tilt" src="../pics/tilt.png" width=800px/>
</p>

Provider Integration development requires a lot of iteration, and the “build, tag, push, update deployment” workflow can be very tedious. Tilt makes this process much simpler by watching for updates and automatically building and deploying them. To build a kind cluster and to start Tilt, run:

```shell
make tilt-up
```
> To access the Tilt UI please go to: `http://localhost:10350`


Once your kind management cluster is up and running, you can deploy a workload cluster. This could be done through the Tilt UI, by pressing one of the buttons in the top right corner, e.g. "Create Hivelocity Cluster". This triggers the `make create-workload-cluster`, which uses the environment variables (we defined in the tilt-settings.yaml) and the cluster-template. Additionally, it installs cilium as CNI.

If you update the API in some way, you need to run `make generate` in order to generate everything related to kubebuilder and the CRDs.

To tear down the workload cluster press the "Delete Workload Cluster" button. After a few minutes the resources should be deleted.

To tear down the kind cluster, use:

```shell
$ make delete-cluster
```

To delete the registry, use: `make delete-registry` or `make delete-cluster-registry`.

If you have any trouble finding the right command, then you can use `make help` to get a list of all available make targets.

## Submitting PRs and testing

Pull requests and issues are highly encouraged! For more information, please have a look in the [Contribution Guidelines](../../CONTRIBUTING.md)

There are two important commands that you should make use of before creating the PR.

With `make verify` you can run all linting checks and others. Make sure that all of these checks pass - otherwise the PR cannot be merged. Note that you need to commit all changes for the last checks to pass.

With `make test` all unit tests are triggered. If they fail out of nowhere, then please re-run them. They are not 100% stable and sometimes there are tests failing due to something related to Kubernetes' `envtest`.

With `make generate` new CRDs are generated, this is necessary if you change the api.

### Running local e2e test

If you are interested in running the E2E tests locally, then you can use the following commands:
```
export HIVELOCITY_API_KEY=<api-key>
export CAPHV_LATEST_VERSION=<latest-version>
make test-e2e
```

## Updating controller during e2e test

You discovered an error in the code, or you want to add some logging statement while the e2e is running?

Usualy development gets done with Tilt, but you can update the image during an e2e test with these steps:

Re-create the image:
```
make e2e-image
```

```
config=$(yq '.clusters[0].name' .mgt-cluster-kubeconfig.yaml)
kind_cluster_name="${config#kind-}"

# refresh the image
kind load docker-image ghcr.io/hivelocity/caphv-staging:e2e --name=$kind_cluster_name

# restart the pod
kubectl rollout restart deployment -n capi-hivelocity-system caphv-controller-manager
```

... check `make watch`
