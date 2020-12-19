module github.com/jimpick/lotus-retrieve-api-daemon

go 1.15

require (
	github.com/filecoin-project/go-address v0.0.5-0.20201103152444-f2023ef3f5bb
	github.com/filecoin-project/go-fil-markets v1.0.10
	github.com/filecoin-project/go-jsonrpc v0.1.2
	github.com/filecoin-project/lotus v1.2.1
	github.com/filecoin-project/specs-actors v0.9.13 // indirect
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-graphsync v0.5.1
	github.com/jimpick/lotus-query-ask-api-daemon v0.0.0-20201218023526-5265063b5560
	github.com/jimpick/lotus-utils v0.0.2
	github.com/libp2p/go-libp2p-core v0.7.0
	github.com/libp2p/go-libp2p-daemon v0.2.2
	github.com/libp2p/go-libp2p-kad-dht v0.11.0
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/libp2p/go-libp2p-record v0.1.3
	github.com/urfave/cli/v2 v2.3.0
	github.com/wangjia184/sortedset v0.0.0-20160527075905-f5d03557ba30 // indirect
	go.uber.org/fx v1.13.1
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace github.com/filecoin-project/lotus => ./extern/lotus-modified

replace github.com/filecoin-project/go-fil-markets => ./extern/go-fil-markets-modified

replace github.com/supranational/blst => ./extern/lotus-modified/extern/blst

replace github.com/filecoin-project/filecoin-ffi => ./extern/lotus-modified/extern/filecoin-ffi

replace github.com/libp2p/go-libp2p => github.com/jimpick/go-libp2p v0.3.2-0.20201217033239-c003a802f4a7

replace github.com/libp2p/go-reuseport-transport => github.com/jimpick/go-reuseport-transport v0.0.5-0.20201019202422-85fd62f8a44c

replace github.com/libp2p/go-ws-transport => github.com/jimpick/go-ws-transport v0.1.1-0.20201116042118-5dd07d9df8ce
