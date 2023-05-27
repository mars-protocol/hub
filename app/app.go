package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"

	// tendermint
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"

	// cosmos SDK

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	api "github.com/cosmos/cosmos-sdk/server/api"
	config "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"

	// core modules
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	// customized core modules
	customgov "github.com/mars-protocol/hub/v2/x/gov"
	customgovkeeper "github.com/mars-protocol/hub/v2/x/gov/keeper"

	// ibc modules
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	ibctestingtypes "github.com/cosmos/ibc-go/v7/testing/types"

	// wasm modules
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	// mars modules
	"github.com/mars-protocol/hub/v2/x/envoy"
	envoykeeper "github.com/mars-protocol/hub/v2/x/envoy/keeper"
	envoytypes "github.com/mars-protocol/hub/v2/x/envoy/types"
	"github.com/mars-protocol/hub/v2/x/incentives"
	incentiveskeeper "github.com/mars-protocol/hub/v2/x/incentives/keeper"
	incentivestypes "github.com/mars-protocol/hub/v2/x/incentives/types"
	"github.com/mars-protocol/hub/v2/x/safety"
	safetykeeper "github.com/mars-protocol/hub/v2/x/safety/keeper"
	safetytypes "github.com/mars-protocol/hub/v2/x/safety/types"

	"github.com/mars-protocol/hub/v2/app/upgrades"
	v2 "github.com/mars-protocol/hub/v2/app/upgrades/v2"

	marswasm "github.com/mars-protocol/hub/v2/app/wasm"
	marsdocs "github.com/mars-protocol/hub/v2/docs"
)

const (
	AccountAddressPrefix = "mars"
	Name                 = "mars"

	// BondDenom is the staking token's denomination
	BondDenom = "umars"

	// If EnabledSpecificProposals is "", and this is "true", then enable all
	// x/wasm proposals.
	// If EnabledSpecificProposals is "", and this is not "true", then disable
	// all x/wasm proposals.
	ProposalsEnabled = "true"

	// If set to non-empty string it must be comma-separated list of values that
	// are all a subset of "EnableAllProposals" (takes precedence over
	// ProposalsEnabled)
	// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
	EnableSpecificProposals = ""
)

var (
	// DefaultNodeHome the default home directory for the app daemon
	DefaultNodeHome string

	// ModuleBasics defines the module `BasicManager`, which is in charge of
	// setting up basic, non- dependent module elements, such as codec
	// registration and genesis verification
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		crisis.AppModuleBasic{},
		distr.AppModuleBasic{}, // distr AppModuleBasic is not customized, so we just use the vanilla one
		evidence.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		genutil.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
				upgradeclient.LegacyProposalHandler,
				upgradeclient.LegacyCancelProposalHandler,
				ibcclientclient.UpdateClientProposalHandler,
				ibcclientclient.UpgradeProposalHandler,
			},
		),
		params.AppModuleBasic{},
		slashing.AppModuleBasic{},
		staking.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		ica.AppModuleBasic{},
		wasm.AppModuleBasic{},
		incentives.AppModuleBasic{},
		safety.AppModuleBasic{},
		envoy.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		govtypes.ModuleName:            {authtypes.Burner},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:            nil,
		wasm.ModuleName:                {authtypes.Burner},
		incentivestypes.ModuleName:     nil,
		safetytypes.ModuleName:         nil,
		envoytypes.ModuleName:          nil,
	}

	// scheduled upgrades and forks
	Upgrades = []upgrades.Upgrade{v2.Upgrade}
	Forks    = []upgrades.Fork{}
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".mars")
}

//------------------------------------------------------------------------------
// Mars app
//------------------------------------------------------------------------------

// `MarsApp` must implement `simapp.App` and `servertypes.Application` interfaces
var (
	_ runtime.AppI            = (*MarsApp)(nil)
	_ servertypes.Application = (*MarsApp)(nil)
	_ ibctesting.TestingApp   = (*MarsApp)(nil)
)

// MarsApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type MarsApp struct {
	// baseapp
	*baseapp.BaseApp

	// codecs
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec // make codec public for testing purposes
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// invariant check period
	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	AuthzKeeper           authzkeeper.Keeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	GovKeeper             customgovkeeper.Keeper // replaces the vanilla gov keeper with our custom one
	ParamsKeeper          paramskeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper // must be a pointer, so we can `SetRouter` on it correctly
	IBCTransferKeeper     ibctransferkeeper.Keeper
	ICAControllerKeeper   icacontrollerkeeper.Keeper
	ICAHostKeeper         icahostkeeper.Keeper
	WasmKeeper            wasm.Keeper
	IncentivesKeeper      incentiveskeeper.Keeper
	SafetyKeeper          safetykeeper.Keeper
	EnvoyKeeper           envoykeeper.Keeper

	// make scoped keepers public for testing purposes
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedIBCTransferKeeper   capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper          capabilitykeeper.ScopedKeeper

	// module manager and configurator
	ModuleManager *module.Manager
	configurator  module.Configurator
}

// NewMarsApp creates and initializes a new `MarsApp` instance
func NewMarsApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	enabledProposals []wasm.ProposalType,
	appOpts servertypes.AppOptions,
	wasmOpts []wasm.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *MarsApp {
	encodingConfig := MakeEncodingConfig()

	appCodec, legacyAmino := encodingConfig.Codec, encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(Name, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		authzkeeper.StoreKey,
		banktypes.StoreKey,
		capabilitytypes.StoreKey,
		distrtypes.StoreKey,
		evidencetypes.StoreKey,
		feegrant.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		consensusparamtypes.StoreKey,
		slashingtypes.StoreKey,
		stakingtypes.StoreKey,
		upgradetypes.StoreKey,
		ibcexported.StoreKey,
		ibctransfertypes.StoreKey,
		icacontrollertypes.StoreKey,
		icahosttypes.StoreKey,
		wasm.StoreKey,
		incentivestypes.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	// load state streaming if enabled
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, logger, keys); err != nil {
		logger.Error("failed to load state streaming", "err", err)
		os.Exit(1)
	}

	app := &MarsApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	// **** create keepers ****

	// address of the gov module account
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	app.ParamsKeeper = initParamsKeeper(
		appCodec,
		legacyAmino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, keys[consensusparamtypes.StoreKey], authtypes.NewModuleAddress(govtypes.ModuleName).String())
	bApp.SetParamStore(&app.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	app.ScopedIBCKeeper = app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	app.ScopedIBCTransferKeeper = app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	app.ScopedICAControllerKeeper = app.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	app.ScopedICAHostKeeper = app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	app.ScopedWasmKeeper = app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		maccPerms,
		AccountAddressPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(keys[authzkeeper.StoreKey], appCodec, app.MsgServiceRouter(), app.AccountKeeper)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		codec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		getSubspace(app, banktypes.ModuleName),
		getBlockedModuleAccountAddrs(app), // NOTE: fee collector & safety fund are excluded from blocked addresses
	)
	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		keys[crisistypes.StoreKey],
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(keys[authzkeeper.StoreKey], appCodec, app.MsgServiceRouter(), app.AccountKeeper)

	// get skipUpgradeHeights from the app options
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	// set the governance module account as the authority for conducting upgrades
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// staking keeper and its dependencies
	// NOTE: the order here (e.g. evidence keeper depends on slashing keeper, so
	// must be defined after it)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		getSubspace(app, stakingtypes.ModuleName),
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		keys[slashingtypes.StoreKey],
		app.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		keys[evidencetypes.StoreKey],
		app.StakingKeeper,
		app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain
	// these hooks.
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		getSubspace(app, ibcexported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper, app.ScopedIBCKeeper,
	)
	app.IBCTransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		getSubspace(app, ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		app.ScopedIBCTransferKeeper,
	)
	app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec,
		keys[icacontrollertypes.StoreKey],
		getSubspace(app, icacontrollertypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.ScopedICAControllerKeeper,
		app.MsgServiceRouter(),
	)
	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		keys[icahosttypes.StoreKey],
		getSubspace(app, icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.ScopedICAHostKeeper,
		app.MsgServiceRouter(),
	)

	// create static IBC router, add transfer route, then set and seal it
	app.IBCKeeper.SetRouter(initIBCRouter(app))

	// load configs for wasm module
	wasmDir := filepath.Join(homePath, "data")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	supportedCapabilities := "iterator,staking,stargate,cosmwasm_1_1,cosmwasm_1_2"

	// register wasm bindings of Mars modules here
	wasmOpts = append(marswasm.RegisterCustomPlugins(&app.DistrKeeper), wasmOpts...)

	// create wasm keeper
	app.WasmKeeper = wasm.NewKeeper(
		appCodec,
		keys[wasm.StoreKey],
		getSubspace(app, wasm.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.DistrKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.ScopedWasmKeeper,
		app.IBCTransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedCapabilities,
		wasmOpts...,
	)

	// mars module keepers
	app.IncentivesKeeper = incentiveskeeper.NewKeeper(
		codec, keys[incentivestypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		app.StakingKeeper,
		authority,
	)
	app.SafetyKeeper = safetykeeper.NewKeeper(
		app.AccountKeeper,
		app.BankKeeper,
		authority,
	)
	app.EnvoyKeeper = envoykeeper.NewKeeper(
		app.appCodec,
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.ICAControllerKeeper,
		app.MsgServiceRouter(),
		[]string{authority},
	)

	// finally, create gov keeper
	//
	// here we use the customized gov keeper, which requires an additional
	// `wasmKeeper` parameter compared to the vanilla govkeeper
	app.GovKeeper = customgovkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		getSubspace(app, govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		app.WasmKeeper,
		initGovRouter(app),
		app.MsgServiceRouter(),
		// the vanilla gov module by default has a 255-character limit for
		// proposal metadata, because it assumes proposals will be stored off-
		// chain and only an IPFS hash will be uploaded on-chain.
		//
		// in this, imo, a big mistake! governance proposals are an important
		// part of a blockchain's history, so their data should be persisted on-
		// chain. it's not like they will take a huge storage space anyways.
		//
		// at Mars we require all proposal metadata to be stored on-chain, and
		// they must conform to a schema (see the customgov module's README.)
		// for this to work, we use u64::MAX as the max allowed length.
		govtypes.Config{
			MaxMetadataLen: ^uint64(0), // ^ is the bitwise NOT operator
		},
	)

	// **** module options ****

	// NOTE: We may consider parsing `appOpts` inside module constructors. For
	// the moment we prefer to be more strict in what arguments the modules
	// expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later
	// modified must be passed by reference here.
	app.ModuleManager = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper,
			app.StakingKeeper,
			app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
		customgov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		params.NewAppModule(app.ParamsKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(app.IBCTransferKeeper),
		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		incentives.NewAppModule(app.IncentivesKeeper),
		safety.NewAppModule(app.SafetyKeeper),
		envoy.NewAppModule(app.EnvoyKeeper),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)), // always be last to make sure that it checks for all invariants and not only part of them
	)

	// During begin block, slashing happens after `distr.BeginBlocker` so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// `CanWithdrawInvariant` invariant.
	// NOTE: staking module is required if `HistoricalEntries` param > 0
	app.ModuleManager.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		wasm.ModuleName,
		incentivestypes.ModuleName,
		safetytypes.ModuleName,
		envoytypes.ModuleName,
	)

	app.ModuleManager.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		wasm.ModuleName,
		incentivestypes.ModuleName,
		safetytypes.ModuleName,
		envoytypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any
	// capabilities so that other modules that want to create or claim
	// capabilities afterwards in `InitChain` can do so safely.
	app.ModuleManager.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		feegrant.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		wasm.ModuleName,
		incentivestypes.ModuleName,
		safetytypes.ModuleName,
		envoytypes.ModuleName,
	)

	// Uncomment if you want to set a custom migration order here.
	// app.ModuleManager.SetOrderMigrations(custom order)

	app.ModuleManager.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.ModuleManager.RegisterServices(app.configurator)

	// RegisterUpgradeHandlers is used for registering any on-chain upgrades.
	// Make sure it's called after `app.ModuleManager` and `app.configurator` are set.
	app.RegisterUpgradeHandlers()

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.ModuleManager.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// add test gRPC service for testing gRPC queries in isolation
	// testdata_pulsar.RegisterQueryServer(app.GRPCQueryRouter(), testdata_pulsar.QueryImpl{})

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.ModuleManager.Modules, overrideModules)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.setAnteHandler(encodingConfig.TxConfig, wasmConfig, keys[wasm.StoreKey])

	// must be before Loading version
	// requires the snapshot store to be created and registered as a BaseAppOption
	// see cmd/wasmd/root.go: 206 - 214 approx
	if manager := app.SnapshotManager(); manager != nil {
		err := manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))
		}
	}
	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.setAnteHandler(encodingConfig.TxConfig, wasmConfig, keys[wasm.StoreKey])

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

//------------------------------------------------------------------------------
// Implement `marsapptypes.App` interface for MarsApp
//
// `ExportAppStateAndValidators` is implemented in ./export.go so no need to
// reimplement here.
//------------------------------------------------------------------------------

// Name returns the app's name
func (app *MarsApp) Name() string {
	return app.BaseApp.Name()
}

// LegacyAmino returns the app's legacy amino codec
func (app *MarsApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// BeginBlocker application updates every begin block
func (app *MarsApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	app.beginBlockForks(ctx)
	return app.ModuleManager.BeginBlock(ctx, req)
}

// EndBlocker application updates on every end block
func (app *MarsApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.ModuleManager.EndBlock(ctx, req)
}

// InitChainer application updates at chain initialization
func (app *MarsApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap())

	return app.ModuleManager.InitGenesis(ctx, app.Codec, genesisState)
}

// LoadHeight loads a particular height
func (app *MarsApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs return the app's module account addresses
func (app *MarsApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// NOTE: we don't actually use simulation manager anywhere in this project, so
// we simply return nil.
func (app *MarsApp) SimulationManager() *module.SimulationManager {
	return nil
}

//------------------------------------------------------------------------------
// Implement `servertypes.Application` interface for `MarsApp`
//
// `RegisterGRPCServer` is already implemented by `BaseApp` and inherited by MarsApp
//------------------------------------------------------------------------------

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *MarsApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// register new tx routes from grpc-gateway
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register new tendermint queries routes from grpc-gateway
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register grpc-gateway routes for all modules
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API
	if apiConfig.Swagger {
		apiSvr.Router.Handle("/swagger.yml", http.FileServer(http.FS(marsdocs.Swagger)))
		apiSvr.Router.HandleFunc("/swagger/", marsdocs.Handler(Name, "/swagger.yml"))
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *MarsApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(
		app.BaseApp.GRPCQueryRouter(),
		clientCtx,
		app.BaseApp.Simulate,
		app.interfaceRegistry,
	)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *MarsApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

//------------------------------------------------------------------------------
// Implement `ibctesting.TestingApp` interface for `MarsApp`
//------------------------------------------------------------------------------

func (app *MarsApp) AppCodec() codec.Codec {
	return app.appCodec
}

func (app *MarsApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

func (app *MarsApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

func (app *MarsApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

func (app *MarsApp) GetStakingKeeper() ibctestingtypes.StakingKeeper {
	return app.StakingKeeper
}

func (app *MarsApp) GetTxConfig() client.TxConfig {
	return MakeEncodingConfig().TxConfig
}

//------------------------------------------------------------------------------
// Upgrades and forks
//------------------------------------------------------------------------------

func (app *MarsApp) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk: %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range Upgrades {
		if upgradeInfo.Name == upgrade.UpgradeName {
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &upgrade.StoreUpgrades))
		}
	}
}

func (app *MarsApp) setupUpgradeHandlers() {
	for _, upgrade := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(app.ModuleManager, app.configurator),
		)
	}
}

func (app *MarsApp) beginBlockForks(ctx sdk.Context) {
	for _, fork := range Forks {
		if ctx.BlockHeight() == fork.UpgradeHeight {
			fork.BeginForkLogic(ctx)
			return
		}
	}
}

//------------------------------------------------------------------------------
// Helpers
//------------------------------------------------------------------------------

// returns a param subspace for a given module name
func getSubspace(app *MarsApp, moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// getEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
// produce a list of enabled proposals to pass into wasmd app.
func getEnabledProposals() []wasm.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasm.EnableAllProposals
		}
		return wasm.DisableAllProposals
	}

	chunks := strings.Split(EnableSpecificProposals, ",")

	proposals, err := wasm.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}

	return proposals
}

// getBlockedModuleAccountAddrs returns all the app's blocked module account
// addresses.
//
// Specifically, we allow the following module accounts to receive funds:
//
//   - `fee_collector` and `safety_fund`, so that protocol revenue can be sent
//     from outposts to the hub via IBC fungible token transfers
//
//   - `incentives`, so that the incentives module can draw funds from the
//     community pool in order to
//     create new incentives schedules upon successful governance proposals
//
//   - `envoy`, so that it can draw funds from the community pool;
//     additionally, if an ICS-20 packet times out, it can receive the refund.
//
// Further note on the 2nd point: the distrkeeper's `DistributeFromFeePool`
// function uses bankkeeper's `SendCoinsFromModuleToAccount` instead of
// `SendCoinsFromModuleToModule`. If it had used `FromModuleToModule`
// then we won't need to allow incentives and envoy module accounts to receive
// funds here.
//
// Forked from: https://github.com/cosmos/gaia/pull/1493
func getBlockedModuleAccountAddrs(app *MarsApp) map[string]bool {
	modAccAddrs := app.ModuleAccountAddrs()

	delete(modAccAddrs, authtypes.NewModuleAddress(authtypes.FeeCollectorName).String())
	delete(modAccAddrs, authtypes.NewModuleAddress(incentivestypes.ModuleName).String())
	delete(modAccAddrs, authtypes.NewModuleAddress(safetytypes.ModuleName).String())
	delete(modAccAddrs, authtypes.NewModuleAddress(envoytypes.ModuleName).String())

	return modAccAddrs
}

// TxConfig returns WasmApp's TxConfig
func (app *MarsApp) TxConfig() client.TxConfig {
	return app.txConfig
}

func (app *MarsApp) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// GetMaccPerms returns a copy of the module account permissions
//
// NOTE: This is solely to be used for testing purposes.
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range GetMaccPerms() {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	// allow the following addresses to receive funds
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	return modAccAddrs
}

// initializes params keeper and its subspaces
func initParamsKeeper(codec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(codec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.StoreKey)
	paramsKeeper.Subspace(icahosttypes.StoreKey)
	paramsKeeper.Subspace(wasm.ModuleName)

	return paramsKeeper
}

// initializes governance proposal router
func initGovRouter(app *MarsApp) govv1beta1.Router {
	govRouter := govv1beta1.NewRouter()

	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler)
	govRouter.AddRoute(paramsproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper))
	govRouter.AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))
	govRouter.AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))
	govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, getEnabledProposals()))

	return govRouter
}

// initIBCRouter initialzies IBC router.
//
// NOTE: We cannot wrap modules in the fee middleware yet until channel
// upgradability is implemented. See discussion here:
// https://discord.com/channels/955868717269516318/955877042883285023/1062113420712882278
func initIBCRouter(app *MarsApp) *ibcporttypes.Router {
	var icaControllerStack ibcporttypes.IBCModule
	icaControllerStack = envoy.NewIBCModule(app.EnvoyKeeper)
	icaControllerStack = icacontroller.NewIBCMiddleware(icaControllerStack, app.ICAControllerKeeper)

	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, ibctransfer.NewIBCModule(app.IBCTransferKeeper))
	ibcRouter.AddRoute(icacontrollertypes.SubModuleName, icaControllerStack)
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icahost.NewIBCModule(app.ICAHostKeeper))
	ibcRouter.AddRoute(wasm.ModuleName, wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper))

	return ibcRouter
}

// GetEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
// produce a list of enabled proposals to pass into wasmd app.
func GetEnabledProposals() []wasm.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasm.EnableAllProposals
		}
		return wasm.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificProposals, ",")
	proposals, err := wasm.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *MarsApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

func (app *MarsApp) setAnteHandler(txConfig client.TxConfig, wasmConfig wasmtypes.WasmConfig, txCounterStoreKey storetypes.StoreKey) {
	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				SignModeHandler: txConfig.SignModeHandler(),
				FeegrantKeeper:  app.FeeGrantKeeper,
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			IBCKeeper:         app.IBCKeeper,
			WasmConfig:        &wasmConfig,
			TXCounterStoreKey: txCounterStoreKey,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}
	app.SetAnteHandler(anteHandler)
}

func (app *MarsApp) setPostHandler() {
	postHandler, err := posthandler.NewPostHandler(
		posthandler.HandlerOptions{},
	)
	if err != nil {
		panic(err)
	}

	app.SetPostHandler(postHandler)
}
