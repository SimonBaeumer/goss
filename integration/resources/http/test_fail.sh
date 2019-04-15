#!/bin/bash
set -euo pipefail


IMAGE="php:7.3-apache"
ID=$(docker run \
        --rm -d \
        -v $(pwd)/"${GOSS_EXE}":/bin/goss \
        -v $(pwd):/app \
        -v $(pwd)/httpd:/var/www/html \
        "${IMAGE}")

function clean {
    printf "\n"
    echo "Stop container..."
    docker stop "${ID}"
}
trap "clean ${ID}" EXIT

sleep 1 # Wait for httpd
docker exec "${ID}" /bin/sh -c 'goss -g /app/goss_fail.yaml validate'