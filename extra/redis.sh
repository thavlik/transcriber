#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
secret() {
    echo $(./get-secret-value.sh redis-cred $1)
}

# This image is just redis:7 with ca-certificates installed.
# It is needed because DigitalOcean redis uses SSL.
IMAGE=thavlik/redis:7

export REDIS_URI=$(secret uri)
docker run -it $IMAGE \
    redis-cli -u $REDIS_URI
    