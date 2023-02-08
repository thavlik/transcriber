#!/bin/bash
cd $(dirname $0)
set -euo pipefail
source ./env.sh
./build.sh
../cmd/transcriber test transcribe /mnt/c/Users/13169/Desktop/judy.mp3
