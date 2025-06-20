#!/bin/sh

set -e
export PATH=$PATH:/opt/alertmanager/bin

if [ "${1:0:1}" = '-' ]; then
  set -- alertmanager "$@"
fi

exec "$@"