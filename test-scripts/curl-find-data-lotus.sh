#! /bin/bash

eval export $(ssh lotus2 /home/lotus2/lotus/lotus auth api-info --perm admin)
echo FULLNODE_API_INFO=$FULLNODE_API_INFO
TOKEN=$(echo $FULLNODE_API_INFO | sed 's,:.*$,,')
CID=bafykbzaced3v6jdz436uh2shndde7nwmjlmp6riomr6ps3fbapvaqb6dqpi2o
MINER=t07283

curl -X POST  -H "Content-Type: application/json"  \
  -H "Authorization: Bearer $TOKEN"  \
  --data "{ \"jsonrpc\": \"2.0\", \"method\": \"Filecoin.ClientFindData\", \"params\": [{ \"/\": \"$CID\" }, null], \"id\": 1 }"  'http://10.0.1.52:7234/rpc/v0' | jq -C . | less -RM
