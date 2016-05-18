#!/bin/bash

if [ "$1" == "kubelet" ]; then
    exec /usr/bin/share-mnt /var/lib/kubelet /sys -- kubelet-start.sh "$@"
else
    exec "$@"
fi
