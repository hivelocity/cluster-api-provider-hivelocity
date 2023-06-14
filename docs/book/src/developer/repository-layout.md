# Repository Layout

```
Here is your text with corrections to typos and grammar mistakes:

# Repository Layout

```
├── api                   # This folder is used to store types and their related resources (Go code) present in CAPHV.
|                         # The API folder has subfolders for each supported API version.
|
├── _artifacts            # This directory is created during e2e tests. It contains yaml files and logs.
|
├── bin                   # Binaries, mostly for envtests.
|
├── config                # This is a Kubernetes manifest folder containing application resource configurations as
|                         # kustomize YAML definitions. These are generated from other folders in the repo using
|                         # `make generate-manifests`.
|                         # More details are in the upstream docs:
|                         # https://cluster-api.sigs.k8s.io/developer/repository-layout.html#manifest-generation
|
├── controllers           # This folder contains reconciler types which provide access to CAPHV controllers.
|                         # These types can be used by users to run any of the Cluster API controllers in an external program.
|
├── docs                  # Source for: https://hivelocity.github.io/cluster-api-provider-hivelocity/
|
├── hack                  # This folder has scripts used for building, testing, and the developer workflow.
|
├── main.go               # This contains the main function for CAPHV. The code gets compiled to the binary `manager`.
|
├── Makefile              # Configuration for building CAPHV. Try `make help` to get an overview.
|
├── pkg                   # Go code used to implement the controller.
├── templates
│   ├── cilium            # CNI for workload clusters (you can choose a different CNI, too).
│   └── cluster-templates # YAML files for the management cluster.
|
├── test                  # Config and code for e2e tests.
|
├── Tiltfile              # Configuration for https://tilt.dev/.
├── tilt-provider.yaml
├── tilt-settings.yaml
├── .envrc.example        # We use direnv to set environment variables. This is optional. See https://direnv.net/.
├── .github               # GitHub config: PR templates, CI, ...
```
