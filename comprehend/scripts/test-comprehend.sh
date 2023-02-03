#!/bin/bash
cd $(dirname $0)
source ./env.sh
./build.sh
export TEXT="The ligamenta flavum is a short but thick ligament that connects the laminae of adjacent vertebrae"
../cmd/comprehend test comprehend
