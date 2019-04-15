#!/bin/bash
set -euo pipefail

function cleanup {
    echo "Removing network..."
    docker network rm goss-integration-interface-fail
}
trap cleanup EXIT


docker network create --subnet=172.22.11.0/16 goss-integration-interface-fail

docker run \
    --rm \
    -v $(pwd)/"${GOSS_EXE}":/bin/goss \
    -v $(pwd):/app \
    --net goss-integration-interface-fail --ip 172.22.11.100 \
    centos:7 \
    /bin/sh -c 'goss -g /app/goss_fail.yaml validate'
