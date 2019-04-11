#!/bin/bash
set -euo pipefail

docker run -v $(pwd)/${GOSS_EXE}:/bin/goss -v $(pwd):/app centos:7 /bin/sh -c 'goss -g /app/goss.yaml validate'
