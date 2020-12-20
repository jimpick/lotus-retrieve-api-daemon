package api

import (
	"context"

	"github.com/filecoin-project/go-address"
	lapi "github.com/filecoin-project/lotus/api"
	marketevents "github.com/filecoin-project/lotus/markets/loggers"
	"github.com/ipfs/go-cid"
)

// RetrieveAPI implements API passing calls to user-provided function values.
type RetrieveAPI interface {
	// ClientHasLocal indicates whether a certain CID is locally stored.
	ClientHasLocal(ctx context.Context, root cid.Cid) (bool, error)
	// ClientFindData identifies peers that have a certain file, and returns QueryOffers (one per peer).
	ClientFindData(ctx context.Context, root cid.Cid, piece *cid.Cid) ([]lapi.QueryOffer, error)
	// ClientMinerQueryOffer returns a QueryOffer for the specific miner and file.
	ClientMinerQueryOffer(ctx context.Context, miner address.Address, root cid.Cid, piece *cid.Cid) (lapi.QueryOffer, error)
	// ClientRetrieve initiates the retrieval of a file, as specified in the order.
	ClientRetrieve(ctx context.Context, order lapi.RetrievalOrder, ref *lapi.FileRef) error
	// ClientRetrieveWithEvents initiates the retrieval of a file, as specified in the order, and provides a channel
	// of status updates.
	ClientRetrieveWithEvents(ctx context.Context, order lapi.RetrievalOrder, ref *lapi.FileRef) (<-chan marketevents.RetrievalEvent, error)
}
