#!/bin/bash
cd "$(dirname "$0")"
export NODE_PATH=/usr/local/lib/node_modules
export INPUT_URL="https://go.drugbank.com/drugs/DB00571" #"https://go.drugbank.com/drugs/DB00186"
node ./query-drug.js
