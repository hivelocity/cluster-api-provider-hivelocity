#!/bin/bash

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

set -e

kubeconfig=".mgt-cluster-kubeconfig.yaml"
if [ -s "$kubeconfig" ]; then
    config=$(yq '.clusters[0].name' .mgt-cluster-kubeconfig.yaml)
    kind_cluster_name="${config#kind-}"
    if [[ $(kind get clusters 2>/dev/null) == *"$kind_cluster_name"* ]]; then
        echo "$kubeconfig already exists"
        exit 0
    fi
    echo "$kubeconfig is outdated. Removing it"
    rm .mgt-cluster-kubeconfig.yaml
fi


cluster_name=$(yq .managementClusterName test/e2e/config/hivelocity-ci-envsubst.yaml)
if [[ $(kind get clusters 2>/dev/null) != *"$cluster_name"* ]]; then
    echo "kind cluster $cluster_name does not exist any more."
    exit 1
fi
kind get kubeconfig --name="$cluster_name" > $kubeconfig
chmod a=,u=rw $kubeconfig
