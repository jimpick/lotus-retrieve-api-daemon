#! /bin/bash

set -e

# ovh-0.v6z.me
IP=66.70.191.245
ping -c 1 $IP
export GOFLAGS=-tags=clientretrieve
eval export $(ssh ovh-0 /home/ubuntu/lotus/lotus auth api-info --perm admin | sed "s,0.0.0.0,$IP,")
echo FULLNODE_API_INFO=$FULLNODE_API_INFO
cd ..
go run . daemon
