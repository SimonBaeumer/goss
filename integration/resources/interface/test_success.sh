#!/bin/bash
set -euo pipefail

function cleanup {
    echo "Removing network..."
    docker network rm goss-integration-test
}
trap cleanup EXIT


docker network create --subnet=172.21.11.0/16 goss-integration-test

docker run \
    --rm \
    -v $(pwd)/"${GOSS_EXE}":/bin/goss \
    -v $(pwd):/app \
    --net goss-integration-test --ip 172.21.11.100 \
    centos:7 \
    /bin/sh -c 'goss -g /app/goss.yaml validate'
