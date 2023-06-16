# E2E tests

End-to-end tests run in CI, so that changes get automatically tested.

## Running local e2e test

If you are interested in running the e2e tests locally, then you can use the following commands:
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

check [make watch](../topics/make-watch.md)