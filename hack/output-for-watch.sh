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

kubectl get clusters -A

echo

kubectl get machines -A

echo

kubectl get hivelocitymachine -A

echo

kubectl get events -A --sort-by=metadata.creationTimestamp | tail -8

echo

./hack/tail-caphv-controller-logs.sh


ip=$(kubectl get machine -l cluster.x-k8s.io/control-plane  -o  jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}' | head -1)
if [ -z "$ip" ]; then
    echo "Could not get IP of control-plane"
    exit 1
fi

echo

if netcat -w 1 -z "$ip" 22; then
    echo "$ip ssh port is reachable"
else
    echo "ssh port for $ip is not reachable"
fi

echo

tmpdir="${TMPDIR:-/tmp}"
cluster_name=$(yq .kustomize_substitutions.CLUSTER_NAME tilt-settings.yaml)
kubeconfig="$tmpdir/$cluster_name-workload-cluster-kubeconfig.yaml"
kubectl get secrets "$cluster_name-kubeconfig" -ojsonpath='{.data.value}' | base64 -d > "$kubeconfig"

if [ ! -s "$kubeconfig" ]; then
    echo "failed to get kubeconfig of workload cluster"
    exit
fi

echo "KUBECONFIG=$kubeconfig kubectl cluster-info"
KUBECONFIG=$kubeconfig kubectl cluster-info 2>/dev/null || echo "failed to connect to workload cluster"
