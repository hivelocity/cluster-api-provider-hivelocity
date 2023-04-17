#!/usr/bin/env python3

import re
import sys
import json

keys_to_skip = ['controller', 'controllerGroup', 'controllerKind', 'reconcileID',
                'HivelocityMachine', 'HivelocityCluster', 'Cluster',
                'namespace', 'name', 'Machine']

rows_to_skip = [
    'controller-runtime.webhook'
]

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
    for r in rows_to_skip:
        if r in line:
            return

    if not line.startswith('{'):
        sys.stdout.write(line)
        return
    data = json.loads(line)
    for key in keys_to_skip:
        data.pop(key, None)
    t = data.pop('time', '')
    t = re.sub(r'^.*T(.+)*\..+$', r'\1', t) # '2023-04-17T12:12:53.423Z

    level = data.pop('level', '').ljust(5)
    file = data.pop('file', '')
    message = data.pop('message', '')
    sys.stdout.write(f'{t} {level} "{message}" {file} {data}\n')


if __name__ == '__main__':
    main()
