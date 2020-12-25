#! /bin/bash

eval export $(ssh lotus2 /home/lotus2/lotus/lotus auth api-info --perm admin)
echo FULLNODE_API_INFO=$FULLNODE_API_INFO
TOKEN=$(echo $FULLNODE_API_INFO | sed 's,:.*$,,')

CLIENT=t3qkztmkfopk63qsel2xk3ek4w22epn3jnnlubnwjha2sl7rjhiuduwx24xivmhtdz7st3zmteuemeefply55q
OWNER=t07281
AMOUNT=16777216

DATA="{ \"jsonrpc\": \"2.0\", \"method\": \"Filecoin.PaychGet\", \"params\": [ \"$CLIENT\", \"$OWNER\", \"$AMOUNT\" ], \"id\": 1 }"

curl -X POST -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN"  \
  --data "$DATA" \
  'http://10.0.1.52:7234/rpc/v0' | jq -C .

