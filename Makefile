# Copyright 2023 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.DEFAULT_GOAL:=help

# Go.
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Directories.
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
EXP_DIR := exp
BIN_DIR := bin
TEST_DIR := test
TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/$(BIN_DIR)
export PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)
# Default path for Kubeconfig File.

# Files
CAPHV_WORKER_CLUSTER_KUBECONFIG ?= "/tmp/workload-kubeconfig"

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.25.0

# Binaries.
GOLANGCI_LINT := $(abspath $(TOOLS_BIN_DIR)/golangci-lint)

# Release variables
STAGING_REGISTRY ?= ghcr.io/hivelocity/cluster-api-provider-hivelocity-staging
PROD_REGISTRY := ghcr.io/hivelocity/cluster-api-provider-hivelocity
IMG ?= $(STAGING_REGISTRY):latest


.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Binaries / Software

.PHONY: install-ctlptl
install-ctlptl: ## Installs CTLPTL (CLI for declaratively setting up local Kubernetes clusters)
	MINIMUM_CTLPTL_VERSION=$(MINIMUM_CTLPTL_VERSION) ./hack/ensure-ctlptl.sh

.PHONY: check-go
check-go: ## Checks go version
	MINIMUM_GO_VERSION=$(MINIMUM_GO_VERSION) ./hack/ensure-go.sh

install-kind: ## Installs Kind (Kubernetes-in-Docker)
	MINIMUM_KIND_VERSION=$(MINIMUM_KIND_VERSION) ./hack/ensure-kind.sh

.PHONY: install-kubectl
install-kubectl: ## Installs Kubectl (CLI for kubernetes)
	MINIMUM_KUBECTL_VERSION=$(MINIMUM_KUBECTL_VERSION) ./hack/ensure-kubectl.sh

.PHONY: install-tilt
install-tilt: ## Installs Tilt (watches files, builds containers, ships to k8s)
	MINIMUM_TILT_VERSION=$(MINIMUM_TILT_VERSION) ./hack/ensure-tilt.sh

.PHONY: install-clusterctl
install-clusterctl: ## Installs clusterctl
	MINIMUM_CLUSTERCTL_VERSION=$(MINIMUM_CLUSTERCTL_VERSION) ./hack/ensure-clusterctl.sh

install-dev-prerequisites: ## Installs all necessary dependencies
	@echo "Start checking dependencies"
	$(MAKE) install-ctlptl
	$(MAKE) check-go
	$(MAKE) install-kind
	$(MAKE) install-kubectl
	$(MAKE) install-tilt
	$(MAKE) install-clusterctl
	@echo "Finished: All dependencies up to date"

KUSTOMIZE := $(abspath $(TOOLS_BIN_DIR)/kustomize)
kustomize: $(KUSTOMIZE) ## Build a local copy of kustomize
$(KUSTOMIZE): # Build kustomize from tools folder.
	cd $(TOOLS_DIR) && go build -tags=tools -o $(KUSTOMIZE) sigs.k8s.io/kustomize/kustomize/v4

TILT := $(abspath $(TOOLS_BIN_DIR)/tilt)
tilt: $(TILT) ## Build a local copy of tilt
$(TILT):
	@mkdir -p $(TOOLS_BIN_DIR)
	MINIMUM_TILT_VERSION=0.31.2 hack/ensure-tilt.sh

ENVSUBST := $(abspath $(TOOLS_BIN_DIR)/envsubst)
envsubst: $(ENVSUBST) ## Build a local copy of envsubst
$(ENVSUBST): $(TOOLS_DIR)/go.mod # Build envsubst from tools folder.
	cd $(TOOLS_DIR) && go build -tags=tools -o $(ENVSUBST) github.com/drone/envsubst/v2/cmd/envsubst

CTLPTL := $(abspath $(TOOLS_BIN_DIR)/ctlptl)
ctlptl: $(CTLPTL) ## Build a local copy of ctlptl
$(CTLPTL):
	cd $(TOOLS_DIR) && go build -tags=tools -o $(CTLPTL) github.com/tilt-dev/ctlptl/cmd/ctlptl

install-crds: generate-manifests $(KUSTOMIZE) ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall-crds: generate-manifests $(KUSTOMIZE) ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

##@ Development

.PHONY: generate
generate: ## Run all generate-manifests, generate-go-deepcopyand generate-go-conversions targets
	$(MAKE) generate-manifests generate-go-deepcopy

generate-manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) \
			paths=./api/... \
			paths=./controllers/... \
			crd \
			rbac:roleName=manager-role \
			webhook \
			output:crd:artifacts:config=config/crd/bases

generate-go-deepcopy: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="./hack/boilerplate/boilerplate.generatego.txt" paths="./api/..."

dry-run: generate
	cd config/manager && $(KUSTOMIZE) edit set image controller=${STAGING_REGISTRY}:${TAG}
	mkdir -p dry-run
	$(KUSTOMIZE) build config/default > dry-run/manifests.yaml

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test -coverpkg=./... ./... -coverprofile cover.out

.PHONY: ensure-boilerplate
ensure-boilerplate: ## Ensures that a boilerplate exists in each file by adding missing boilerplates
	./hack/ensure-boilerplate.sh

##@ Build

.PHONY: build
build: generate fmt vet ## Build manager binary.
	go build -o bin/manager main.go

.PHONY: run
run: generate fmt vet ## Run a controller from your host.
	go run ./main.go

# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: test ## Build docker image with the manager.
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}

# PLATFORMS defines the target platforms for  the manager image be build to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - able to use docker buildx . More info: https://docs.docker.com/build/buildx/
# - have enable BuildKit, More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image for your registry (i.e. if you do not inform a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To properly provided solutions that supports more than one platform you should use this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: test ## Build and push docker image for the manager for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- docker buildx create --name project-v3-builder
	docker buildx use project-v3-builder
	- docker buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- docker buildx rm project-v3-builder
	rm Dockerfile.cross

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: generate-manifests $(KUSTOMIZE) ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: generate-manifests $(KUSTOMIZE) ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: generate-manifests $(KUSTOMIZE) ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Tool Versions
CONTROLLER_TOOLS_VERSION ?= v0.10.0

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary. If wrong version is installed, it will be overwritten.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen && $(LOCALBIN)/controller-gen --version | grep -q $(CONTROLLER_TOOLS_VERSION) || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest


##@ Lint and Verify

.PHONY: modules
modules: ## Runs go mod to ensure modules are up to date.
	go mod tidy
	cd $(TOOLS_DIR); go mod tidy

golangci-lint: $(GOLANGCI_LINT) ## Build a local copy of golangci-lint
$(GOLANGCI_LINT): .github/workflows/pr-golangci-lint.yml # Download golanci-lint using hack script into tools folder.
	hack/ensure-golangci-lint.sh \
		-b $(TOOLS_DIR)/$(BIN_DIR) \
		$(shell cat .github/workflows/pr-golangci-lint.yml | grep "version:.*golangci-lint/releases" | sed -r 's/.*version: (v[^ ]*) .*/\1/')

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint Golang codebase
	$(GOLANGCI_LINT) run -v $(GOLANGCI_LINT_EXTRA_ARGS)

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Lint the Go codebase and run auto-fixers if supported by the linter.
	GOLANGCI_LINT_EXTRA_ARGS=--fix $(MAKE) lint

.PHONY: format-tiltfile
format-tiltfile: ## Format the Tiltfile
	./hack/verify-starlark.sh fix

yamllint: ## Lints YAML Files
	yamllint -c .github/linters/yaml-lint.yaml --strict --format parsable .

ALL_VERIFY_CHECKS = boilerplate shellcheck tiltfile modules gen

.PHONY: verify
verify: lint yamllint $(addprefix verify-,$(ALL_VERIFY_CHECKS)) ## Run all verify-* targets
	@echo "All verify checks passed, congrats!"

.PHONY: verify-modules
verify-modules: modules  ## Verify go modules are up to date
	@if !(git diff --quiet HEAD -- go.sum go.mod $(TOOLS_DIR)/go.mod $(TOOLS_DIR)/go.sum $(TEST_DIR)/go.mod $(TEST_DIR)/go.sum); then \
		git diff; \
		echo "go module files are out of date"; exit 1; \
	fi
	@if (find . -name 'go.mod' | xargs -n1 grep -q -i 'k8s.io/client-go.*+incompatible'); then \
		find . -name "go.mod" -exec grep -i 'k8s.io/client-go.*+incompatible' {} \; -print; \
		echo "go module contains an incompatible client-go version"; exit 1; \
	fi

.PHONY: verify-gen
verify-gen: generate  ## Verfiy go generated files are up to date
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "generated files are out of date, run make generate"; exit 1; \
	fi

.PHONY: verify-boilerplate
verify-boilerplate: ## Verify boilerplate text exists in each file
	./hack/verify-boilerplate.sh

.PHONY: verify-shellcheck
verify-shellcheck: ## Verify shell files
	./hack/verify-shellcheck.sh

.PHONY: verify-tiltfile
verify-tiltfile: ## Verify Tiltfile format
	./hack/verify-starlark.sh


wait-and-get-secret:
	# Wait for the kubeconfig to become available.
	@test $${CLUSTER_NAME?Environment variable is required}
	${TIMEOUT} 5m bash -c "while ! kubectl get secrets | grep $(CLUSTER_NAME)-kubeconfig; do sleep 1; done"
	# Get kubeconfig and store it locally.
	kubectl get secrets $(CLUSTER_NAME)-kubeconfig -o json | jq -r .data.value | base64 --decode > $(CAPHV_WORKER_CLUSTER_KUBECONFIG)
	${TIMEOUT} 15m bash -c "while ! kubectl --kubeconfig=$(CAPHV_WORKER_CLUSTER_KUBECONFIG) get nodes | grep control-plane; do sleep 1; done"


install-manifests-cilium:
	# Deploy cilium
	helm repo add cilium https://helm.cilium.io/
	helm repo update cilium
	KUBECONFIG=$(CAPHV_WORKER_CLUSTER_KUBECONFIG) helm upgrade --install cilium cilium/cilium --version 1.12.2 \
  	--namespace kube-system \
	-f templates/cilium/cilium.yaml


install-manifests-ccm:
	# Deploy Hivelocity Cloud Controller Manager
	helm repo add hivelocity https://charts.hivelocity.com
	helm repo update hivelocity
	KUBECONFIG=$(CAPHV_WORKER_CLUSTER_KUBECONFIG) helm upgrade --install ccm hivelocity/ccm-hivelocity --version 1.0.11 \
	--namespace kube-system \
	--set secret.name=hivelocity \
	--set secret.key=hivelocity
	@echo 'run "kubectl --kubeconfig=$(CAPHV_WORKER_CLUSTER_KUBECONFIG) ..." to work with the new target cluster'


create-workload-cluster: $(KUSTOMIZE) $(ENVSUBST) ## Creates a workload-cluster. ENV Variables need to be exported or defined in the tilt-settings.yaml
	# Create workload Cluster. This is usually called via Tilt.
	@test $${CLUSTER_NAME?Environment variable is required}
	@test $${HIVELOCITY_API_KEY?Environment variable is required}
	kubectl create secret generic hivelocity --from-literal=hivelocity=$(HIVELOCITY_API_KEY) --save-config --dry-run=client -o yaml | kubectl apply -f -
	$(KUSTOMIZE) build templates/cluster-templates/hivelocity --load-restrictor LoadRestrictionsNone  > templates/cluster-templates/cluster-template-hivelocity.yaml
	cat templates/cluster-templates/cluster-template-hivelocity.yaml | $(ENVSUBST) - | kubectl apply -f -
	$(MAKE) wait-and-get-secret
	$(MAKE) install-manifests-cilium
	$(MAKE) install-manifests-ccm

.PHONY: cluster
cluster: $(CTLPTL) ## Creates kind-dev Cluster
	./hack/kind-dev.sh

.PHONY: delete-cluster
delete-cluster: $(CTLPTL) ## Deletes Kind-dev Cluster (default)
	$(CTLPTL) delete cluster kind-caphv

.PHONY: delete-registry
delete-registry: $(CTLPTL) ## Deletes Kind-dev Cluster and the local registry
	$(CTLPTL) delete registry caphv-registry

.PHONY: delete-cluster-registry
delete-cluster-registry: $(CTLPTL) ## Deletes Kind-dev Cluster and the local registry
	$(CTLPTL) delete cluster kind-caphv
	$(CTLPTL) delete registry caphv-registry

##@ Clean

.PHONY: clean
clean: ## Remove all generated files
	$(MAKE) clean-bin

.PHONY: clean-bin
clean-bin: ## Remove all generated helper binaries
	rm -rf $(BIN_DIR)
	rm -rf $(TOOLS_BIN_DIR)

.PHONY: clean-release
clean-release: ## Remove the release folder
	rm -rf $(RELEASE_DIR)

.PHONY: clean-release-git
clean-release-git: ## Restores the git files usually modified during a release
	git restore ./*manager_image_patch.yaml ./*manager_pull_policy.yaml

##@ Release

## latest git tag for the commit, e.g., v0.3.10
RELEASE_TAG ?= $(shell git describe --abbrev=0 2>/dev/null)
# the previous release tag, e.g., v0.3.9, excluding pre-release tags
PREVIOUS_TAG ?= $(shell git tag -l | grep -E "^v[0-9]+\.[0-9]+\.[0-9]." | sort -V | grep -B1 $(RELEASE_TAG) | head -n 1 2>/dev/null)
RELEASE_DIR ?= out
RELEASE_NOTES_DIR := _releasenotes

$(RELEASE_DIR):
	mkdir -p $(RELEASE_DIR)/

$(RELEASE_NOTES_DIR):
	mkdir -p $(RELEASE_NOTES_DIR)/

.PHONY: test-release
test-release:
	$(MAKE) set-manifest-image MANIFEST_IMG=$(PROD_REGISTRY) MANIFEST_TAG=$(TAG)
	$(MAKE) set-manifest-pull-policy PULL_POLICY=IfNotPresent
	$(MAKE) release-manifests

.PHONY: release
release: clean-release  ## Builds and push container images using the latest git tag for the commit.
	@if [ -z "${RELEASE_TAG}" ]; then echo "RELEASE_TAG is not set"; exit 1; fi
	@if ! [ -z "$$(git status --porcelain)" ]; then echo "Your local git repository contains uncommitted changes, use git clean before proceeding."; exit 1; fi
	git checkout "${RELEASE_TAG}"
	# Set the manifest image to the production bucket.
	$(MAKE) set-manifest-image MANIFEST_IMG=$(PROD_REGISTRY) MANIFEST_TAG=$(RELEASE_TAG)
	$(MAKE) set-manifest-pull-policy PULL_POLICY=IfNotPresent
	## Build the manifests
	$(MAKE) release-manifests clean-release-git

.PHONY: release-manifests
release-manifests: generate $(KUSTOMIZE) $(RELEASE_DIR) cluster-templates ## Builds the manifests to publish with a release
	$(KUSTOMIZE) build config/default > $(RELEASE_DIR)/infrastructure-components.yaml
	## Build caphv-components (aggregate of all of the above).
	cp metadata.yaml $(RELEASE_DIR)/metadata.yaml
	cp templates/cluster-templates/cluster-template* $(RELEASE_DIR)/

.PHONY: release-notes
release-notes: $(RELEASE_NOTES_DIR) $(RELEASE_NOTES)
	go run ./hack/tools/release/notes.go --from=$(PREVIOUS_TAG) > $(RELEASE_NOTES_DIR)/$(RELEASE_TAG).md

.PHONY: set-manifest-image
set-manifest-image:
	$(info Updating kustomize image patch file for default resource)
	sed -i'' -e 's@image: .*@image: '"${MANIFEST_IMG}:$(MANIFEST_TAG)"'@' ./config/default/manager_image_patch.yaml

.PHONY: set-manifest-pull-policy
set-manifest-pull-policy:
	$(info Updating kustomize pull policy file for default resource)
	sed -i'' -e 's@imagePullPolicy: .*@imagePullPolicy: '"$(PULL_POLICY)"'@' ./config/default/manager_pull_policy.yaml


.PHONY: tilt-up
tilt-up: $(ENVSUBST) $(KUSTOMIZE) $(TILT) cluster  ## Start a mgt-cluster & Tilt. Installs the CRDs and deploys the controllers
	EXP_CLUSTER_RESOURCE_SET=true $(TILT) up
