export GOFLAGS=-tags=clientretrieve

lotus-retrieve-api-daemon:
	go build .

gen:
	#cd extern/lotus-modified && go run ./gen/main_clientretrieve.go && go generate ./...
	rm -f extern/lotus-modified/api/cbor_gen.go
	cd extern/lotus-modified && go run ./gen/main_clientretrieve.go
	cd extern/lotus-modified && go generate ./...

clean:
	rm -f lotus-retrieve-api-daemon

env:
	go env

why-ffi:
	go list $(GOFLAGS) -e -json -compiled=true -test=true -deps=true . | jq -C '. | select(.Imports) | select(.Imports[] | in({"github.com/filecoin-project/filecoin-ffi": ""})) | .ImportPath'

why-sector-storage:
	go list $(GOFLAGS) -e -json -compiled=true -test=true -deps=true . | jq -C '. | select(.Imports) | select(.Imports[] | in({"github.com/filecoin-project/lotus/extern/sector-storage": ""})) | .ImportPath'

list:
	go list $(GOFLAGS) -e -json -compiled=true -test=true -deps=true . | jq -C .

list-wasm:
	cd wasm/bundlemain; GOOS=js GOARCH=wasm go list $(GOFLAGS) -e -json -compiled=true -test=true -deps=true . | jq -C .

why-ffi-wasm:
	cd wasm/bundlemain; GOOS=js GOARCH=wasm go list $(GOFLAGS) -e -json -compiled=true -test=true -deps=true . | jq -C '. | select(.Imports) | select(.Imports[] | in({"github.com/filecoin-project/filecoin-ffi": ""})) | .ImportPath'

why-badger-wasm:
	cd wasm/bundlemain; GOOS=js GOARCH=wasm go list $(GOFLAGS) -e -json -compiled=true -test=true -deps=true . | jq -C '. | select(.Imports) | select(.Imports[] | in({"github.com/dgraph-io/badger/v2": ""})) | .ImportPath'


checkout:
	mkdir -p extern
	cd extern; git clone -b jim-retrieve-api-daemon-modified git@github.com:jimpick/lotus.git lotus-modified
	cd extern; git clone -b jim-retrieve-api-daemon-modified git@github.com:jimpick/go-fil-markets.git go-fil-markets-modified
	cd extern; git clone -b jim-retrieve-api-daemon-modified git@github.com:jimpick/go-data-transfer.git go-data-transfer-modified
	cd extern; git clone -b jim-retrieve-api-daemon-modified git@github.com:jimpick/go-graphsync.git go-graphsync-modified

run:
	go run . daemon
