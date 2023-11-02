 # Setup a Development Environment

## Tagging Machines

First you need to tag machines, so that the CAPHV controller can select and provision them. The devices must have `caphv-use=allow` tag so that the controller can use them.

See [Provisioning Machines](../topics/provisioning-machines.md) for more about this.

`make tilt-up` and other Makefile targets call `go run ./test/claim-devices-or-fail` to claim all devices
with a particular tag. But be careful, since all machines with this tag will get all existing tags removed.
Don't accidentally "free" running machines this way.

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

## API Key

Create a HIVELOCITY_API_KEY and add it to your `.envrc` file.

You can create this key via the web UI: [my.hivelocity.net/account](https://my.hivelocity.net/account).

## Kind

During development the management-cluster runs in local [kind cluster](https://kind.sigs.k8s.io/).

The tool `kind` gets installed into `hack/tools/bin/` automatically via the Makefile.
