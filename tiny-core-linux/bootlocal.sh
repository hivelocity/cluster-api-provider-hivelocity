#!/bin/sh
# Attention: This is /bin/sh running in Tiny Core Linux.
# Bash Syntax does not work here.
set -x
set -e
echo "#############################################################################"
echo "#############################################################################"
echo
echo "Running startup-script.start on $(hostname) at $(date -u)"
echo
echo "#############################################################################"
echo "#############################################################################"

# The first network requests seem to fail, so try it several times
while ! sudo -u tc tce-load -wi curl; do
    date
    echo "Network not up yet. Sleeping and retrying"
    sleep 5
done

authorized_keys=$(grep -o 'authorized_keys="[^"]*"' /proc/cmdline | cut -d= -f2- | tr -d '"')
if [ -n "$authorized_keys" ]; then
    echo "Installing and starting ssh daemon"
    sudo -u tc tce-load -wi openssh
    cp /usr/local/etc/ssh/sshd_config.orig /usr/local/etc/ssh/sshd_config
    /usr/local/etc/init.d/openssh start >/var/log/openssh-start.log 2>&1
    echo "Creating /root/.ssh/authorized_keys from kernel commandline"
    mkdir -p /root/.ssh/

    # Otherwise: bad ownership if sshd tries to read the authorized_keys file.
    chmod 644 /root

    echo "$authorized_keys" >/root/.ssh/authorized_keys
else
    echo 'Not installing and starting ssh daemon, because no authorized_keys="..." kernel parameter was set'
fi

echo "interfaces:"
ifconfig

create_partition() {

    echo "#############################################################################"
    echo "#############################################################################"
    echo "Start of create partition"

    # Create a new partition table and a single primary partition
    # Note: The following commands are piped into fdisk to perform non-interactive disk partitioning
    (
        echo o # Create a new empty DOS partition table
        echo n # Add a new partition
        echo p # Primary partition
        echo 1 # Partition number 1
        echo   # First sector (Accept default: 1)
        echo   # Last sector (Accept default: end of the disk)
        echo w # Write changes
    ) | fdisk "$DEVICE"

    # Wait for the partition table to be re-read

    # reread partition table
    #partprobe
    #mdev -s

    sudo -u tc tce-load -wi parted
    partprobe
    udevadm trigger

    if [ ! -e "$PART" ]; then
        echo "Could not find partition $PART. Exit"
        exit 0
    fi

    yes | mkfs.ext4 "$PART"
    echo "#############################################################################"
    echo "#############################################################################"
}

create_public_interface() {
    echo "#############################################################################"
    echo "Setting up public IP"
    ip=$(jq -r '.interfaces.public.ip' /metadata.json)
    netmask=$(jq -r '.interfaces.public.netmask' /metadata.json)
    gateway=$(jq -r '.interfaces.public.gateway' /metadata.json)
    ifconfig eth0 "$ip" netmask "$netmask"
    route add default gw "$gateway"
    echo
    route -n
    echo "#############################################################################"
}

DEVICE=$(grep -o 'targetdrive=[^ ]*' /proc/cmdline | cut -d= -f2-)
if [ -z "$DEVICE" ]; then
    echo "Could not find targetdrive=... in /proc/cmdline."
    echo "Please specify the targetdrive while booting."
    echo "Example: kernel ... targetdrive=/dev/sda"
    exit 1
fi

PART="${DEVICE}1"

URL=$(grep -o 'sourceurl=[^ ]*' /proc/cmdline | cut -d= -f2-)
if [ -z "$URL" ]; then
    echo "Could not find sourceurl=... in /proc/cmdline"
    echo "Please specify the sourceurl while booting."
    echo "Example: kernel ... sourceurl=https://user:password@example.com/machine-image.tgz"
    exit 1
fi

metadata_url=$(grep -o 'metadata_url=[^ ]*' /proc/cmdline | cut -d= -f2-)
if [ -z "$metadata_url" ]; then
    echo "Could not find metadata_url=... in /proc/cmdline"
    echo "Please specify the metadata_url while booting."
    echo "Example: kernel ... metadata_url={metadata_url}"
    echo "See https://developers.hivelocity.net/docs/custom-ipxe"
    exit 1
fi

sudo -u tc tce-load -wi jq
sudo -u tc tce-load -wi libzstd

if [ "${metadata_url#http}" = "$metadata_url" ]; then
    echo "Metadata URL does not start with http"
    echo "metadata_url: $metadata_url"
    metadata_url=
else
    echo
    echo "#############################################################################"
    echo "metadata_url: $metadata_url"
    echo "$metadata_url" >/metadata.url
    curl -sSL --fail -o /metadata.json "$metadata_url"
    cat /metadata.json
    echo "#############################################################################"
    echo
    create_public_interface
fi

if [ ! -e "$DEVICE" ]; then
    echo "$DEVICE not found. Exit"
    echo "existing devices:"
    exit 1
fi

create_partition

mount -t ext4 "$PART" /mnt

echo "#############################################################################"
echo "#############################################################################"
echo
echo "If you want to scroll up, then you can freeze the terminal with ctrl-s."
echo "Then you can scroll up with shift-PageUp"
echo "Unfreeze the terminal with ctrl-q"
echo

curl -SL --fail --retry 20 -o- "$URL" | tar -C /mnt -xzf-

for i in dev proc sys dev/pts; do
    mkdir -p /mnt/$i
    mount -o bind /$i /mnt/$i
done
cp /proc/mounts /mnt/etc/mtab

if [ -e /metadata.json ]; then
    cp /metadata.json /mnt/
fi

# These files were missing in a image, and systemd failed to start.
mkdir -p /mnt/run
mkdir -p /mnt/run/lock

chroot /mnt /bin/bash <<EOF
grub-install $DEVICE
update-grub
EOF

echo "Installed the image to $PART."

finish_url=$(jq -r '.finishHook.url' /metadata.json) || true
if [ -z "$finish_url" ]; then
    echo "failed to find finish_url to finish provisioning in metadata.json"
    echo "BTW, the metadata_url times out after 30 minutes."
    echo "See https://developers.hivelocity.net/docs/custom-ipxe"
    cat /metadata.json
    exit 1
fi
# Tell Hivelocity API, that machine is provisioned
# We hope to get out of the stuck endless reloading state like this.
# This switches of iPXE for this device. The next boot should be from the disc.
curl -XPOST -sSL --fail "$finish_url"

echo "Looks good!"
echo "Restarting"

# Wait some seconds, so error message can be read.
sleep 5

# Leave Tinycore Linux, switch to downloaded machine image.
reboot
