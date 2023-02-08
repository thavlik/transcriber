#!/bin/bash
cd $(dirname $0)
set -euo pipefail
source ./env.sh
./build.sh
../cmd/transcriber batch-transcribe
