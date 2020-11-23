module github.com/jimpick/lotus-retrieve-api-daemon

go 1.15

require (
	github.com/filecoin-project/go-fil-markets v1.0.6
	github.com/filecoin-project/go-jsonrpc v0.1.2-0.20201008195726-68c6a2704e49
	github.com/filecoin-project/lotus v1.2.1
	github.com/libp2p/go-libp2p-core v0.7.0
	github.com/libp2p/go-libp2p-kad-dht v0.11.0
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/libp2p/go-libp2p-record v0.1.3
	github.com/urfave/cli/v2 v2.3.0
	go.uber.org/fx v1.13.1
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace github.com/filecoin-project/lotus => extern/lotus-modified

replace github.com/filecoin-project/go-fil-markets => extern/go-fil-markets-modified
