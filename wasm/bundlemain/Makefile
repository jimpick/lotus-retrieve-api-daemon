export GOFLAGS=-tags=clientretrieve
export GOOS=js
export GOARCH=wasm

main.wasm:
	go build -o main.wasm

nodemon:
	nodemon --ext go \
		--watch . --watch ../.. \
		--watch ../../extern/lotus-modified \
		--watch ../../extern/go-fil-markets-modified \
		--watch ../../../../browser-markets/go-ws-transport-0xproject-feat-wss-dialing \
		--watch ../../../../browser-markets/go-libp2p-daemon-ws \
		--watch ../../../../browser-markets/go-jsonrpc-wasm \
		--exec go build -o main.wasm \
		--signal SIGTERM

clean:
	rm -f *.wasm *.wasm.gz
