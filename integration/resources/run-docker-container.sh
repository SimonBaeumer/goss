#!/usr/bin/env bash

GOSS_EXE=${GOSS_EXE:-"../../release/goss-linux-amd64"}

docker run --rm -it -v $(pwd)/"${GOSS_EXE}":/bin/goss -v $(pwd):/app centos:7 /bin/bash