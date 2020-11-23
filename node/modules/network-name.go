package modules

import (
	"github.com/filecoin-project/lotus/api"
	"go.uber.org/fx"

	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/modules/helpers"
)

func NetworkName(mctx helpers.MetricsCtx, lc fx.Lifecycle, nodeAPI api.FullNode) (dtypes.NetworkName, error) {
	ctx := helpers.LifecycleCtx(mctx, lc)
	return nodeAPI.StateNetworkName(ctx)
}
