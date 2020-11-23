package modules

import (
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/modules/helpers"
	graphsyncimpl "github.com/ipfs/go-graphsync/impl"
	gsnet "github.com/ipfs/go-graphsync/network"
	"github.com/ipfs/go-graphsync/storeutil"
	"github.com/libp2p/go-libp2p-core/host"
	"go.uber.org/fx"
)

// Graphsync creates a graphsync instance from the given loader and storer
func Graphsync(mctx helpers.MetricsCtx, lc fx.Lifecycle, clientBs dtypes.ClientBlockstore, h host.Host) (dtypes.Graphsync, error) {
	graphsyncNetwork := gsnet.NewFromLibp2pHost(h)
	loader := storeutil.LoaderForBlockstore(clientBs)
	storer := storeutil.StorerForBlockstore(clientBs)
	gs := graphsyncimpl.New(helpers.LifecycleCtx(mctx, lc), graphsyncNetwork, loader, storer, graphsyncimpl.RejectAllRequestsByDefault())
	return gs, nil
}
