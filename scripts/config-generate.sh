#!/bin/sh

set -e

WORK_DIR=$(cd -P "$(dirname "$0")"/ > /dev/null; pwd)
WORK_CONFIG=$(dirname ${WORK_DIR})/conf

if [ ! -f ${WORK_CONFIG}/env ];then
    echo "${WORK_CONFIG}/env file does not exist" && exit 0
fi

export $(grep -v "^#" ${WORK_CONFIG}/env | xargs)

config_dev()
{
    envsubst < ${WORK_CONFIG}/config.env.yml > ${WORK_CONFIG}/config.yml
    return 0
}

config_docker()
{
    export SECRET_KEY=
    export MYSQL_HOST=127.0.0.1
    export MYSQL_USER=root
    export MYSQL_PASSWORD=
    export MYSQL_PORT=3306
    export MYSQL_DATABASE=yonecloud
    export MYSQL_PREFIX=yone

    envsubst < ${WORK_CONFIG}/config.env.yml > ${WORK_CONFIG}/config.docker.yml
    return 0
}

config_docker_rm()
{
    [ -f ${WORK_CONFIG}/config.docker.yml ] && rm -f ${WORK_CONFIG}/config.docker.yml
    return 0
}

usage()
{
	cat <<EOF
Usage: $0 <option>
Option: dev | docker | docker-rm

Example:
  /bin/sh $0 dev
  /bin/sh $0 docker
  /bin/sh $0 docker-rm
EOF
    return 0
}

if [ "$#" -eq 1 ];then
    case "$1" in
        dev)
            config_dev
            ;;
        docker)
            config_docker
            ;;
        docker-rm)
            config_docker_rm
            ;;
        all)
            config_dev
            config_docker
            ;;
        *)
            usage
            ;;
    esac
else
    usage
fi