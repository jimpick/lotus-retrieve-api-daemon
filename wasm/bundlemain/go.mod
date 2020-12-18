module github.com/jimpick/lotus-retrieve-api-daemon/wasm/bundlemain

go 1.15

require (
	github.com/filecoin-project/go-address v0.0.5-0.20201103152444-f2023ef3f5bb
	github.com/filecoin-project/go-fil-markets v1.0.10
	github.com/filecoin-project/go-jsonrpc v0.1.2
	github.com/filecoin-project/lotus v1.2.1
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-datastore v0.4.5
	github.com/ipfs/go-graphsync v0.5.1
	github.com/jimpick/lotus-retrieve-api-daemon v0.0.0-00010101000000-000000000000
	github.com/libp2p/go-libp2p-core v0.7.0
	github.com/libp2p/go-libp2p-daemon v0.2.2
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/libp2p/go-libp2p-record v0.1.3
	github.com/libp2p/go-ws-transport v0.3.1
	github.com/multiformats/go-multiaddr v0.3.1
	go.uber.org/fx v1.13.1
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace github.com/jimpick/lotus-retrieve-api-daemon => ../..

replace github.com/filecoin-project/lotus => ../../extern/lotus-modified

replace github.com/filecoin-project/go-fil-markets => ../../extern/go-fil-markets-modified

replace github.com/libp2p/go-libp2p => github.com/jimpick/go-libp2p v0.3.2-0.20201217033239-c003a802f4a7

replace github.com/libp2p/go-reuseport-transport => github.com/jimpick/go-reuseport-transport v0.0.5-0.20201019202422-85fd62f8a44c

replace github.com/filecoin-project/go-jsonrpc => github.com/jimpick/go-jsonrpc v0.0.0-20201109011442-669bac3b0e93

replace github.com/libp2p/go-libp2p-daemon => ../../../../browser-markets/go-libp2p-daemon-ws

replace github.com/libp2p/go-ws-transport => ../../../../browser-markets/go-ws-transport-0xproject-feat-wss-dialing

replace github.com/multiformats/go-multiaddr => github.com/jimpick/go-multiaddr v0.3.2-0.20201116042404-3634c019a1d6
