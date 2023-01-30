#!/bin/bash
cd $(dirname $0)
source ./env.sh
cd ../cmd
go build -o transcriber
