#!/bin/bash
cd $(dirname $0)
secret() {
    echo $(./get-secret-value.sh openai-cred $1)
}
export OPENAI_SECRET_KEY=$(secret secretkey)
cd ../cmd
go build -o define
test_disease() {
    out=$(./define test is-disease $@)
    echo "$@: $out"
}
test_disease ligamentum flavum
test_disease diabetes
test_disease arthritis
test_disease cancer
test_disease skull
test_disease blood vessels
test_disease stroke