#!/bin/sh

set -e

if [ "${1:0:1}" = '-' ]; then
  set -- /opt/prometheus-plug "$@"
fi

exec "$@"