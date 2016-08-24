#!/bin/bash

mount --rbind /host/dev /dev

FQDN=$(hostname --fqdn || hostname)

exec "$@" --hostname-override ${FQDN}
