#!/bin/bash
cd $(dirname $0)
cd ../cmd
go build -o pdbmesh
time ./pdbmesh test convert DB00207.pdb
