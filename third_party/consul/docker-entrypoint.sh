#!/bin/sh

set -e
export PATH=$PATH:/opt/consul/bin

if [ -z "$CONSUL_BIND" ]; then
  CONSUL_BIND="0.0.0.0"
fi

if [ -z "$CONSUL_CLIENT" ]; then
  CONSUL_CLIENT="0.0.0.0"
fi

if [ -z "$CONSUL_DATA" ]; then
  CONSUL_DATA_DIR=/opt/consul/data
fi

if [ -z "$CONSUL_CONFIG" ]; then
  CONSUL_CONFIG_DIR=/opt/consul/config
fi

if [ "${1:0:1}" = '-' ]; then
  set -- consul "$@"
fi

if [ "$1" = 'agent' ]; then
  shift
  set -- consul agent \
    -data-dir="$CONSUL_DATA_DIR" \
    -config-dir="$CONSUL_CONFIG_DIR" \
    -bind=$CONSUL_BIND \
    -client=$CONSUL_CLIENT \
    "$@"
else
  set -- consul "$@"
fi

exec "$@"