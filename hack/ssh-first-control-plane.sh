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

export KUBECONFIG=.mgt-cluster-kubeconfig.yaml
if [ ! -s "$KUBECONFIG" ]; then
    echo "$KUBECONFIG does not exist or is empty."
    exit 1
fi
ip=$(kubectl get machine -A -l cluster.x-k8s.io/control-plane  -o  jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}' | head -1)
if [ -z "$ip" ]; then
    echo "Could not get IP of control-plane"
    exit 1
fi

ssh_file=$HOME/.ssh/hivelocity
if [ ! -e $ssh_file ]; then
    echo "$ssh_file does not exist"
    exit 1
fi

ssh -i "$ssh_file" "root@$ip"
