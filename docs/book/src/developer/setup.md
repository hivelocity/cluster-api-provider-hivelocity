# Setup a Development Environment

## Labeling Machines

First you need to label machines, so that the CAPHV controller can select and provision them.

See [[Provisioning Machines](../topics/provisioning-machines.md) for more about this.

You can use `go run ./test/claim-devices-or-fail` to claim all devices
with a particular label. But be careful, since all machines with this label will get all existing labels removed. Don't accidentally "free" running machines this way.

## SSH-Key

The machines, which get provisioned by CAPHV have a SSH public key installed, so that you can connect
to the machine via ssh.

By default the ssh-key gets stored in the two files `~/.ssh/hivelocity` and `~/.ssh/hivelocity.pub`.

If the key pair does not exist yet, it will get created by the Makefile.

The keys get uploaded with `go run ./cmd upload-ssh-pub-key ssh-key-hivelocity-pub ~/.ssh/hivelocity.pub`. But for most cases this gets automatically handled for you by the scripts.



## hack directory

The directory `hack` contains files to support you. For example the directory `hack/tools/bin` contains
binaries (`envsubst` and `kustomize`) which are needed to create yaml templates.

To make your setup self-contained, these binaries get created if you execute a matching makefile target for the first time.
