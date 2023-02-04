#!/bin/bash
cd $(dirname $0)
secret() {
    echo $(./get-secret-value.sh bing-cred $1)
}
export BING_API_KEY=$(secret apikey)
cd ../cmd
go build -o imgsearch
./imgsearch test imgsearch ligamentum flavum
