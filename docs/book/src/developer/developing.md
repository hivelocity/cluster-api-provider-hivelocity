
# Developing Cluster API Provider Hivelocity

Developing our provider is quite easy. First, you need to install some base requirements. Second, you need to follow the quickstart documents to set up everything related to Hivelocity. Third, you need to configure your Tilt set-up. After having done those three steps, you can start developing with the local Kind cluster and the Tilt UI to create one of the different workload clusters that are already pre-configured.


## Why Tilt

Provider Integration development requires a lot of iteration, and the “build, tag, push, update deployment” workflow can be very tedious.

Tilt makes this process much simpler by watching for updates and automatically building and deploying them.

You just need to update the Go code, and if it compiles, the caphv-controller-manager will get updated.

## Updating the API of the CRDs

If you update the API in some way, you need to run `make generate` in order to generate everything related to kubebuilder and the CRDs.

## Deleting the cluster

To tear down the workload cluster press the "Delete Workload Cluster" button. After a few minutes the resources should be deleted.

To tear down the kind cluster, use:

```shell
$ make delete-cluster
```

To delete the registry, use: `make delete-registry` or `make delete-cluster-registry`.

## make help

If you have any trouble finding the right command, then you can use `make help` to get a list of all available make targets.

## Submitting PRs and testing

<aside class="note info">

Pull requests and issues are highly encouraged! For more information, please have a look in the [Contribution Guidelines](../reference/CONTRIBUTING.md)

</aside>

There are two important commands that you should make use of before creating the PR.

With `make verify` you can run all linting checks and others. Make sure that all of these checks pass - otherwise the PR cannot be merged. Note that you need to commit all changes for the last checks to pass.

With `make test` all unit tests are triggered. If they fail out of nowhere, then please re-run them. They are not 100% stable and sometimes there are tests failing due to something related to Kubernetes' `envtest`.

With `make generate` new CRDs are generated, this is necessary if you change the api.

With [make watch](../topics/make-watch.md) you can monitor the progress.
