module github.com/jimpick/lotus-retrieve-api-daemon

go 1.15

require (
	github.com/filecoin-project/go-fil-markets v1.0.10
	github.com/filecoin-project/go-jsonrpc v0.1.2
	github.com/filecoin-project/lotus v1.2.1
	github.com/filecoin-project/specs-actors v0.9.13
	github.com/libp2p/go-libp2p-core v0.7.0
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

replace github.com/supranational/blst => ./extern/lotus-old/extern/fil-blst/blst

replace github.com/filecoin-project/fil-blst => ./extern/lotus-old/extern/fil-blst

replace github.com/filecoin-project/filecoin-ffi => ./extern/lotus-old/extern/filecoin-ffi
