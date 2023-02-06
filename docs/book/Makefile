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

# Directories.
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
TOOLS_DIR := $(realpath ../../hack/tools)
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
BIN_DIR := bin
MDBOOK_INSTALL := $(realpath ../../hack/ci-install-mdbook.sh)
# API_DIRS to skip.
API_DIRS :=
API_SRCS := $(foreach dir, $(API_DIRS), $(call rwildcard,../../$(dir),*.go))
GENERATED_SRCS := $(foreach dir, $(API_DIRS),../../$(dir)/zz_generated.deepcopy.go ../../$(dir)/zz_generated.conversion.go ../../$(dir)/zz_generated.defaults.go)
API_FILTERED_SRCS := $(filter-out $(GENERATED_SRCS), $(API_SRCS))
GENCRDAPIREFERENCEDOCS_SRCS := $(call rwildcard,.,gen-crd-docs/*.*)

export PATH := $(TOOLS_BIN_DIR):$(PATH)

# Binaries.
MDBOOK_TABULATE := $(TOOLS_BIN_DIR)/mdbook-tabulate
$(MDBOOK_TABULATE): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/mdbook-tabulate  sigs.k8s.io/cluster-api/hack/tools/mdbook/tabulate

MDBOOK_EMBED := $(TOOLS_BIN_DIR)/mdbook-embed
$(MDBOOK_EMBED): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/mdbook-embed sigs.k8s.io/cluster-api/hack/tools/mdbook/embed

MDBOOK_RELEASELINK := $(TOOLS_BIN_DIR)/mdbook-releaselink
$(MDBOOK_RELEASELINK): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/mdbook-releaselink sigs.k8s.io/cluster-api/hack/tools/mdbook/releaselink

# Only set MDBOOK if it's not set as an environment variable
MDBOOK ?= $(TOOLS_BIN_DIR)/mdbook
$(MDBOOK):
	$(MDBOOK_INSTALL) v0.4.25 $(TOOLS_BIN_DIR)

EMOJITSU := $(TOOLS_BIN_DIR)/emojitsu
$(EMOJITSU):
	wget -O $(TOOLS_BIN_DIR)/emojitsu --no-verbose https://github.com/shonfeder/emojitsu/releases/download/0.1.1/gh-actions-emojitsu
	chmod +x $(TOOLS_BIN_DIR)/emojitsu

GENCRDAPIREFERENCEDOCS := $(TOOLS_BIN_DIR)/gen-crd-docs
$(GENCRDAPIREFERENCEDOCS): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/gen-crd-docs github.com/ahmetb/gen-crd-api-reference-docs

src/crd:
	mkdir -p src/crd

src/crd/index.md: $(API_FILTERED_SRCS) src/crd
	$(MAKE) gen_crd_docs API_DIR="$(shell go list -m)" OUT_FILE=$@

.PHONY: gen_crd_docs
gen_crd_docs: $(GENCRDAPIREFERENCEDOCS) src/crd $(GENCRDAPIREFERENCEDOCS_SRCS)
	$(GENCRDAPIREFERENCEDOCS) -template-dir gen-crd-docs/template -config "gen-crd-docs/config.json" \
		-api-dir $(API_DIR) \
		-out-file $(OUT_FILE)

BOOK_DEPS := $(MDBOOK) $(MDBOOK_EMBED) $(MDBOOK_RELEASELINK) $(MDBOOK_TABULATE) $(EMOJITSU) src/crd/index.md

.PHONY: serve
serve: $(BOOK_DEPS)
	$(MDBOOK) serve

.PHONY: build
build: $(BOOK_DEPS)
	$(MDBOOK) build
