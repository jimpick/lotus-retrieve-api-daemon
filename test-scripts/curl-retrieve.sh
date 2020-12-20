#! /bin/bash

# lotus client retrieve --miner t07283 bafykbzaced3v6jdz436uh2shndde7nwmjlmp6riomr6ps3fbapvaqb6dqpi2o $PWD/wiki.zip.aa.aa-$(date +%s)

CLIENT=t3qkztmkfopk63qsel2xk3ek4w22epn3jnnlubnwjha2sl7rjhiuduwx24xivmhtdz7st3zmteuemeefply55q
MINER=t07283
CID=bafykbzaced3v6jdz436uh2shndde7nwmjlmp6riomr6ps3fbapvaqb6dqpi2o

mkdir -p downloads
DEST=$PWD/downloads/wiki.zip.aa.aa-$(date +'%s')
echo "Downloading to: $DEST"

OFFER="$(curl -X POST  -H "Content-Type: application/json" --data "{ \"jsonrpc\": \"2.0\", \"method\": \"Filecoin.ClientMinerQueryOffer\", \"params\": [\"$MINER\", { \"/\": \"$CID\" }, null], \"id\": 1 }"  'http://127.0.0.1:1238/rpc/v0')"
echo $OFFER | jq

ROOT=$(echo $OFFER | jq .result.Root)
SIZE=$(echo $OFFER | jq .result.Size)
MIN_PRICE=$(echo $OFFER | jq .result.MinPrice)
UNSEAL_PRICE=$(echo $OFFER | jq .result.UnsealPrice)
PAYMENT_INTERVAL=$(echo $OFFER | jq .result.PaymentInterval)
PAYMENT_INTERVAL_INCREASE=$(echo $OFFER | jq .result.PaymentIntervalIncrease)
MINER_OWNER=$(echo $OFFER | jq .result.Miner)
MINER_PEER=$(echo $OFFER | jq .result.MinerPeer)

ORDER="{ \"Root\": $ROOT, \"Piece\": null, \"Size\": $SIZE, \"Total\": $MIN_PRICE, \"UnsealPrice\": $UNSEAL_PRICE, \"PaymentInterval\": $PAYMENT_INTERVAL, \"PaymentIntervalIncrease\": $PAYMENT_INTERVAL_INCREASE, \"Client\": \"$CLIENT\", \"Miner\": $MINER_OWNER, \"MinerPeer\": $MINER_PEER }"
echo $ORDER | jq

FILEREF="{ \"Path\": \"$DEST\", \"IsCAR\": false }"
echo $FILEREF | jq

DATA="{ \"jsonrpc\": \"2.0\", \"method\": \"Filecoin.ClientRetrieve\", \"params\": [ $ORDER, $FILEREF ], \"id\": 1 }"

curl -X POST -H "Content-Type: application/json" \
       	--data "$DATA" \
       	'http://127.0.0.1:1238/rpc/v0' 
ls -lh $DEST
