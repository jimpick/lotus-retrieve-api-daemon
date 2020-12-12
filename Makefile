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
	go list $(GOFLAGS) -e -json -compiled=true -test=true -deps=true . | jq -C '. | select(.Imports) | select(.Imports[] | contains("github.com/filecoin-project/filecoin-ffi")) | .ImportPath'

why-seed:
	go list $(GOFLAGS) -e -json -compiled=true -test=true -deps=true . | jq -C '. | select(.Imports) | select(.Imports[] | contains("github.com/filecoin-project/lotus/cmd/lotus-seed/seed")) | .ImportPath'

why-sector-storage:
	go list $(GOFLAGS) -e -json -compiled=true -test=true -deps=true . | jq -C '. | select(.Imports) | select(.Imports[] | contains("github.com/filecoin-project/lotus/extern/sector-storage")) | .ImportPath'


list:
	go list $(GOFLAGS) -e -json -compiled=true -test=true -deps=true . | jq -C .

checkout:
	mkdir -p extern
	cd extern; git clone -b jim-retrieve-api-daemon-modified git@github.com:jimpick/lotus.git lotus-modified
	cd extern; git clone -b jim-retrieve-api-daemon-modified git@github.com:jimpick/go-fil-markets.git go-fil-markets-modified
	cd extern; git clone -b jim/retrieve-daemon git@github.com:jimpick/lotus.git lotus-old
	cd extern; git clone -b jim/v0.6.1-extra-logging git@github.com:jimpick/go-fil-markets.git go-fil-markets-old
	cd extern; git clone -b jim/extra-logging git@github.com:jimpick/go-data-transfer.git go-data-transfer-old
	cd extern; git clone -b jim/more-logging git@github.com:jimpick/go-graphsync.git go-graphsync-old

run:
	go run . daemon
