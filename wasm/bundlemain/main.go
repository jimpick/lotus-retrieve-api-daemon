package main

import (
	"fmt"

	"github.com/jimpick/lotus-retrieve-api-daemon/wasm/retrievalservice"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-daemon/p2pclient"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	// remote libp2p node for non-wss
	// controlMaddr, _ := multiaddr.NewMultiaddr("/dns4/libp2p-caddy-p2pd.localhost/tcp/9059/wss")
	controlMaddr, _ := multiaddr.NewMultiaddr("/dns4/p2pd.v6z.me/tcp/9059/wss")
	listenMaddr, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	p2pclientNode, err := p2pclient.NewClient(controlMaddr, listenMaddr)
	fmt.Printf("Jim p2pclientNode %v\n", p2pclientNode)
	nodeID, nodeAddrs, err := p2pclientNode.Identify()
	peerInfo := peer.AddrInfo{
		ID:    nodeID,
		Addrs: nodeAddrs,
	}
	fmt.Printf("Jim peerInfo %v\n", peerInfo)
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("p2pclient->p2pd node address:", addrs[0])

	// APIs
	retrievalservice.Start(p2pclientNode)

	println("WASM Go Initialized")

	c := make(chan struct{}, 0)
	<-c // wait forever
}
