#!/bin/bash
set -euo pipefail

IMAGE_NAME="goss-integration-file"

docker build -t "${IMAGE_NAME}" .

docker run --rm -v $(pwd)/"${GOSS_EXE}":/bin/goss -v $(pwd):/app "${IMAGE_NAME}" /bin/sh -c 'goss -g /app/goss.yaml validate'
