#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
secret() {
    echo $(./get-secret-value.sh mongo-cred $1)
}
export MONGO_USERNAME=$(secret username)
export MONGO_PASSWORD=$(secret password)
export MONGO_HOST=$(secret host)
export MONGO_DATABASE=$(secret database)
export MONGO_URL="mongodb+srv://${MONGO_USERNAME}:${MONGO_PASSWORD}@${MONGO_HOST}/admin?tls=true&authSource=admin&replicaSet=db-mongodb-nyc1-29482"
echo "connecting to $MONGO_HOST"
docker run -it \
    mongo:5 mongosh ${MONGO_URL}

