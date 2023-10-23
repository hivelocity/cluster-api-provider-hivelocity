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

set -euo pipefail
ID=${1:-}
if [ -z "$ID" ]; then
     echo "Please provide device-id"
     exit 1
fi
outdir=machine-state-$ID-$(date +"%Y-%m-%d--%H-%M-%S")
mkdir $outdir
echo "writing to $outdir"

CURL="curl -sSL --fail-with-body"

echo "Get Power"
$CURL \
     --url https://core.hivelocity.net/api/v2/device/$ID/power \
     --header "X-API-KEY: $HIVELOCITY_API_KEY" \
     --header 'accept: application/json'  | yq -P > $outdir/power.yaml

echo "Device"
$CURL \
     --url https://core.hivelocity.net/api/v2/device/$ID \
     --header "X-API-KEY: $HIVELOCITY_API_KEY" \
     --header 'accept: application/json'  | yq -P > $outdir/device.yaml

echo "Bare Meta Device"
$CURL \
     --url https://core.hivelocity.net/api/v2/bare-metal-devices/$ID \
     --header "X-API-KEY: $HIVELOCITY_API_KEY" \
     --header 'accept: application/json' | yq -P > $outdir/bare-metal-device.yaml

echo "IPMI Data"
$CURL \
     --url https://core.hivelocity.net/api/v2/device/$ID/ipmi \
     --header "X-API-KEY: $HIVELOCITY_API_KEY" \
     --header 'accept: application/json' | yq -P > $outdir/ipmi.yaml


echo "Events"
$CURL \
     --url https://core.hivelocity.net/api/v2/device/$ID/events \
     --header "X-API-KEY: $HIVELOCITY_API_KEY" \
     --header 'accept: application/json' | yq -P > $outdir/events.yaml

echo "done. See $outdir"