#!/bin/bash
set -euo pipefail

ID=$(docker run -d -v $(pwd)/"${GOSS_EXE}":/bin/goss -v $(pwd):/app httpd:2.4)

function clean {
    printf "\n"
    echo "Stop container..."
    docker stop "${ID}"
    echo "Remove container..."
    docker rm "${ID}"
}
trap "clean ${ID}" EXIT

docker exec "${ID}" /bin/sh -c 'goss -g /app/goss.yaml validate'