#!/bin/bash
cd $(dirname $0)
secret() {
    echo $(./get-secret-value.sh openai-cred $1)
}
export OPENAI_SECRET_KEY=$(secret secretkey)
cd ../cmd
go build -o define
./define test completion "define ligamentum flavum"
