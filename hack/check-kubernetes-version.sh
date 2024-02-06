#!/bin/bash

# Copyright 2022 The Kubernetes Authors.
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

MIN_VERSION=28

if [ -z "$KUBERNETES_VERSION" ]; then
    echo "env var KUBERNETES_VERSION is not set"
    exit 1
fi

# Extract the major version number using regex
if [[ ! $KUBERNETES_VERSION =~ v1\.([0-9]+)\. ]]; then
    echo "The KUBERNETES_VERSION=$KUBERNETES_VERSION format is not recognized."
    exit 1
fi

VERSION_NUMBER=${BASH_REMATCH[1]}

# Compare the extracted number
if [ "$VERSION_NUMBER" -lt $MIN_VERSION ]; then
    echo "KUBERNETES_VERSION=$KUBERNETES_VERSION is less than $MIN_VERSION."
    exit 1
fi
