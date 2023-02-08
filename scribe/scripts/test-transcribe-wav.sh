#!/bin/bash
cd $(dirname $0)
set -euo pipefail
source ./env.sh
./build.sh
../cmd/transcriber test transcribe ../test.wav
