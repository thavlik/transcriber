#!/bin/bash
kubectl get secret -n ts $1 -o json \
    | jq .data.$2 \
    | xargs echo \
    | base64 -d
