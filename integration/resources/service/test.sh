#!/bin/bash -vx
set -euo pipefail

IMAGE="gossint:centos-service"
GOSS_EXE=${GOSS_EXE:-""}
GOSS_FILE=${2:-"goss.yaml"}

if [[ -z "${GOSS_EXE}" ]]; then
    echo "Error: Provide a path to GOSS_EXE via env variable"
    exit 1
fi

docker build -t "${IMAGE}" .

ID=$(docker run \
    --rm \
    -v $(pwd)/"${GOSS_EXE}":/bin/goss \
    -v $(pwd):/app \
    -v /sys/fs/cgroup:/sys/fs/cgroup:ro \
    --privileged \
    -d \
    "${IMAGE}")
trap "docker stop ${ID}" EXIT

# Start httpd service
docker exec "${ID}" /bin/sh -c "systemctl start httpd "

# Execute goss tests in container
docker exec "${ID}" /bin/sh -c "goss -g /app/${GOSS_FILE} validate"
