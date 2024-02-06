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

# shellcheck source=../hack/ci-e2e-sshkeys.sh
source "${REPO_ROOT}/hack/ci-e2e-sshkeys.sh"

# We need to export the SSH_KEY_NAME as a environment variable
SSH_KEY_NAME=caphv-e2e-$(
    head /dev/urandom | tr -dc A-Za-z0-9 | head -c 12
    echo ''
)
echo "SSH Key Name : $SSH_KEY_NAME"
export SSH_KEY_PATH=/tmp/${SSH_KEY_NAME}
echo "SSH Key Path : $SSH_KEY_PATH"
export HIVELOCITY_SSH_KEY=${SSH_KEY_NAME}
create_ssh_key ${SSH_KEY_PATH}
trap 'remove_ssh_key ${SSH_KEY_NAME}' EXIT

if [ ! -e "$SSH_KEY_PATH" ]; then
    echo "$SSH_KEY_PATH does not exit."
    exit 1
fi

go run ./cmd upload-ssh-pub-key $HIVELOCITY_SSH_KEY "$SSH_KEY_PATH.pub"

CONTROL_PLANE_TAG=$(yq .variables.HIVELOCITY_CONTROL_PLANE_DEVICE_TYPE "$E2E_CONF_FILE")
WORKER_TAG=$(yq .variables.HIVELOCITY_WORKER_DEVICE_TYPE "$E2E_CONF_FILE")

echo "#####################################################################################################"
echo "All devices with one of these device tags will get claimed. They will get provisioned soon:"
echo "$CONTROL_PLANE_TAG $WORKER_TAG"
echo "#####################################################################################################"

go run test/claim-devices-or-fail/claim-devices-or-fail.go ${CONTROL_PLANE_TAG} ${WORKER_TAG}

mkdir -p "$ARTIFACTS"
echo "+ run tests!"

if [[ "${CI:-""}" == "true" ]]; then
    make set-manifest-image "MANIFEST_IMG=ghcr.io/hivelocity/caphv-staging" "MANIFEST_TAG=${TAG}"
    make set-manifest-pull-policy PULL_POLICY=IfNotPresent
fi

make -C test/e2e/ run GINKGO_NODES="${GINKGO_NODES}" GINKGO_FOCUS="${GINKGO_FOKUS}"
