package node

import (
	"context"
	"errors"
	"time"

	ci "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	record "github.com/libp2p/go-libp2p-record"

	"github.com/filecoin-project/go-fil-markets/discovery"
	discoveryimpl "github.com/filecoin-project/go-fil-markets/discovery/impl"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/journal"
	"github.com/filecoin-project/lotus/lib/peermgr"
	"github.com/filecoin-project/lotus/node/config"
	"github.com/filecoin-project/lotus/node/modules"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/modules/helpers"
	"github.com/filecoin-project/lotus/node/modules/lp2p"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/jimpick/lotus-retrieve-api-daemon/api"
	"github.com/jimpick/lotus-retrieve-api-daemon/node/impl"
	rmodules "github.com/jimpick/lotus-retrieve-api-daemon/node/modules"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

// From node/builder.go

// special is a type used to give keys to modules which
//  can't really be identified by the returned type
type special struct{ id int }

type invoke int

//nolint:golint
var (
	DefaultTransportsKey = special{0}  // Libp2p option
	DiscoveryHandlerKey  = special{2}  // Private type
	AddrsFactoryKey      = special{3}  // Libp2p option
	SmuxTransportKey     = special{4}  // Libp2p option
	RelayKey             = special{5}  // Libp2p option
	SecurityKey          = special{6}  // Libp2p option
	BaseRoutingKey       = special{7}  // fx groups + multiret
	NatPortMapKey        = special{8}  // Libp2p option
	ConnectionManagerKey = special{9}  // Libp2p option
	AutoNATSvcKey        = special{10} // Libp2p option
	BandwidthReporterKey = special{11} // Libp2p option
)

// Invokes are called in the order they are defined.
//nolint:golint
const (
	// InitJournal at position 0 initializes the journal global var as soon as
	// the system starts, so that it's available for all other components.
	InitJournalKey = invoke(iota)

	// libp2p

	PstoreAddSelfKeysKey
	RunPeerMgrKey

	// daemon
	ExtractApiKey

	_nInvokes // keep this last
)

type Settings struct {
	// modules is a map of constructors for DI
	//
	// In most cases the index will be a reflect. Type of element returned by
	// the constructor, but for some 'constructors' it's hard to specify what's
	// the return type should be (or the constructor returns fx group)
	modules map[interface{}]fx.Option

	// invokes are separate from modules as they can't be referenced by return
	// type, and must be applied in correct order
	invokes []fx.Option

	nodeType repo.RepoType

	Online bool // Online option applied
	Config bool // Config option applied

}

func Repo(r repo.Repo) Option {
	return func(settings *Settings) error {
		lr, err := r.Lock(settings.nodeType)
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
		Override(new(lp2p.BaseIpfsRouting), lp2p.DHTRouting(dht.ModeAuto)),

		Override(DiscoveryHandlerKey, lp2p.DiscoveryHandler),
		Override(AddrsFactoryKey, lp2p.AddrsFactory(nil, nil)),
		Override(SmuxTransportKey, lp2p.SmuxTransport(true)),
		Override(RelayKey, lp2p.NoRelay()),
		Override(SecurityKey, lp2p.Security(true, false)),

		Override(BaseRoutingKey, lp2p.BaseRouting),
		Override(new(routing.Routing), lp2p.Routing),

		Override(NatPortMapKey, lp2p.NatPortMap),
		Override(BandwidthReporterKey, lp2p.BandwidthCounter),

		Override(ConnectionManagerKey, lp2p.ConnectionManager(50, 200, 20*time.Second, nil)),
		Override(AutoNATSvcKey, lp2p.AutoNATService),

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
			s.nodeType = repo.RetrieveAPI
			return nil
		},
		func(s *Settings) error {
			resAPI := &impl.RetrieveAPI{}
			s.invokes[ExtractApiKey] = fx.Populate(resAPI)
			*out = resAPI
			return nil
		},
	)
}

func defaults() []Option {
	return []Option{
		// global system journal.
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
		modules: map[interface{}]fx.Option{},
		invokes: make([]fx.Option, _nInvokes),
	}

	// apply module options in the right order
	if err := Options(Options(defaults()...), Options(opts...))(&settings); err != nil {
		return nil, xerrors.Errorf("applying node options failed: %w", err)
	}

	// gather constructors for fx.Options
	ctors := make([]fx.Option, 0, len(settings.modules))
	for _, opt := range settings.modules {
		ctors = append(ctors, opt)
	}

	// fill holes in invokes for use in fx.Options
	for i, opt := range settings.invokes {
		if opt == nil {
			settings.invokes[i] = fx.Options()
		}
	}

	app := fx.New(
		fx.Options(ctors...),
		fx.Options(settings.invokes...),

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
