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

INFRA_SHORT = caphv

# TODO: change "syself" to "hivelocity", when we have the perms for uploading
IMAGE_PREFIX ?= ghcr.io/hivelocity

INFRA_PROVIDER = hivelocity

STAGING_IMAGE = $(INFRA_SHORT)-staging
BUILDER_IMAGE = ghcr.io/hivelocity/$(INFRA_SHORT)-builder
BUILDER_IMAGE_VERSION = $(shell cat .builder-image-version.txt)

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec
.DEFAULT_GOAL:=help
GOTEST ?= go test

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

#############
# Variables #
#############

# Certain aspects of the build are done in containers for consistency (e.g. protobuf generation)
# If you have the correct tools installed and you want to speed up development you can run
# make BUILD_IN_CONTAINER=false target
# or you can override this with an environment variable
BUILD_IN_CONTAINER ?= true

# Boiler plate for building Docker containers.
TAG ?= dev
ARCH ?= amd64
# Allow overriding the imagePullPolicy
PULL_POLICY ?= Always
# Build time versioning details.
LDFLAGS := $(shell hack/version.sh)

TIMEOUT := $(shell command -v timeout || command -v gtimeout)

# Directories
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
EXP_DIR := exp
TEST_DIR := test
BIN_DIR := bin
TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/$(BIN_DIR)
export PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)
export GOBIN := $(abspath $(TOOLS_BIN_DIR))

# Files
WORKER_CLUSTER_KUBECONFIG ?= ".workload-cluster-kubeconfig.yaml"
MGT_CLUSTER_KUBECONFIG ?= ".mgt-cluster-kubeconfig.yaml"

# Kubebuilder.
export KUBEBUILDER_ENVTEST_KUBERNETES_VERSION ?= 1.28.0

##@ Binaries
############
# Binaries #
############
CONTROLLER_GEN := $(abspath $(TOOLS_BIN_DIR)/controller-gen)
controller-gen: $(CONTROLLER_GEN) ## Build a local copy of controller-gen
$(CONTROLLER_GEN): # Build controller-gen from tools folder.
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.14.0

KUSTOMIZE := $(abspath $(TOOLS_BIN_DIR)/kustomize)
kustomize: $(KUSTOMIZE) ## Build a local copy of kustomize
$(KUSTOMIZE): # Build kustomize from tools folder.
	go install sigs.k8s.io/kustomize/kustomize/v4@v4.5.7

TILT := $(abspath $(TOOLS_BIN_DIR)/tilt)
tilt: $(TILT) ## Build a local copy of tilt
$(TILT):
	@mkdir -p $(TOOLS_BIN_DIR)
	MINIMUM_TILT_VERSION=0.33.3 hack/ensure-tilt.sh

ENVSUBST := $(abspath $(TOOLS_BIN_DIR)/envsubst)
envsubst: $(ENVSUBST) ## Build a local copy of envsubst
$(ENVSUBST): # Build envsubst from tools folder.
	go install github.com/drone/envsubst/v2/cmd/envsubst@latest

SETUP_ENVTEST := $(abspath $(TOOLS_BIN_DIR)/setup-envtest)
setup-envtest: $(SETUP_ENVTEST) ## Build a local copy of setup-envtest
$(SETUP_ENVTEST): # Build setup-envtest from tools folder.
	go install sigs.k8s.io/controller-runtime/tools/setup-envtest@v0.0.0-20231206145619-1ea2be573f78

CTLPTL := $(abspath $(TOOLS_BIN_DIR)/ctlptl)
ctlptl: $(CTLPTL) ## Build a local copy of ctlptl
$(CTLPTL):
	go install github.com/tilt-dev/ctlptl/cmd/ctlptl@v0.8.25

CLUSTERCTL := $(abspath $(TOOLS_BIN_DIR)/clusterctl)
clusterctl: $(CLUSTERCTL) ## Build a local copy of clusterctl
$(CLUSTERCTL):
	curl -sSLf https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.6.0/clusterctl-$$(go env GOOS)-$$(go env GOARCH) -o $(CLUSTERCTL)
	chmod a+rx $(CLUSTERCTL)

HELM := $(abspath $(TOOLS_BIN_DIR)/helm)
helm: $(HELM) ## Build a local copy of helm
$(HELM):
	curl -sSL https://get.helm.sh/helm-v3.13.2-linux-amd64.tar.gz | tar xz -C $(TOOLS_BIN_DIR) --strip-components=1 linux-amd64/helm
	chmod a+rx $(HELM)
KIND := $(abspath $(TOOLS_BIN_DIR)/kind)
kind: $(KIND) ## Build a local copy of kind
$(KIND):
	go install sigs.k8s.io/kind@v0.20.0

KUBECTL := $(abspath $(TOOLS_BIN_DIR)/kubectl)
kubectl: $(KUBECTL) ## Build a local copy of kubectl
$(KUBECTL):
	curl -fsSL "https://dl.k8s.io/release/v1.27.3/bin/$$(go env GOOS)/$$(go env GOARCH)/kubectl" -o $(KUBECTL)
	chmod a+rx $(KUBECTL)

go-binsize-treemap := $(abspath $(TOOLS_BIN_DIR)/go-binsize-treemap)
go-binsize-treemap: $(go-binsize-treemap) # Build go-binsize-treemap from tools folder.
$(go-binsize-treemap):
	go install github.com/nikolaydubina/go-binsize-treemap@v0.2.0

go-cover-treemap := $(abspath $(TOOLS_BIN_DIR)/go-cover-treemap)
go-cover-treemap: $(go-cover-treemap) # Build go-cover-treemap from tools folder.
$(go-cover-treemap):
	go install github.com/nikolaydubina/go-cover-treemap@v1.3.0

GOTESTSUM := $(abspath $(TOOLS_BIN_DIR)/gotestsum)
gotestsum: $(GOTESTSUM) # Build gotestsum from tools folder.
$(GOTESTSUM):
	go install gotest.tools/gotestsum@v1.11.0

VIDDY := $(abspath $(TOOLS_BIN_DIR)/viddy)
viddy: $(VIDDY)
$(VIDDY):
	go install github.com/sachaos/viddy@latest

all-tools: $(GOTESTSUM) $(go-cover-treemap) $(go-binsize-treemap) $(KIND) $(KUBECTL) $(CLUSTERCTL) $(CTLPTL) $(SETUP_ENVTEST) $(ENVSUBST) $(KUSTOMIZE) $(CONTROLLER_GEN)
	echo 'done'

##@ Development
###############
# Development #
###############
install-crds: generate-manifests $(KUSTOMIZE) ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) apply -f -


uninstall-crds: generate-manifests $(KUSTOMIZE) ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) delete -f -

deploy-controller: generate-manifests $(KUSTOMIZE) ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMAGE_PREFIX}/$(STAGING_IMAGE):${TAG}
	$(KUSTOMIZE) build config/default | $(KUBECTL) apply -f -

undeploy-controller: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | $(KUBECTL) delete -f -

install-essentials: ## This gets the secret and installs a CNI and the CCM. Usage: MAKE install-essentials
	$(MAKE) wait-and-get-secret
	$(MAKE) install-cilium-in-wl-cluster
	$(MAKE) install-ccm-in-wl-cluster

wait-and-get-secret:
	# Wait for the kubeconfig to become available.
	${TIMEOUT} --foreground 5m bash -c "while ! $(KUBECTL) get secrets | grep $(CLUSTER_NAME)-kubeconfig; do sleep 1; done"
	# Get kubeconfig and store it locally.
	$(KUBECTL) get secrets $(CLUSTER_NAME)-kubeconfig -o json | jq -r .data.value | base64 --decode > $(WORKER_CLUSTER_KUBECONFIG)
	${TIMEOUT} --foreground 22m bash -c "while ! $(KUBECTL) --kubeconfig=$(WORKER_CLUSTER_KUBECONFIG) get nodes | grep control-plane; do sleep 5; done"

install-cilium-in-wl-cluster: $(HELM)
	# Deploy cilium
	$(HELM) repo add cilium https://helm.cilium.io/
	$(HELM) repo update cilium
	KUBECONFIG=$(WORKER_CLUSTER_KUBECONFIG) $(HELM) upgrade --install cilium cilium/cilium --version 1.14.4 \
  	--namespace kube-system \
	-f templates/cilium/cilium.yaml

install-ccm-in-wl-cluster:
	$(HELM) repo add $(INFRA_PROVIDER) https://hivelocity.github.io/hivelocity-cloud-controller-manager/
	$(HELM) repo update $(INFRA_PROVIDER)
	KUBECONFIG=$(WORKER_CLUSTER_KUBECONFIG) $(HELM) upgrade --install ccm-hivelocity hivelocity/ccm-hivelocity \
		--version 0.1.2 \
		--namespace kube-system \
		--set secret.name=$(INFRA_PROVIDER) \
		--set secret.key=$(INFRA_PROVIDER)
	@echo 'run "kubectl --kubeconfig=$(WORKER_CLUSTER_KUBECONFIG) ..." to work with the new target cluster'

add-ssh-pub-key:
	@./hack/ensure-env-variables.sh HIVELOCITY_API_KEY HIVELOCITY_SSH_KEY
	go run ./cmd upload-ssh-pub-key $$HIVELOCITY_SSH_KEY $(HOME)/.ssh/$(INFRA_PROVIDER).pub

env-vars-for-wl-cluster:
	@if [ -n "$$HIVELOCITY_WORKER_MACHINE_TYPE" ]; then echo "please rename HIVELOCITY_WORKER_MACHINE_TYPE to HIVELOCITY_WORKER_DEVICE_TYPE"; exit 1; fi
	@if [ -n "$$HIVELOCITY_CONTROL_PLANE_MACHINE_TYPE" ]; then echo "please rename HIVELOCITY_CONTROL_PLANE_MACHINE_TYPE to HIVELOCITY_CONTROL_PLANE_DEVICE_TYPE"; exit 1; fi

	@./hack/ensure-env-variables.sh CLUSTER_NAME CONTROL_PLANE_MACHINE_COUNT HIVELOCITY_CONTROL_PLANE_DEVICE_TYPE \
	HIVELOCITY_API_KEY HIVELOCITY_SSH_KEY HIVELOCITY_WORKER_DEVICE_TYPE KUBERNETES_VERSION WORKER_MACHINE_COUNT \
	HIVELOCITY_REGION
	@hack/check-kubernetes-version.sh


	@regex="^[-A-Za-z0-9_.]*$$"; if [[ ! $$HIVELOCITY_WORKER_DEVICE_TYPE =~ $$regex ]]; then \
		echo "HIVELOCITY_WORKER_DEVICE_TYPE=$$HIVELOCITY_WORKER_DEVICE_TYPE needs to be a valid Kubernetes label value." ;\
		exit 1 ;\
	fi
	@regex="^[-A-Za-z0-9_.]*$$"; if [[ ! $$HIVELOCITY_CONTROL_PLANE_DEVICE_TYPE =~ $$regex ]]; then \
		echo "HIVELOCITY_CONTROL_PLANE_DEVICE_TYPE=$$HIVELOCITY_CONTROL_PLANE_DEVICE_TYPE needs to be a valid Kubernetes label value." ;\
		exit 1 ;\
	fi

create-workload-cluster: env-vars-for-wl-cluster $(HOME)/.ssh/$(INFRA_PROVIDER).pub $(CLUSTERCTL) $(KUSTOMIZE) $(ENVSUBST) install-crds ## Creates a workload-cluster.
	# Create workload Cluster.
	rm -f $(WORKER_CLUSTER_KUBECONFIG)
	go run ./cmd upload-ssh-pub-key $$HIVELOCITY_SSH_KEY $(HOME)/.ssh/$(INFRA_PROVIDER).pub

	# If the secret already exists, then it is likely that the cluster is already running,
	# and the user wants to connect to the running cluster.
	# In this case, don't remove the labels from the machines, otherwise the running cluster will be broken.
	$(KUBECTL) get secret $(INFRA_PROVIDER) >/dev/null 2>&1 || \
	 	go run ./test/claim-devices-or-fail $$HIVELOCITY_CONTROL_PLANE_DEVICE_TYPE $$HIVELOCITY_WORKER_DEVICE_TYPE

	$(KUBECTL) create secret generic $(INFRA_PROVIDER) --from-literal=$(INFRA_PROVIDER)=$(HIVELOCITY_API_KEY) --save-config --dry-run=client -o yaml | $(KUBECTL) apply -f -
	$(KUSTOMIZE) build templates/cluster-templates/$(INFRA_PROVIDER) --load-restrictor LoadRestrictionsNone  > templates/cluster-templates/cluster-template-$(INFRA_PROVIDER).yaml
	cat templates/cluster-templates/cluster-template-$(INFRA_PROVIDER).yaml | $(ENVSUBST) - > templates/cluster-templates/cluster-template-$(INFRA_PROVIDER).yaml.apply
	$(KUBECTL) apply -f templates/cluster-templates/cluster-template-$(INFRA_PROVIDER).yaml.apply
	$(MAKE) wait-and-get-secret
	$(MAKE) install-cilium-in-wl-cluster
	$(MAKE) install-ccm-in-wl-cluster

move-to-workload-cluster: $(CLUSTERCTL)
	$(CLUSTERCTL) init --kubeconfig=$(WORKER_CLUSTER_KUBECONFIG) --core cluster-api --bootstrap kubeadm --control-plane kubeadm --infrastructure $(INFRA_PROVIDER)
	$(KUBECTL) --kubeconfig=$(WORKER_CLUSTER_KUBECONFIG) -n $(INFRA_SHORT)-system wait deploy/$(INFRA_SHORT)-controller-manager --for condition=available && sleep 15s
	$(CLUSTERCTL) move --to-kubeconfig=$(WORKER_CLUSTER_KUBECONFIG)

.PHONY: delete-workload-cluster
delete-workload-cluster: ## Deletes the example workload Kubernetes cluster
	./hack/ensure-env-variables.sh CLUSTER_NAME
	@echo 'Your workload cluster will now be deleted, this can take up to 20 minutes'
	$(KUBECTL) patch cluster $(CLUSTER_NAME) --type=merge -p '{"spec":{"paused": false}}'
	$(KUBECTL) delete cluster $(CLUSTER_NAME)
	${TIMEOUT} --foreground 22m bash -c "while $(KUBECTL) get cluster | grep $(NAME); do sleep 1; done"
	@echo 'Cluster deleted'

create-mgt-cluster: $(CLUSTERCTL) $(KUBECTL) cluster ## Start a mgt-cluster with the latest version of all capi components and the infra provider.
	#TODO: activate after official release of "hivelocity" $(CLUSTERCTL) init --core cluster-api --bootstrap kubeadm --control-plane kubeadm --infrastructure $(INFRA_PROVIDER)
	$(KUBECTL) create secret generic $(INFRA_PROVIDER) --from-literal=$(INFRA_PROVIDER)=$(HIVELOCITY_API_KEY)
	$(KUBECTL) patch secret $(INFRA_PROVIDER) -p '{"metadata":{"labels":{"clusterctl.cluster.x-k8s.io/move":""}}}'

.PHONY: cluster
cluster: $(CTLPTL) $(KUBECTL) ## Creates kind-dev Cluster
	@# Fail early: Test if HIVELOCITY_API_KEY is set.
	./hack/ensure-env-variables.sh HIVELOCITY_API_KEY
	./hack/kind-dev.sh

.PHONY: delete-mgt-cluster
delete-mgt-cluster: $(CTLPTL) ## Deletes Kind-dev Cluster (default)
	$(CTLPTL) delete cluster kind-$(INFRA_SHORT)

.PHONY: delete-registry
delete-registry: $(CTLPTL) ## Deletes Kind-dev Cluster and the local registry
	$(CTLPTL) delete registry $(INFRA_SHORT)-registry

.PHONY: delete-mgt-cluster-registry
delete-mgt-cluster-registry: $(CTLPTL) ## Deletes Kind-dev Cluster and the local registry
	$(CTLPTL) delete cluster kind-$(INFRA_SHORT)
	$(CTLPTL) delete registry $(INFRA_SHORT)-registry

##@ Clean
#########
# Clean #
#########
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
	git restore ./*manager_config_patch.yaml ./*manager_pull_policy.yaml

##@ Releasing
#############
# Releasing #
#############
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
	$(MAKE) set-manifest-image MANIFEST_IMG=$(IMAGE_PREFIX)/$(STAGING_IMAGE) MANIFEST_TAG=$(TAG)
	$(MAKE) set-manifest-pull-policy PULL_POLICY=IfNotPresent
	$(MAKE) release-manifests

.PHONY: release-manifests
release-manifests: generate-manifests generate-go-deepcopy $(KUSTOMIZE) $(RELEASE_DIR) cluster-templates ## Builds the manifests to publish with a release
	$(KUSTOMIZE) build config/default > $(RELEASE_DIR)/infrastructure-components.yaml
	## Build $(INFRA_SHORT)-components (aggregate of all of the above).
	cp metadata.yaml $(RELEASE_DIR)/metadata.yaml
	cp templates/cluster-templates/cluster-template* $(RELEASE_DIR)/

.PHONY: release
release: clean-release  ## Builds and push container images using the latest git tag for the commit.
	@if [ -z "${RELEASE_TAG}" ]; then echo "RELEASE_TAG is not set"; exit 1; fi
	@if ! [ -z "$$(git status --porcelain)" ]; then echo "Your local git repository contains uncommitted changes, use git clean before proceeding."; exit 1; fi
	git checkout "${RELEASE_TAG}"
	# Set the manifest image to the production bucket.
	$(MAKE) set-manifest-image MANIFEST_IMG=$(IMAGE_PREFIX)/$(INFRA_SHORT) MANIFEST_TAG=$(RELEASE_TAG)
	$(MAKE) set-manifest-pull-policy PULL_POLICY=IfNotPresent
	## Build the manifests
	$(MAKE) release-manifests clean-release-git

.PHONY: release-notes
release-notes: $(RELEASE_NOTES_DIR) $(RELEASE_NOTES)
	go run ./hack/tools/release/notes.go --from=$(PREVIOUS_TAG) > $(RELEASE_NOTES_DIR)/$(RELEASE_TAG).md

##@ Images
##########
# Images #
##########

.PHONY: set-manifest-image
set-manifest-image:
	$(info Updating kustomize image patch file for default resource)
	sed -i'' -e 's@image: .*@image: '"${MANIFEST_IMG}:$(MANIFEST_TAG)"'@' ./config/default/manager_config_patch.yaml

.PHONY: set-manifest-pull-policy
set-manifest-pull-policy:
	$(info Updating kustomize pull policy file for default resource)
	sed -i'' -e 's@imagePullPolicy: .*@imagePullPolicy: '"$(PULL_POLICY)"'@' ./config/default/manager_pull_policy.yaml

builder-image-promote-latest:
	./hack/ensure-env-variables.sh USERNAME PASSWORD
	skopeo copy --src-creds=$(USERNAME):$(PASSWORD) --dest-creds=$(USERNAME):$(PASSWORD) \
		docker://$(BUILDER_IMAGE):$(BUILDER_IMAGE_VERSION) \
		docker://$(BUILDER_IMAGE):latest

##@ Binary
##########
# Binary #
##########
$(INFRA_SHORT): ## Build controller binary.
	go build -mod=vendor -o bin/manager main.go

run: ## Run a controller from your host.
	go run ./main.go

##@ Testing
###########
# Testing #
###########
ARTIFACTS ?= _artifacts
$(ARTIFACTS):
	mkdir -p $(ARTIFACTS)/


$(MGT_CLUSTER_KUBECONFIG):
	./hack/get-kubeconfig-of-management-cluster.sh

$(WORKER_CLUSTER_KUBECONFIG):
	./hack/get-kubeconfig-of-workload-cluster.sh

.PHONY: k9s-workload-cluster
k9s-workload-cluster: $(WORKER_CLUSTER_KUBECONFIG)
	KUBECONFIG=$(WORKER_CLUSTER_KUBECONFIG) k9s

.PHONY: bash-with-kubeconfig-set-to-workload-cluster
bash-with-kubeconfig-set-to-workload-cluster: $(WORKER_CLUSTER_KUBECONFIG)
	KUBECONFIG=$(WORKER_CLUSTER_KUBECONFIG) bash

.PHONY: tail-controller-logs
tail-controller-logs: ## Show the last lines of the controller logs
	@hack/tail-controller-logs.sh

.PHONY: ssh-first-control-plane
ssh-first-control-plane: ## ssh into the first control-plane
	@hack/ssh-first-control-plane.sh


KUBEBUILDER_ASSETS ?= $(shell $(SETUP_ENVTEST) use --use-env --bin-dir $(abspath $(TOOLS_BIN_DIR)) -p path $(KUBEBUILDER_ENVTEST_KUBERNETES_VERSION))

E2E_DIR ?= $(ROOT_DIR)/test/e2e
E2E_CONF_FILE_SOURCE ?= $(E2E_DIR)/config/$(INFRA_PROVIDER).yaml
E2E_CONF_FILE ?= $(E2E_DIR)/config/$(INFRA_PROVIDER)-ci-envsubst.yaml

.PHONY: test-unit
test-unit: $(SETUP_ENVTEST) $(GOTESTSUM) ## Run unit and integration tests
	@mkdir -p $(shell pwd)/.coverage
	KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" $(GOTESTSUM) --junitfile=.coverage/junit.xml --format testname -- -covermode=atomic -coverprofile=.coverage/cover.out -p=4 ./controllers/... ./pkg/... ./api/...

.PHONY: e2e-image
e2e-image: ## Build the e2e manager image
	docker build --pull --build-arg ARCH=$(ARCH) --build-arg LDFLAGS="$(LDFLAGS)" -t $(IMAGE_PREFIX)/$(STAGING_IMAGE):e2e -f images/$(INFRA_SHORT)/Dockerfile .

.PHONY: $(E2E_CONF_FILE)
e2e-conf-file: $(E2E_CONF_FILE)
$(E2E_CONF_FILE): $(ENVSUBST) $(E2E_CONF_FILE_SOURCE)
	mkdir -p $(shell dirname $(E2E_CONF_FILE))
	MANAGEMENT_CLUSTER_NAME="$(INFRA_SHORT)-e2e-$$(date +"%Y%m%d-%H%M%S")-$$USER" $(ENVSUBST) < $(E2E_CONF_FILE_SOURCE) > $(E2E_CONF_FILE)

.PHONY: test-e2e
test-e2e: $(E2E_CONF_FILE) $(if $(SKIP_IMAGE_BUILD),,e2e-image) $(ARTIFACTS)
	rm -f $(WORKER_CLUSTER_KUBECONFIG)
	GINKGO_FOKUS="'\[Basic\]'" GINKGO_NODES=2 E2E_CONF_FILE=$(E2E_CONF_FILE) ./hack/ci-e2e-capi.sh


##@ Report
##########
# Report #
##########
report-cover-html: ## Create a html report
	@mkdir -p $(shell pwd)/.reports
	go tool cover -html .coverage/cover.out -o .reports/coverage.html

report-binsize-treemap: $(go-binsize-treemap) ## Creates a treemap of the binary
	@mkdir -p $(shell pwd)/.reports
	go tool nm -size bin/manager | $(go-binsize-treemap) -w 1024 -h 256 > .reports/$(INFRA_SHORT)-binsize-treemap-sm.svg
	go tool nm -size bin/manager | $(go-binsize-treemap) -w 1024 -h 1024 > .reports/$(INFRA_SHORT)-binsize-treemap.svg
	go tool nm -size bin/manager | $(go-binsize-treemap) -w 2048 -h 2048 > .reports/$(INFRA_SHORT)-binsize-treemap-lg.svg

report-binsize-treemap-all: $(go-binsize-treemap) report-binsize-treemap
	@mkdir -p $(shell pwd)/.reports
	go tool nm -size bin/manager | $(go-binsize-treemap) -w 4096 -h 4096 > .reports/$(INFRA_SHORT)-binsize-treemap-xl.svg
	go tool nm -size bin/manager | $(go-binsize-treemap) -w 8192 -h 8192 > .reports/$(INFRA_SHORT)-binsize-treemap-xxl.svg

report-cover-treemap: $(go-cover-treemap) ## Creates a treemap of the coverage
	@mkdir -p $(shell pwd)/.reports
	$(go-cover-treemap) -w 1080 -h 360 -coverprofile .coverage/cover.out > .reports/$(INFRA_SHORT)-cover-treemap-sm.svg
	$(go-cover-treemap) -w 2048 -h 1280 -coverprofile .coverage/cover.out > .reports/$(INFRA_SHORT)-cover-treemap-lg.svg
	$(go-cover-treemap) --only-folders -coverprofile .coverage/cover.out > .reports/$(INFRA_SHORT)-cover-treemap-folders.svg

##@ Verify
##########
# Verify #
##########


.PHONY: verify-modules
verify-modules: modules  ## Verify go modules are up to date
	@if !(git diff --quiet HEAD -- go.sum go.mod); then \
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

.PHONY: verify-starlark
verify-starlark: ## Verify Starlark Code
	./hack/verify-starlark.sh

.PHONY: verify-manifests ## Verify Manifests
verify-manifests:
	./hack/verify-manifests.sh
.PHONY: verify-container-images
verify-container-images: ## Verify container images
	trivy image -q --exit-code 1 --ignore-unfixed --severity MEDIUM,HIGH,CRITICAL $(IMAGE_PREFIX)/$(INFRA_SHORT):latest

##@ Generate
############
# Generate #
############
.PHONY: generate-boilerplate
generate-boilerplate: ## Generates missing boilerplates
	./hack/ensure-boilerplate.sh

# support go modules
generate-modules: ## Generates missing go modules
ifeq ($(BUILD_IN_CONTAINER),true)
	docker run  --rm \
		-v $(shell go env GOPATH)/pkg:/go/pkg \
		-v $(shell pwd):/src \
		$(BUILDER_IMAGE):$(BUILDER_IMAGE_VERSION) $@;
else
	./hack/golang-modules-update.sh
endif

$(HOME)/.ssh/$(INFRA_PROVIDER).pub:
	echo "Creating SSH key-pair to access the nodes which get created by $(INFRA_PROVIDER)"
	ssh-keygen -f ~/.ssh/$(INFRA_PROVIDER)

generate-modules-ci: generate-modules
	@if ! (git diff --exit-code ); then \
		echo "\nChanges found in generated files"; \
		exit 1; \
	fi

generate-manifests: $(CONTROLLER_GEN) ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) \
			paths=./api/... \
			paths=./controllers/... \
			crd:crdVersions=v1 \
			rbac:roleName=manager-role \
			output:crd:dir=./config/crd/bases \
			output:webhook:dir=./config/webhook \
			webhook

generate-go-deepcopy: $(CONTROLLER_GEN) ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) \
		object:headerFile="./hack/boilerplate/boilerplate.generatego.txt" \
		paths="./api/..."

generate-api-ci: generate-manifests generate-go-deepcopy
	@if ! (git diff --exit-code ); then \
		echo "\nChanges found in generated files"; \
		exit 1; \
	fi

cluster-templates: $(KUSTOMIZE)
	$(KUSTOMIZE) build templates/cluster-templates/$(INFRA_PROVIDER) --load-restrictor LoadRestrictionsNone  > templates/cluster-templates/cluster-template.yaml

##@ Format
##########
# Format #
##########
.PHONY: format-golang
format-golang: ## Format the Go codebase and run auto-fixers if supported by the linter.
ifeq ($(BUILD_IN_CONTAINER),true)
	docker run  --rm -t -i \
		-v $(shell go env GOPATH)/pkg:/go/pkg \
		-v $(shell pwd):/src \
		$(BUILDER_IMAGE):$(BUILDER_IMAGE_VERSION) $@;
else
	go version
	golangci-lint version
	golangci-lint run --fix
endif

.PHONY: format-starlark
format-starlark: ## Format the Starlark codebase
	./hack/verify-starlark.sh fix

.PHONY: format-yaml
format-yaml: ## Lint YAML files
ifeq ($(BUILD_IN_CONTAINER),true)
	docker run  --rm -t -i \
		-v $(shell go env GOPATH)/pkg:/go/pkg \
		-v $(shell pwd):/src \
		$(BUILDER_IMAGE):$(BUILDER_IMAGE_VERSION) $@;
else
	yamlfixer --version
	yamlfixer -c .yamllint.yaml .
endif

##@ Lint
########
# Lint #
########
.PHONY: lint-golang
lint-golang: ## Lint Golang codebase
ifeq ($(BUILD_IN_CONTAINER),true)
	docker run  --rm -t -i \
		-v $(shell go env GOPATH)/pkg:/go/pkg \
		-v $(shell pwd):/src \
		$(BUILDER_IMAGE):$(BUILDER_IMAGE_VERSION) $@;
else
	go version
	golangci-lint version
	golangci-lint run
endif

.PHONY: lint-golang-ci
lint-golang-ci:
ifeq ($(BUILD_IN_CONTAINER),true)
	docker run  --rm -t -i \
		-v $(shell go env GOPATH)/pkg:/go/pkg \
		-v $(shell pwd):/src \
		$(BUILDER_IMAGE):$(BUILDER_IMAGE_VERSION) $@;
else
	go version
	golangci-lint version
	golangci-lint run --out-format=github-actions
endif

.PHONY: lint-yaml
lint-yaml: ## Lint YAML files
ifeq ($(BUILD_IN_CONTAINER),true)
	docker run  --rm -t -i \
		-v $(shell go env GOPATH)/pkg:/go/pkg \
		-v $(shell pwd):/src \
		$(BUILDER_IMAGE):$(BUILDER_IMAGE_VERSION) $@;
else
	yamllint --version
	yamllint -c .yamllint.yaml --strict .
endif

.PHONY: lint-yaml-ci
lint-yaml-ci:
ifeq ($(BUILD_IN_CONTAINER),true)
	docker run  --rm -t -i \
		-v $(shell go env GOPATH)/pkg:/go/pkg \
		-v $(shell pwd):/src \
		$(BUILDER_IMAGE):$(BUILDER_IMAGE_VERSION) $@;
else
	yamllint --version
	yamllint -c .yamllint.yaml . --format github
endif

DOCKERFILES=$(shell find . -not \( -path ./hack -prune \) -not \( -path ./vendor -prune \) -type f -regex ".*Dockerfile.*"  | tr '\n' ' ')
.PHONY: lint-dockerfile
lint-dockerfile: ## Lint Dockerfiles
ifeq ($(BUILD_IN_CONTAINER),true)
	docker run  --rm -t -i \
		-v $(shell go env GOPATH)/pkg:/go/pkg \
		-v $(shell pwd):/src \
		$(BUILDER_IMAGE):$(BUILDER_IMAGE_VERSION) $@;
else
	hadolint --version
	hadolint -t error $(DOCKERFILES)
endif

lint-links: ## Link Checker
ifeq ($(BUILD_IN_CONTAINER),true)
	docker run --rm -t -i \
		-v $(shell pwd):/src \
		$(BUILDER_IMAGE):$(BUILDER_IMAGE_VERSION) $@;
else
	lychee --config .lychee.toml ./*.md  ./docs/**/*.md
endif

##@ Main Targets
################
# Main Targets #
################
.PHONY: lint
lint: lint-golang lint-yaml lint-dockerfile lint-links ## Lint Codebase

.PHONY: format
format: format-starlark format-golang format-yaml ## Format Codebase

.PHONY: generate
generate: generate-manifests generate-go-deepcopy generate-boilerplate generate-modules ## Generate Files

ALL_VERIFY_CHECKS = boilerplate shellcheck starlark manifests
.PHONY: verify
verify: generate lint $(addprefix verify-,$(ALL_VERIFY_CHECKS)) ## Verify all

.PHONY: modules
modules: generate-modules ## Update go.mod & go.sum

.PHONY: boilerplate
boilerplate: generate-boilerplate ## Ensure that your files have a boilerplate header

.PHONY: builder-image-push
builder-image-push: ## Build $(INFRA_SHORT)-builder to a new version. For more information see README.
	BUILDER_IMAGE=$(BUILDER_IMAGE) ./hack/upgrade-builder-image.sh

.PHONY: test
test: test-unit ## Runs all unit and integration tests.

.PHONY: tilt-up
tilt-up: env-vars-for-wl-cluster $(ENVSUBST) $(KUBECTL) $(KUSTOMIZE) $(TILT) cluster  ## Start a mgt-cluster & Tilt. Installs the CRDs and deploys the controllers
	EXP_CLUSTER_RESOURCE_SET=true $(TILT) up --port=10352

.PHONY: watch
watch: $(VIDDY) ## Watch CRDs cluster, machines and Events.
	$(VIDDY) -n 3 hack/output-for-watch.sh
