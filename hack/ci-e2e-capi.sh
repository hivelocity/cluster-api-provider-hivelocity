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

set -o errexit
set -o pipefail

BASE_DIR=$(dirname "${BASH_SOURCE[0]}")
REPO_ROOT=$(realpath "$BASE_DIR/..")

cd "${REPO_ROOT}" || exit 1

# Make sure the tools binaries are on the path.
export PATH="${REPO_ROOT}/hack/tools/bin:${PATH}"
export ARTIFACTS="${ARTIFACTS:-${REPO_ROOT}/_artifacts}"

export HIVELOCITY_SSH_KEY=ssh-key-hivelocity-pub
export SSH_KEY_NAME=$HOME/.ssh/hivelocity
if [ ! -e "$SSH_KEY_NAME" ]; then
    echo "$SSH_KEY_NAME does not exit."
    exit 1
fi

go run ./cmd upload-ssh-pub-key ssh-key-hivelocity-pub "$HOME/.ssh/hivelocity.pub"

# TODO: read count of machines from config, and then claim them, and fail if not enough
#       machines could get claimed.
#       Up to now all machines having these labels will be claimed. The count does
#       not get checked yet.
CONTROL_PLANE_TAG=$(yq .variables.HIVELOCITY_CONTROL_PLANE_MACHINE_TYPE test/e2e/config/hivelocity-ci-envsubst.yaml)
WORKER_TAG=$(yq .variables.HIVELOCITY_WORKER_MACHINE_TYPE test/e2e/config/hivelocity-ci-envsubst.yaml)

go run test/claim-devices-or-fail/claim-devices-or-fail.go cat ${CONTROL_PLANE_TAG} ${WORKER_TAG}

mkdir -p "$ARTIFACTS"
echo "+ run tests!"

if [[ "${CI:-""}" == "true" ]]; then
    make set-manifest-image "MANIFEST_IMG=${IMAGE_PREFIX}/caphv-staging" "MANIFEST_TAG=${TAG}"
    make set-manifest-pull-policy PULL_POLICY=IfNotPresent
fi

make -C test/e2e/ run GINKGO_NODES="${GINKGO_NODES}" GINKGO_FOCUS="${GINKGO_FOKUS}"

