#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
goout=../api
rm $goout/client.gen.go || true
rm $goout/server.gen.go || true
mkdir -p $goout
name=iam-defs

# Build the definitions image
docker build -t $name .

# Check if there's a container that is already running.
# This can happen if this script exits with error, and
# the only penalty is having to wait longer until these
# containers exit.
if [[ -n $(docker ps | grep $name) ]]; then
    echo "Stopping currently running container \"$name\""
    set +e
    docker stop $name
    docker container wait $name
    docker container rm $name --force
    set -e
fi

docker run -d --name $name --rm $name tail -f /dev/null &>/dev/null
docker exec $name cat server.gen.go > $goout/server.gen.go
docker exec $name cat client.gen.go > $goout/client.gen.go
docker stop $name &>/dev/null

echo "Successfully generated definitions"
