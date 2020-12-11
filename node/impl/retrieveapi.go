package impl

import (
	"github.com/filecoin-project/lotus/node/impl/client"
	"github.com/jimpick/lotus-retrieve-api-daemon/api"
)

type RetrieveAPI struct {
	client.API
}

var _ api.RetrieveAPI = &RetrieveAPI{}
