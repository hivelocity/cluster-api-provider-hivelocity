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

set -o errexit
set -o pipefail


create_ssh_key() {
    echo "generating new ssh key"
    ssh-keygen -t ed25519 -f ${SSH_KEY_PATH} -N '' 2>/dev/null <<< y >/dev/null
}

remove_ssh_key() {
    local ssh_fingerprint=$1
    echo "removing ssh key"
    rm -f ${SSH_KEY_PATH}

    ${REPO_ROOT}/hack/log/redact.sh || true
}
