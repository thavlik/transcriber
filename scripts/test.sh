#!/bin/bash
cd $(dirname $0)/../cmd
go build -o transcriber
./transcriber test ../test.wav
