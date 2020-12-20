#! /bin/bash

set -e

ping -c 1 10.0.1.52
export GOFLAGS=-tags=clientretrieve
eval export $(ssh lotus2 /home/lotus2/lotus/lotus auth api-info --perm admin | sed 's,0.0.0.0,10.0.1.52,')
echo FULLNODE_API_INFO=$FULLNODE_API_INFO
cd ..
go run . daemon
