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

if [ ! -s "$KUBECONFIG" ]; then
    ./hack/get-kubeconfig-of-management-cluster.sh
fi

function print_heading() {
    green='\033[0;32m'
    nc='\033[0m' # No Color
    echo -e "${green}${1}${nc}"
}

print_heading Hivelocity

kubectl get clusters -A

print_heading machines

kubectl get machines -A

print_heading hivelocitymachine

kubectl get hivelocitymachine -A

print_heading events

kubectl get events -A --sort-by=metadata.creationTimestamp | tail -8

print_heading logs

./hack/tail-caphv-controller-logs.sh

echo

ip=$(kubectl get machine -A -l cluster.x-k8s.io/control-plane  -o  jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}' | head -1)
if [ -z "$ip" ]; then
    echo "❌ Could not get IP of control-plane"
    exit 1
fi


if netcat -w 2 -z "$ip" 22; then
    echo "👌 $ip ssh port is reachable"
else
    echo "❌ ssh port for $ip is not reachable"
fi

echo

./hack/get-kubeconfig-of-workload-cluster.sh

kubeconfig=".workload-cluster-kubeconfig.yaml"


print_heading "KUBECONFIG=$kubeconfig kubectl cluster-info"
if KUBECONFIG=$kubeconfig kubectl cluster-info >/dev/null 2>&1; then
    echo "👌 cluster is reachable"
else
    echo "❌ cluster is not reachable"
    exit
fi

echo

KUBECONFIG=$kubeconfig kubectl get -n kube-system deployment cilium-operator || echo "❌ cilium-operator not installed?"

KUBECONFIG=$kubeconfig kubectl get -n kube-system deployment ccm-hivelocity || echo "❌ ccm not installed?"

print_heading "workload-cluster nodes"

KUBECONFIG=$kubeconfig kubectl get nodes
