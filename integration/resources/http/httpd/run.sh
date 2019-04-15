#!/bin/bash

IMAGE="goss-httpd-testing"
docker build -t  "${IMAGE}" .
docker run --rm -v $(pwd):/var/www/html -p 8001:80 "${IMAGE}"