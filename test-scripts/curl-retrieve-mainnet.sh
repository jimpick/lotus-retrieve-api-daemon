#! /bin/bash

# lotus client retrieve --miner=f0105208 bafykbzacecp5ijz6n2hoy6nwferyxajnnftfny52kntbmqv5ywo734yivekw4 $PWD/retrieval.bin-$(date +%s)

CLIENT=f3vp7m3244tjtxrvg4n2lfedtqnnnzhyno3ym6vnl4wzozztik4f2kvzfbfbgzcga7g3mckddw6x4ahp5n4iwa
MINER=f0105208

# 4KiB
#CID=bafk2bzacebhlhbcnhmvover42qq5bx773c522skieho6nhtbz7d2ow3f4sw24

# 4MiB
CID=bafykbzacecp5ijz6n2hoy6nwferyxajnnftfny52kntbmqv5ywo734yivekw4

# 105MiB
#CID=bafykbzacecdt3sxmqj2xi6ntatf4ikycu6korkyv3pd632zot64lah5oonftw

mkdir -p downloads
DEST=$PWD/downloads/$CID.bin-$(date +'%s')
echo "Downloading to: $DEST"

OFFER="$(curl -X POST  -H "Content-Type: application/json" --data "{ \"jsonrpc\": \"2.0\", \"method\": \"Filecoin.ClientMinerQueryOffer\", \"params\": [\"$MINER\", { \"/\": \"$CID\" }, null], \"id\": 1 }"  'http://127.0.0.1:1238/rpc/v0')"
#echo $OFFER | jq
ROOT=$(echo $OFFER | jq .result.Root)
SIZE=$(echo $OFFER | jq .result.Size)
MIN_PRICE=$(echo $OFFER | jq .result.MinPrice)
UNSEAL_PRICE=$(echo $OFFER | jq .result.UnsealPrice)
PAYMENT_INTERVAL=$(echo $OFFER | jq .result.PaymentInterval)
PAYMENT_INTERVAL_INCREASE=$(echo $OFFER | jq .result.PaymentIntervalIncrease)
MINER_OWNER=$(echo $OFFER | jq .result.Miner)
MINER_PEER=$(echo $OFFER | jq .result.MinerPeer)

ORDER="{ \"Root\": $ROOT, \"Piece\": null, \"Size\": $SIZE, \"Total\": $MIN_PRICE, \"UnsealPrice\": $UNSEAL_PRICE, \"PaymentInterval\": $PAYMENT_INTERVAL, \"PaymentIntervalIncrease\": $PAYMENT_INTERVAL_INCREASE, \"Client\": \"$CLIENT\", \"Miner\": $MINER_OWNER, \"MinerPeer\": $MINER_PEER }"
#echo $ORDER | jq

FILEREF="{ \"Path\": \"$DEST\", \"IsCAR\": false }"
DATA="{ \"jsonrpc\": \"2.0\", \"method\": \"Filecoin.ClientRetrieve\", \"params\": [ $ORDER, $FILEREF ], \"id\": 1 }"

curl -X POST -H "Content-Type: application/json" \
       	--data "$DATA" \
       	'http://127.0.0.1:1238/rpc/v0' 
