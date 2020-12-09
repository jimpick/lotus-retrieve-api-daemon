checkout:
	mkdir -p extern
	cd extern; git clone -b jim-retrieve-api-daemon-modified git@github.com:jimpick/lotus.git lotus-modified
	cd extern; git clone -b jim/retrieve-daemon git@github.com:jimpick/lotus.git lotus-old
	cd extern; git clone -b jim-retrieve-api-daemon-modified git@github.com:jimpick/go-fil-markets.git go-fil-markets-modified
	cd extern; git clone -b jim/v0.6.1-extra-logging git@github.com:jimpick/go-fil-markets.git go-fil-markets-old
	cd extern; git clone -b jim/extra-logging git@github.com:jimpick/go-data-transfer.git go-data-transfer-old
	cd extern; git clone -b jim/more-logging git@github.com:jimpick/go-graphsync.git go-graphsync-old

run:
	go run . daemon
