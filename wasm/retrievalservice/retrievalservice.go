package retrievalservice

import (
	"context"
	"errors"
	"fmt"
	"syscall/js"
	"time"

	"github.com/filecoin-project/go-fil-markets/discovery"
	discoveryimpl "github.com/filecoin-project/go-fil-markets/discovery/impl"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/journal"
	"github.com/filecoin-project/lotus/node/config"
	"github.com/filecoin-project/lotus/node/modules/lp2p"
	"github.com/filecoin-project/lotus/node/modules/moduleapi"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/jimpick/lotus-retrieve-api-daemon/api"
	. "github.com/jimpick/lotus-utils/fxnodesetup"
	"github.com/libp2p/go-libp2p-daemon/p2pclient"

	ci "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/routing"

	// dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	record "github.com/libp2p/go-libp2p-record"

	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/lib/peermgr"

	// _ "github.com/filecoin-project/lotus/lib/sigs/bls"
	// _ "github.com/filecoin-project/lotus/lib/sigs/secp"

	"github.com/filecoin-project/lotus/node/modules"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/modules/helpers"

	"github.com/jimpick/lotus-retrieve-api-daemon/node/impl"
	rmodules "github.com/jimpick/lotus-retrieve-api-daemon/node/modules"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

var jsHandler js.Value

func Start(p2pclientNode *p2pclient.Client) {
	var retrieveAPI api.RetrieveAPI

	ctx := context.Background()

	r := repo.NewMemory(nil)

	nilRouting, err := lp2p.NilRouting(ctx)
	if err != nil {
		panic(err)
	}

	var fullNodeStruct = apistruct.FullNodeStruct{}
	var nodeAPI lotusapi.FullNode = &fullNodeStruct
	var closer jsonrpc.ClientCloser
	defer func() {
		if closer != nil {
			closer()
		}
	}()

	_, err = New(ctx,
		RetrieveAPI(&retrieveAPI),
		Repo(r),
		Online(),
		Override(new(lp2p.BaseIpfsRouting), nilRouting),
		Override(new(moduleapi.StateModuleAPI), nodeAPI),
		Override(new(moduleapi.ChainModuleAPI), nodeAPI),
		Override(new(moduleapi.PaychModuleAPI), nodeAPI),
		Override(new(*p2pclient.Client), p2pclientNode),
	)
	if err != nil {
		panic(err)
	}

	cbOpt := jsonrpc.WithConnectCallback(func(environment js.Value) {
		fmt.Printf("Jim retrievalservice connectCallBack\n")
		requestsForLotusHandler := environment.Get("requestsForLotusHandler")

		// closer, err := jsonrpc.NewJSMergeClient(context.Background(), requestsForLotusHandler, "Filecoin", []interface{}{&nodeAPI})
		closer, err = jsonrpc.NewJSMergeClient(context.Background(), requestsForLotusHandler, "Filecoin",
			[]interface{}{
				&fullNodeStruct.CommonStruct.Internal,
				&fullNodeStruct.Internal,
			})
		if err != nil {
			fmt.Printf("connecting with lotus failed: %s\n", err)
			panic(err)
		}
	})
	rpcServer := jsonrpc.NewJSServer("connectRetrievalService", cbOpt)
	rpcServer.Register("Filecoin", retrieveAPI)
}

//nolint:golint
var (
	DefaultTransportsKey = Special{0}  // Libp2p option
	DiscoveryHandlerKey  = Special{2}  // Private type
	AddrsFactoryKey      = Special{3}  // Libp2p option
	SmuxTransportKey     = Special{4}  // Libp2p option
	RelayKey             = Special{5}  // Libp2p option
	SecurityKey          = Special{6}  // Libp2p option
	BaseRoutingKey       = Special{7}  // fx groups + multiret
	NatPortMapKey        = Special{8}  // Libp2p option
	ConnectionManagerKey = Special{9}  // Libp2p option
	AutoNATSvcKey        = Special{10} // Libp2p option
	BandwidthReporterKey = Special{11} // Libp2p option
)

// Invokes are called in the order they are defined.
//nolint:golint
const (
	// InitJournal at position 0 initializes the journal global var as soon as
	// the system starts, so that it's available for all other components.
	InitJournalKey = Invoke(iota)

	// libp2p

	PstoreAddSelfKeysKey
	RunPeerMgrKey

	// daemon
	ExtractApiKey

	_nInvokes // keep this last
)

func Repo(r repo.Repo) Option {
	return func(settings *Settings) error {
		lr, err := r.Lock(settings.NodeType)
		if err != nil {
			return err
		}

		return Options(
			Override(new(repo.LockedRepo), modules.LockedRepo(lr)), // module handles closing

			Override(new(dtypes.MetadataDS), modules.Datastore),

			Override(new(dtypes.ClientImportMgr), modules.ClientImportMgr),
			Override(new(dtypes.ClientMultiDstore), modules.ClientMultiDatastore),
			Override(new(dtypes.ClientBlockstore), modules.ClientBlockstore),
			Override(new(dtypes.ClientRetrievalStoreManager), modules.ClientRetrievalStoreManager),

			Override(new(ci.PrivKey), lp2p.PrivKey),
			Override(new(ci.PubKey), ci.PrivKey.GetPublic),
			Override(new(peer.ID), peer.IDFromPublicKey),

			Override(new(types.KeyStore), modules.KeyStore),
		)(settings)
	}
}

func libp2p() Option {
	return Options(
		Override(new(peerstore.Peerstore), pstoremem.NewPeerstore),

		Override(DefaultTransportsKey, lp2p.DefaultTransports),

		Override(new(lp2p.RawHost), lp2p.Host),
		Override(new(host.Host), lp2p.RoutedHost),
		// Override(new(lp2p.BaseIpfsRouting), lp2p.DHTRouting(dht.ModeAuto)),

		Override(DiscoveryHandlerKey, lp2p.DiscoveryHandler),
		Override(AddrsFactoryKey, lp2p.AddrsFactory(nil, nil)),
		Override(SmuxTransportKey, lp2p.SmuxTransport(true)),
		Override(RelayKey, lp2p.NoRelay()),
		Override(SecurityKey, lp2p.Security(true, false)),

		Override(BaseRoutingKey, lp2p.BaseRouting),
		Override(new(routing.Routing), lp2p.Routing),

		Override(BandwidthReporterKey, lp2p.BandwidthCounter),

		Override(ConnectionManagerKey, lp2p.ConnectionManager(50, 200, 20*time.Second, nil)),

		Override(PstoreAddSelfKeysKey, lp2p.PstoreAddSelfKeys),
	)
}

// Online sets up basic libp2p node
func Online() Option {
	return Options(
		// make sure that online is applied before Config.
		// This is important because Config overrides some of Online units
		func(s *Settings) error { s.Online = true; return nil },
		ApplyIf(func(s *Settings) bool { return s.Config },
			Error(errors.New("the Online option must be set before Config option")),
		),

		libp2p(),

		Override(new(dtypes.BootstrapPeers), modules.BuiltinBootstrap),
		Override(new(dtypes.NetworkName), rmodules.NetworkName),
		Override(new(*peermgr.PeerMgr), peermgr.NewPeerMgr),
		Override(new(dtypes.Graphsync), modules.Graphsync(config.DefaultFullNode().Client.SimultaneousTransfers)),
		Override(RunPeerMgrKey, modules.RunPeerMgr),
		Override(new(*discoveryimpl.Local), modules.NewLocalDiscovery),
		Override(new(discovery.PeerResolver), modules.RetrievalResolver),
		Override(new(retrievalmarket.RetrievalClient), modules.RetrievalClient),
		Override(new(dtypes.ClientDatastore), modules.NewClientDatastore),
		Override(new(dtypes.ClientDataTransfer), modules.NewClientGraphsyncDataTransfer),
	)
}

func RetrieveAPI(out *api.RetrieveAPI) Option {
	return Options(
		func(s *Settings) error {
			s.NodeType = repo.RetrieveAPI
			return nil
		},
		func(s *Settings) error {
			resAPI := &impl.RetrieveAPI{}
			s.Invokes[ExtractApiKey] = fx.Populate(resAPI)
			*out = resAPI
			return nil
		},
	)
}

func defaults() []Option {
	return []Option{
		Override(new(journal.DisabledEvents), journal.EnvDisabledEvents),
		Override(new(journal.Journal), modules.OpenFilesystemJournal),

		Override(new(helpers.MetricsCtx), context.Background),
		Override(new(record.Validator), modules.RecordValidator),

		Override(new(dtypes.Bootstrapper), dtypes.Bootstrapper(false)),
	}
}

type StopFunc func(context.Context) error

// New builds and starts new Filecoin node
func New(ctx context.Context, opts ...Option) (StopFunc, error) {
	settings := Settings{
		Modules: map[interface{}]fx.Option{},
		Invokes: make([]fx.Option, _nInvokes),
	}

	// apply module options in the right order
	if err := Options(Options(defaults()...), Options(opts...))(&settings); err != nil {
		return nil, xerrors.Errorf("applying node options failed: %w", err)
	}

	// gather constructors for fx.Options
	ctors := make([]fx.Option, 0, len(settings.Modules))
	for _, opt := range settings.Modules {
		ctors = append(ctors, opt)
	}

	// fill holes in invokes for use in fx.Options
	for i, opt := range settings.Invokes {
		if opt == nil {
			settings.Invokes[i] = fx.Options()
		}
	}

	app := fx.New(
		fx.Options(ctors...),
		fx.Options(settings.Invokes...),

		fx.NopLogger,
	)

	// TODO: we probably should have a 'firewall' for Closing signal
	//  on this context, and implement closing logic through lifecycles
	//  correctly
	if err := app.Start(ctx); err != nil {
		// comment fx.NopLogger few lines above for easier debugging
		return nil, xerrors.Errorf("starting node: %w", err)
	}

	return app.Stop, nil
}
