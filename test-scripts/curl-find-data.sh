#! /bin/bash

CID=bafykbzaced3v6jdz436uh2shndde7nwmjlmp6riomr6ps3fbapvaqb6dqpi2o
MINER=t07283

curl -X POST  -H "Content-Type: application/json"  \
  -H "Authorization: Bearer $(cat ~/.lotus/token)"  \
  --data "{ \"jsonrpc\": \"2.0\", \"method\": \"Filecoin.ClientFindData\", \"params\": [{ \"/\": \"$CID\" }, null], \"id\": 1 }"  'http://127.0.0.1:1238/rpc/v0' | jq -C . | less -RM
