#!/bin/sh

set -e
export PATH=$PATH:/opt/prometheus

if [ "${1:0:1}" = '-' ]; then
  set -- prometheus "$@"
fi

exec "$@"