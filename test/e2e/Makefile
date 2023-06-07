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


# If you update this file, please follow:
# https://suva.sh/posts/well-documented-makefiles/

# Use GOPROXY environment variable if set

.DEFAULT_GOAL:=help

GOPROXY := $(shell go env GOPROXY)
ifeq ($(GOPROXY),)
GOPROXY := https://proxy.golang.org
endif
export GOPROXY

REPO_ROOT := $(shell git rev-parse --show-toplevel)

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## --------------------------------------
## Binaries
## --------------------------------------

TOOLS_DIR := $(REPO_ROOT)/hack/tools
BIN_DIR := bin
TOOLS_BIN_DIR := $(TOOLS_DIR)/$(BIN_DIR)
GINKGO := $(abspath $(TOOLS_BIN_DIR)/ginkgo)
KUSTOMIZE := $(TOOLS_BIN_DIR)/kustomize

ginkgo: $(GINKGO)
$(GINKGO): # Build ginkgo from tools folder.
	cd $(TOOLS_DIR) && go build -mod=vendor -tags=tools -o $(GINKGO) github.com/onsi/ginkgo/v2/ginkgo

$(KUSTOMIZE): # Build kustomize from tools folder.
	KUSTOMIZE_VERSION=4.5.7	$(REPO_ROOT)/hack/ensure-kustomize.sh

## --------------------------------------
## Templates
## --------------------------------------

HV_TEMPLATES := $(REPO_ROOT)/test/e2e/data/infrastructure-hivelocity

e2e-cilium-templates:
	helm repo add cilium https://helm.cilium.io/
	helm repo update cilium
	helm template cilium cilium/cilium --version 1.12.2 \
  	--namespace kube-system \
	-f $(REPO_ROOT)/templates/cilium/cilium.yaml > $(REPO_ROOT)/test/e2e/data/cni/cilium/cilium.yaml
	sed -i 's/$${BIN_PATH}/$$BIN_PATH/' $(REPO_ROOT)/test/e2e/data/cni/cilium/cilium.yaml

e2e-ccm-templates:
	helm repo add hivelocity https://hivelocity.github.io/hivelocity-cloud-controller-manager/
	helm repo update hivelocity
	helm template ccm hivelocity/ccm-hivelocity \
		--version 0.1.2 \
	--namespace kube-system > $(REPO_ROOT)/test/e2e/data/ccm/hivelocity-ccm.yaml

.PHONY: cluster-templates
cluster-templates: $(KUSTOMIZE) cluster-templates-v1alpha1 ## Generate cluster templates for all versions

cluster-templates-v1alpha1: $(KUSTOMIZE) ## Generate cluster templates for v1alpha1
	$(KUSTOMIZE) build $(HV_TEMPLATES)/v1alpha1/cluster-template --load-restrictor LoadRestrictionsNone | sed 's/hivelocity_secret_placeholder/$(shell echo -n "$(HIVELOCITY_API_KEY)" | base64 | tr -d '\n')/' > $(HV_TEMPLATES)/v1alpha1/cluster-template.yaml


## --------------------------------------
## Testing
## --------------------------------------

GINKGO_NODES ?= 1
GINKGO_NOCOLOR ?= false
GINKGO_FOCUS ?= ""
GINKGO_SKIP ?= ""
GINKGO_FLAKE_ATTEMPTS ?= 0
ARTIFACTS ?= $(ROOT_DIR)/_artifacts
SKIP_CLEANUP ?= false
SKIP_CREATE_MGMT_CLUSTER ?= false
E2E_DIR ?= $(REPO_ROOT)/test/e2e
E2E_CONF_FILE ?= $(E2E_DIR)/config/hivelocity-ci-envsubst.yaml
KUBETEST_CONF_PATH ?= $(abspath $(E2E_DIR)/data/kubetest/conformance.yaml)

# # to set multiple ginkgo skip flags, if any
# ifneq ($(strip $(GINKGO_SKIP)),)
# _SKIP_ARGS := $(foreach arg,$(strip $(GINKGO_SKIP)),-skip="$(arg)")
# endif

.PHONY: run
run: $(GINKGO)  cluster-templates  ## Run the end-to-end tests
	time $(GINKGO) -v --trace --show-node-events --tags=e2e \
		--timeout=3h \
		--randomize-all -race $(GINKGO_ADDITIONAL_ARGS) \
		--output-dir="$(ARTIFACTS)" \
		--junit-report="junit.e2e_suite.1.xml" \
		--focus=$(GINKGO_FOCUS) --skip=$(GINKGO_SKIP) \
		--nodes=$(GINKGO_NODES) --no-color=$(GINKGO_NOCOLOR) \
		--flake-attempts=$(GINKGO_FLAKE_ATTEMPTS) ./ -- \
		-e2e.artifacts-folder="$(ARTIFACTS)" \
		-e2e.config="$(E2E_CONF_FILE)" \
		-e2e.skip-resource-cleanup=$(SKIP_CLEANUP) \
		-e2e.use-existing-cluster=$(SKIP_CREATE_MGMT_CLUSTER)
