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

kubectl get events -A -o=wide --sort-by=.lastTimestamp | grep -vP 'LeaderElection|CSRApproved' | tail -8

print_heading conditions

go run github.com/guettli/check-conditions@latest all | grep -vP 'ScalingUp|WaitingForInfrastructure|WaitingForNodeRef|WaitingForAvailableMachines|^Checked.*Duration:'


print_heading logs

./hack/tail-controller-logs.sh

echo

capi_error="$(kubectl logs -n capi-system --since=5m deployments/capi-controller-manager | \
    grep -iP 'error|\be\d\d\d\d\b' | \
    grep -vP 'ignoring DaemonSet-managed Pods|TLS handshake error from' | \
    tail -7)"
if [ -n "$capi_error" ]; then
  print_heading capi controller errors
  echo "$capi_error"
fi

ip=$(kubectl get cluster -A -o=jsonpath='{.items[*].spec.controlPlaneEndpoint.host}' | head -1)
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

kubeconfig_wl=".workload-cluster-kubeconfig.yaml"

print_heading "KUBECONFIG=$kubeconfig_wl kubectl cluster-info"
if KUBECONFIG=$kubeconfig_wl kubectl --request-timeout=2s cluster-info >/dev/null 2>&1; then
    echo "👌 cluster is reachable"
else
    echo "❌ cluster is not reachable"
    exit
fi

echo

KUBECONFIG=$kubeconfig_wl kubectl get -n kube-system deployment cilium-operator || echo "❌ cilium-operator not installed? To install CNI and CCM in wl-cluster: make install-essentials"

KUBECONFIG=$kubeconfig_wl kubectl get -n kube-system deployment ccm-hivelocity || echo "❌ ccm not installed?"

print_heading "workload-cluster nodes"

KUBECONFIG=$kubeconfig_wl kubectl get nodes

if [ "$(kubectl get hivelocitymachine | wc -l)" -ne "$(KUBECONFIG="$kubeconfig_wl" kubectl get nodes | wc -l)" ]; then
    echo "❌ Number of nodes in wl-cluster does not match number of machines in mgt-cluster"
else
    echo "👌 number of nodes in wl-cluster is equal to number of machines in mgt-cluster"
fi

not_approved=$(KUBECONFIG=$kubeconfig_wl kubectl get csr --no-headers  --sort-by='.metadata.creationTimestamp' | grep -v Approved | tail -8 )
if [ -n "$not_approved" ]; then
    echo "❌ (CSRs)certificate signing requests which are not approved"
    echo "$not_approved"
fi
