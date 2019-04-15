#!/bin/bash
set -euo pipefail

trap "docker-compose down && docker-compose rm" EXIT

docker-compose up -d --build > /dev/null

docker-compose exec -T app goss -g /app/goss.yaml validate
