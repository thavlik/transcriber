#!/bin/bash
cd $(dirname $0)
source ./env.sh
./build.sh
../cmd/transcriber test transcribe-wav ../test.wav
