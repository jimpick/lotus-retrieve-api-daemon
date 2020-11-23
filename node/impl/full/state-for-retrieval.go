package full

import (
	"context"

	"go.uber.org/fx"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
)

type StateForRetrievalAPI struct {
	fx.In

	nodeAPI api.FullNode
}

func (a *StateForRetrievalAPI) StateNetworkName(ctx context.Context) (dtypes.NetworkName, error) {
	return a.nodeAPI.StateNetworkName(ctx)
}
