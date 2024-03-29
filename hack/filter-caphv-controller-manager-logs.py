#!/usr/bin/env python3

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

import re
import sys
import json

keys_to_skip = ['controller', 'controllerGroup', 'controllerKind', 'reconcileID',
                'HivelocityCluster', 'Cluster',
                'namespace', 'name', 'Machine', 'stacktrace', 'stack', 'logger']

rows_to_skip = [
    'Bootstrap not ready - requeuing',
    'controller-runtime.webhook', 'certwatcher/certwatcher', 'Registering a validating webhook',
    'Registering a mutating webhook', 'Starting EventSource',
    '"Reconciling finished"',
    '"Reconciling HivelocityMachine"',
    '"Reconcile HivelocityMachine"',
    '"Reconciling HivelocityCluster"',
    '"Creating cluster scope"',
    '"Starting reconciling cluster"',
    '"Completed function"',
    '"Adding request."',
    '"validate update" v1alpha1/hivelocitymachine_webhook',
    '"validate update" v1alpha1/hivelocitycluster_webhook',
    'spec.status was updated',
    'DEBUG "ReconcileState" ',
]

rows_to_skip_regex = [
    r'''DEBUG "Started function" device/device.go:\d+ {'HivelocityMachine': {'name': '\S+', 'namespace': '\S+'}, 'function': 'actionDeviceProvisioned'}''',
    r'''DEBUG "wrote response" admission.*''',
    r'''DEBUG "received request" admission.*''',
    r'''DEBUG.*hivelocity API called.*GET.*''',
]

def printHivelocityMachine(value):
    return 'HivelocityMachine=' + value.get('name', '')

key_to_output = {
    'HivelocityMachine': printHivelocityMachine,
}
def main():

    if len(sys.argv) == 1 or sys.argv[1] in ['-h', '--help']:
        print('''%s [file|-]
    filter the logs of caphv-controller-manager.
    Used for debugging.
    ''' % sys.argv[0])
        sys.exit(1)

    if sys.argv[1] == '-':
        fd = sys.stdin
    else:
        fd = open(sys.argv[1])
    read_logs(fd)


def read_logs(fd):
    for line in fd:
        handle_line(line)

def handle_line(line):
    if not line.startswith('{'):
        sys.stdout.write(line)
        return
    data = json.loads(line)
    for key in keys_to_skip:
        data.pop(key, None)

    formatted_data = []
    for key, func in key_to_output.items():
        if key not in data:
            continue
        value = data.pop(key)
        if not key:
            continue
        value = func(value)
        if not value:
            continue
        formatted_data.append(value)

    t = data.pop('time', '')
    t = re.sub(r'^.*T(.+)*\..+$', r'\1', t) # '2023-04-17T12:12:53.423Z

    level = data.pop('level', '').ljust(5)
    file = data.pop('file', '')
    message = data.pop('message', '')
    if not data:
        data=''

    formatted_data = ' '.join(formatted_data)
    new_line = f'{t} {level} "{message}" {file} {formatted_data} {data}\n'
    for r in rows_to_skip:
        if r in new_line:
            return
    for r in rows_to_skip_regex:
        if re.search(r, new_line):
            return

    sys.stdout.write(new_line)


if __name__ == '__main__':
    main()
