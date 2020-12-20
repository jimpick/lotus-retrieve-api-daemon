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

	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/discovery"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	rmodules "github.com/filecoin-project/lotus/cmd/lotus-retrieve-api-daemon/node/modules"
	"github.com/filecoin-project/lotus/lib/peermgr"
	"github.com/filecoin-project/lotus/node/impl"
	"github.com/filecoin-project/lotus/node/modules"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/modules/helpers"
	"github.com/filecoin-project/lotus/node/modules/lp2p"
	"github.com/filecoin-project/lotus/node/repo"
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
		/*
			c, err := lr.Config()
			if err != nil {
				return err
			}
		*/

		return Options(
			Override(new(repo.LockedRepo), modules.LockedRepo(lr)), // module handles closing

			Override(new(dtypes.MetadataDS), modules.Datastore),
			// Override(new(dtypes.ChainBlockstore), modules.ChainBlockstore),

			Override(new(dtypes.ClientImportMgr), modules.ClientImportMgr),
			Override(new(dtypes.ClientMultiDstore), rmodules.ClientMultiDatastore),

			Override(new(dtypes.ClientBlockstore), modules.ClientBlockstore),
			Override(new(dtypes.ClientRetrievalStoreManager), rmodules.ClientRetrievalStoreManager),
			Override(new(ci.PrivKey), lp2p.PrivKey),
			Override(new(ci.PubKey), ci.PrivKey.GetPublic),
			Override(new(peer.ID), peer.IDFromPublicKey),

			Override(new(types.KeyStore), modules.KeyStore),

			// Override(new(*dtypes.APIAlg), modules.APISecret),

			// ApplyIf(isType(repo.FullNode), ConfigFullNode(c)),
			// ApplyIf(isType(repo.StorageMiner), ConfigStorageMiner(c)),
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

		/*
			Override(new(*dtypes.ScoreKeeper), lp2p.ScoreKeeper),
			Override(new(*pubsub.PubSub), lp2p.GossipSub),
			Override(new(*config.Pubsub), func(bs dtypes.Bootstrapper) *config.Pubsub {
				return &config.Pubsub{
					Bootstrapper: bool(bs),
				}
			}),
		*/

		Override(PstoreAddSelfKeysKey, lp2p.PstoreAddSelfKeys),
		// Override(StartListeningKey, lp2p.StartListening(config.DefaultFullNode().Libp2p.ListenAddresses)),
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

		/*
			// common
			Override(new(*slashfilter.SlashFilter), modules.NewSlashFilter),

			// Full node

			ApplyIf(isType(repo.FullNode),
				// TODO: Fix offline mode

		*/
		Override(new(dtypes.BootstrapPeers), modules.BuiltinBootstrap),
		/*
			Override(new(dtypes.DrandBootstrap), modules.DrandBootstrap),
			Override(new(dtypes.DrandSchedule), modules.BuiltinDrandConfig),

			Override(HandleIncomingMessagesKey, modules.HandleIncomingMessages),

			Override(new(ffiwrapper.Verifier), ffiwrapper.ProofVerifier),
			Override(new(vm.SyscallBuilder), vm.Syscalls),
			Override(new(*store.ChainStore), modules.ChainStore),
			Override(new(*stmgr.StateManager), stmgr.NewStateManager),
			Override(new(*wallet.Wallet), wallet.NewWallet),

			Override(new(dtypes.ChainGCLocker), blockstore.NewGCLocker),
			Override(new(dtypes.ChainGCBlockstore), modules.ChainGCBlockstore),
			Override(new(dtypes.ChainBitswap), modules.ChainBitswap),
			Override(new(dtypes.ChainBlockService), modules.ChainBlockService),

			// Filecoin services
			// We don't want the SyncManagerCtor to be used as an fx constructor, but rather as a value.
			// It will be called implicitly by the Syncer constructor.
			Override(new(chain.SyncManagerCtor), func() chain.SyncManagerCtor { return chain.NewSyncManager }),
			Override(new(*chain.Syncer), modules.NewSyncer),
			Override(new(exchange.Client), exchange.NewClient),
			Override(new(*messagepool.MessagePool), modules.MessagePool),

			Override(new(modules.Genesis), modules.ErrorGenesis),
			Override(new(dtypes.AfterGenesisSet), modules.SetGenesis),
			Override(SetGenesisKey, modules.DoSetGenesis),
		*/
		Override(new(dtypes.NetworkName), rmodules.NetworkName),
		/*
			Override(new(*hello.Service), hello.NewHelloService),
			Override(new(exchange.Server), exchange.NewServer),
		*/
		Override(new(*peermgr.PeerMgr), peermgr.NewPeerMgr),
		Override(new(dtypes.Graphsync), rmodules.Graphsync),
		/*
			Override(new(*dtypes.MpoolLocker), new(dtypes.MpoolLocker)),

			Override(RunHelloKey, modules.RunHello),
			Override(RunChainExchangeKey, modules.RunChainExchange),
		*/
		Override(RunPeerMgrKey, modules.RunPeerMgr),
		/*
			Override(HandleIncomingBlocksKey, modules.HandleIncomingBlocks),
		*/
		Override(new(*discovery.Local), modules.NewLocalDiscovery),
		Override(new(retrievalmarket.PeerResolver), modules.RetrievalResolver),
		Override(new(retrievalmarket.RetrievalClient), rmodules.RetrievalClient),
		Override(new(dtypes.ClientDatastore), modules.NewClientDatastore),
		Override(new(dtypes.ClientDataTransfer), modules.NewClientGraphsyncDataTransfer),
		/*
				Override(new(modules.ClientDealFunds), modules.NewClientDealFunds),
				Override(new(storagemarket.StorageClient), modules.StorageClient),
				Override(new(storagemarket.StorageClientNode), storageadapter.NewClientNodeAdapter),
				Override(new(beacon.Schedule), modules.RandomSchedule),

				Override(new(*paychmgr.Store), paychmgr.NewStore),
				Override(new(*paychmgr.Manager), paychmgr.NewManager),
				Override(new(*market.FundMgr), market.StartFundManager),
				Override(HandlePaymentChannelManagerKey, paychmgr.HandleManager),
				Override(SettlePaymentChannelsKey, settler.SettlePaymentChannels),
			),

			// miner
			ApplyIf(func(s *Settings) bool { return s.nodeType == repo.StorageMiner },
				Override(new(api.Common), From(new(common.CommonAPI))),
				Override(new(sectorstorage.StorageAuth), modules.StorageAuth),

				Override(new(*stores.Index), stores.NewIndex),
				Override(new(stores.SectorIndex), From(new(*stores.Index))),
				Override(new(dtypes.MinerID), modules.MinerID),
				Override(new(dtypes.MinerAddress), modules.MinerAddress),
				Override(new(*ffiwrapper.Config), modules.ProofsConfig),
				Override(new(stores.LocalStorage), From(new(repo.LockedRepo))),
				Override(new(sealing.SectorIDCounter), modules.SectorIDCounter),
				Override(new(*sectorstorage.Manager), modules.SectorStorage),
				Override(new(ffiwrapper.Verifier), ffiwrapper.ProofVerifier),

				Override(new(sectorstorage.SectorManager), From(new(*sectorstorage.Manager))),
				Override(new(storage2.Prover), From(new(sectorstorage.SectorManager))),

				Override(new(*sectorblocks.SectorBlocks), sectorblocks.NewSectorBlocks),
				Override(new(*storage.Miner), modules.StorageMiner(config.DefaultStorageMiner().Fees)),
				Override(new(dtypes.NetworkName), modules.StorageNetworkName),

				Override(new(dtypes.StagingMultiDstore), modules.StagingMultiDatastore),
				Override(new(dtypes.StagingBlockstore), modules.StagingBlockstore),
				Override(new(dtypes.StagingDAG), modules.StagingDAG),
				Override(new(dtypes.StagingGraphsync), modules.StagingGraphsync),
				Override(new(retrievalmarket.RetrievalProvider), modules.RetrievalProvider),
				Override(new(dtypes.ProviderDataTransfer), modules.NewProviderDAGServiceDataTransfer),
				Override(new(dtypes.ProviderPieceStore), modules.NewProviderPieceStore),
				Override(new(*storedask.StoredAsk), modules.NewStorageAsk),
				Override(new(dtypes.DealFilter), modules.BasicDealFilter(nil)),
				Override(new(modules.ProviderDealFunds), modules.NewProviderDealFunds),
				Override(new(storagemarket.StorageProvider), modules.StorageProvider),
				Override(new(storagemarket.StorageProviderNode), storageadapter.NewProviderNodeAdapter),
				Override(HandleRetrievalKey, modules.HandleRetrieval),
				Override(GetParamsKey, modules.GetParams),
				Override(HandleDealsKey, modules.HandleDeals),
				Override(new(gen.WinningPoStProver), storage.NewWinningPoStProver),
				Override(new(*miner.Miner), modules.SetupBlockProducer),

				Override(new(dtypes.ConsiderOnlineStorageDealsConfigFunc), modules.NewConsiderOnlineStorageDealsConfigFunc),
				Override(new(dtypes.SetConsiderOnlineStorageDealsConfigFunc), modules.NewSetConsideringOnlineStorageDealsFunc),
				Override(new(dtypes.ConsiderOnlineRetrievalDealsConfigFunc), modules.NewConsiderOnlineRetrievalDealsConfigFunc),
				Override(new(dtypes.SetConsiderOnlineRetrievalDealsConfigFunc), modules.NewSetConsiderOnlineRetrievalDealsConfigFunc),
				Override(new(dtypes.StorageDealPieceCidBlocklistConfigFunc), modules.NewStorageDealPieceCidBlocklistConfigFunc),
				Override(new(dtypes.SetStorageDealPieceCidBlocklistConfigFunc), modules.NewSetStorageDealPieceCidBlocklistConfigFunc),
				Override(new(dtypes.ConsiderOfflineStorageDealsConfigFunc), modules.NewConsiderOfflineStorageDealsConfigFunc),
				Override(new(dtypes.SetConsiderOfflineStorageDealsConfigFunc), modules.NewSetConsideringOfflineStorageDealsFunc),
				Override(new(dtypes.ConsiderOfflineRetrievalDealsConfigFunc), modules.NewConsiderOfflineRetrievalDealsConfigFunc),
				Override(new(dtypes.SetConsiderOfflineRetrievalDealsConfigFunc), modules.NewSetConsiderOfflineRetrievalDealsConfigFunc),
				Override(new(dtypes.SetSealingConfigFunc), modules.NewSetSealConfigFunc),
				Override(new(dtypes.GetSealingConfigFunc), modules.NewGetSealConfigFunc),
				Override(new(dtypes.SetExpectedSealDurationFunc), modules.NewSetExpectedSealDurationFunc),
				Override(new(dtypes.GetExpectedSealDurationFunc), modules.NewGetExpectedSealDurationFunc),
			),
		*/
	)
}

// func FullAPI(out *api.Retrieve) Option {
func RetrieveAPI(out *api.Retrieve) Option {
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
		/*
			// global system journal.
			Override(new(journal.DisabledEvents), func() journal.DisabledEvents {
				if env, ok := os.LookupEnv(EnvJournalDisabledEvents); ok {
					if ret, err := journal.ParseDisabledEvents(env); err == nil {
						return ret
					}
				}
				// fallback if env variable is not set, or if it failed to parse.
				return journal.DefaultDisabledEvents
			}),
			Override(new(journal.Journal), modules.OpenFilesystemJournal),
			Override(InitJournalKey, func(j journal.Journal) {
				journal.J = j // eagerly sets the global journal through fx.Invoke.
			}),
		*/

		Override(new(helpers.MetricsCtx), context.Background),
		Override(new(record.Validator), modules.RecordValidator),

		Override(new(dtypes.Bootstrapper), dtypes.Bootstrapper(false)),
		/*
			Override(new(dtypes.ShutdownChan), make(chan struct{})),

			// Filecoin modules
		*/
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
