#!/bin/bash
set -euo pipefail
# To test on different operating systems it is necessary to add more logic.
# OS parameter defines the operating system
# GOSS_FILE defines which goss file in this directory should be executed
#
# Example:
# ./test.sh centos goss.yaml


OS=${1:-"centos"}
GOSS_EXE=${GOSS_EXE:-""}
GOSS_FILE=${2:-"goss.yaml"}

if [[ -z "${GOSS_EXE}" ]]; then
    echo "Error: Provide a path to GOSS_EXE via env variable"
    exit 1
fi

IMAGE=""
case "${OS}" in
"ubuntu"*)
    IMAGE="ubuntu:18.04"
    ;;
"alpine"*)
    IMAGE="alpine:3.9"
    ;;
"centos"*)
    IMAGE="centos:7"
    ;;
*)
    echo "No valid OS was given - available are ubuntu, alpine and centos"
    exit 1
esac

echo "IMAGE: ${IMAGE}"
docker run --rm -v $(pwd)/"${GOSS_EXE}":/bin/goss -v $(pwd)/"${OS}":/app "${IMAGE}" /bin/sh -c "goss -g /app/${GOSS_FILE} validate"
