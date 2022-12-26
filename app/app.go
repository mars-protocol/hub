package app

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"

	// tendermint
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	// cosmos SDK
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	api "github.com/cosmos/cosmos-sdk/server/api"
	config "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"

	// core modules
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
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
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
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
	customgov "github.com/mars-protocol/hub/x/gov"
	customgovkeeper "github.com/mars-protocol/hub/x/gov/keeper"

	// ibc modules
	ica "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v6/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/types"
	ibctransfer "github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v6/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v6/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v6/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"

	// wasm modules
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"

	// mars modules
	"github.com/mars-protocol/hub/x/incentives"
	incentiveskeeper "github.com/mars-protocol/hub/x/incentives/keeper"
	incentivestypes "github.com/mars-protocol/hub/x/incentives/types"
	"github.com/mars-protocol/hub/x/safety"
	safetykeeper "github.com/mars-protocol/hub/x/safety/keeper"
	safetytypes "github.com/mars-protocol/hub/x/safety/types"

	marswasm "github.com/mars-protocol/hub/app/wasm"
	marsdocs "github.com/mars-protocol/hub/docs"
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
		gov.NewAppModuleBasic(govProposalHandlers), // gov AppModuleBasic is not customized, so we just use the vanilla one
		params.AppModuleBasic{},
		slashing.AppModuleBasic{},
		staking.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibcfee.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		ica.AppModuleBasic{},
		wasm.AppModuleBasic{},
		incentives.AppModuleBasic{},
		safety.AppModuleBasic{},
	)

	// governance proposal handlers
	govProposalHandlers = append(
		wasmclient.ProposalHandlers,
		paramsclient.ProposalHandler,
		distrclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		govtypes.ModuleName:            {authtypes.Burner},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		ibcfeetypes.ModuleName:         nil,
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:            nil,
		wasm.ModuleName:                {authtypes.Burner},
		incentivestypes.ModuleName:     nil,
		safetytypes.ModuleName:         nil,
	}
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
	_ simapp.App              = (*MarsApp)(nil)
	_ servertypes.Application = (*MarsApp)(nil)
)

// MarsApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type MarsApp struct {
	// baseapp
	*baseapp.BaseApp

	// codecs
	legacyAmino       *codec.LegacyAmino
	Codec             codec.Codec // make codec public for testing purposes
	interfaceRegistry codectypes.InterfaceRegistry

	// invariant check period
	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper       authkeeper.AccountKeeper
	AuthzKeeper         authzkeeper.Keeper
	BankKeeper          bankkeeper.Keeper
	CapabilityKeeper    *capabilitykeeper.Keeper
	CrisisKeeper        crisiskeeper.Keeper
	DistrKeeper         distrkeeper.Keeper
	EvidenceKeeper      evidencekeeper.Keeper
	FeeGrantKeeper      feegrantkeeper.Keeper
	GovKeeper           customgovkeeper.Keeper // replaces the vanilla gov keeper with our custom one
	ParamsKeeper        paramskeeper.Keeper
	SlashingKeeper      slashingkeeper.Keeper
	StakingKeeper       stakingkeeper.Keeper
	UpgradeKeeper       upgradekeeper.Keeper
	IBCKeeper           *ibckeeper.Keeper // must be a pointer, so we can `SetRouter` on it correctly
	IBCFeeKeeper        ibcfeekeeper.Keeper
	IBCTransferKeeper   ibctransferkeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper
	WasmKeeper          wasm.Keeper
	IncentivesKeeper    incentiveskeeper.Keeper
	SafetyKeeper        safetykeeper.Keeper

	// make scoped keepers public for testing purposes
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedIBCTransferKeeper   capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper          capabilitykeeper.ScopedKeeper

	// the module manager
	mm *module.Manager
}

// NewMarsApp creates and initializes a new `MarsApp` instance
func NewMarsApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig EncodingConfig, appOpts servertypes.AppOptions,
	wasmOpts []wasm.Option, baseAppOptions ...func(*baseapp.BaseApp),
) *MarsApp {
	legacyAmino := encodingConfig.Amino
	codec := encodingConfig.Marshaler // this field will be renamed to `Codec`
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(
		Name,
		logger,
		db,
		encodingConfig.TxConfig.TxDecoder(),
		baseAppOptions...,
	)
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
		slashingtypes.StoreKey,
		stakingtypes.StoreKey,
		upgradetypes.StoreKey,
		ibchost.StoreKey,
		ibcfeetypes.StoreKey,
		ibctransfertypes.StoreKey,
		icacontrollertypes.StoreKey,
		icahosttypes.StoreKey,
		wasm.StoreKey,
		incentivestypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &MarsApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		Codec:             codec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	// **** create keepers ****

	// address of the gov module account
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	app.ParamsKeeper = initParamsKeeper(
		codec,
		legacyAmino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	bApp.SetParamStore(
		app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()),
	)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		codec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	app.ScopedIBCKeeper = app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	app.ScopedIBCTransferKeeper = app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	app.ScopedICAControllerKeeper = app.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	app.ScopedICAHostKeeper = app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	app.ScopedWasmKeeper = app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		codec,
		keys[authtypes.StoreKey],
		getSubspace(app, authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
		AccountAddressPrefix,
	)
	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		codec,
		app.BaseApp.MsgServiceRouter(),
		app.AccountKeeper,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		codec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		getSubspace(app, banktypes.ModuleName),
		getBlockedModuleAccountAddrs(app), // NOTE: fee collector & safety fund are excluded from blocked addresses
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		getSubspace(app, crisistypes.ModuleName),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)
	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		codec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		codec,
		homePath,
		app.BaseApp,
		authority,
	)

	// staking keeper and its dependencies
	// NOTE: the order here (e.g. evidence keeper depends on slashing keeper, so
	// must be defined after it)
	stakingKeeper := stakingkeeper.NewKeeper(
		codec,
		keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		getSubspace(app, stakingtypes.ModuleName),
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		codec,
		keys[distrtypes.StoreKey],
		getSubspace(app, distrtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		authtypes.FeeCollectorName,
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		codec,
		keys[slashingtypes.StoreKey],
		&stakingKeeper,
		getSubspace(app, slashingtypes.ModuleName),
	)
	app.EvidenceKeeper = *evidencekeeper.NewKeeper(
		codec,
		keys[evidencetypes.StoreKey],
		&app.StakingKeeper,
		app.SlashingKeeper,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain
	// these hooks.
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		codec,
		keys[ibchost.StoreKey],
		getSubspace(app, ibchost.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper, app.ScopedIBCKeeper,
	)
	app.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		codec,
		keys[ibcfeetypes.StoreKey],
		app.IBCKeeper.ChannelKeeper, // default to ChannelKeeper, may be replaced with middleware such as ICS-29
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
	)
	app.IBCTransferKeeper = ibctransferkeeper.NewKeeper(
		codec,
		keys[ibctransfertypes.StoreKey],
		getSubspace(app, ibctransfertypes.ModuleName),
		app.IBCFeeKeeper, // default to ChannelKeeper, may be replaced with middleware such as ICS-29
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		app.ScopedIBCTransferKeeper,
	)
	app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		codec,
		keys[icacontrollertypes.StoreKey],
		getSubspace(app, icacontrollertypes.SubModuleName),
		app.IBCFeeKeeper, // default to ChannelKeeper, may be replaced with middleware such as ICS-29
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.ScopedICAControllerKeeper,
		app.MsgServiceRouter(),
	)
	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		codec,
		keys[icahosttypes.StoreKey],
		getSubspace(app, icahosttypes.SubModuleName),
		app.IBCFeeKeeper, // default to ChannelKeeper, may be replaced with middleware such as ICS-29
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

	supportedFeatures := "iterator,staking,stargate"

	// register wasm bindings of Mars modules here
	wasmOpts = append(marswasm.RegisterCustomPlugins(&app.DistrKeeper), wasmOpts...)

	// create wasm keeper
	app.WasmKeeper = wasm.NewKeeper(
		codec,
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
		supportedFeatures,
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

	// finally, create gov keeper
	//
	// here we use the customized gov keeper, which requires an additional
	// `wasmKeeper` parameter compared to the vanilla govkeeper
	app.GovKeeper = customgovkeeper.NewKeeper(
		codec,
		keys[govtypes.StoreKey],
		getSubspace(app, govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		app.WasmKeeper,
		initGovRouter(app),
		app.MsgServiceRouter(),
		govtypes.DefaultConfig(),
	)

	// **** module options ****

	// NOTE: We may consider parsing `appOpts` inside module constructors. For
	// the moment we prefer to be more strict in what arguments the modules
	// expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later
	// modified must be passed by reference here.
	app.mm = module.NewManager(
		auth.NewAppModule(codec, app.AccountKeeper, nil),
		authzmodule.NewAppModule(codec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		bank.NewAppModule(codec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(codec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		distr.NewAppModule(codec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(codec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
		customgov.NewAppModule(codec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		params.NewAppModule(app.ParamsKeeper),
		slashing.NewAppModule(codec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(codec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		ibcfee.NewAppModule(app.IBCFeeKeeper),
		ibctransfer.NewAppModule(app.IBCTransferKeeper),
		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
		wasm.NewAppModule(codec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		incentives.NewAppModule(app.IncentivesKeeper),
		safety.NewAppModule(app.SafetyKeeper),
	)

	// During begin block, slashing happens after `distr.BeginBlocker` so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// `CanWithdrawInvariant` invariant.
	// NOTE: staking module is required if `HistoricalEntries` param > 0
	app.mm.SetOrderBeginBlockers(
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
		ibchost.ModuleName,
		ibcfeetypes.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		wasm.ModuleName,
		incentivestypes.ModuleName,
		safetytypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
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
		ibchost.ModuleName,
		ibcfeetypes.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		wasm.ModuleName,
		incentivestypes.ModuleName,
		safetytypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any
	// capabilities so that other modules that want to create or claim
	// capabilities afterwards in `InitChain` can do so safely.
	app.mm.SetOrderInitGenesis(
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
		ibchost.ModuleName,
		ibcfeetypes.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		wasm.ModuleName,
		incentivestypes.ModuleName,
		safetytypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.mm.RegisterServices(module.NewConfigurator(codec, app.MsgServiceRouter(), app.GRPCQueryRouter()))

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeeGrantKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			IBCKeeper:         app.IBCKeeper,
			WasmConfig:        wasmConfig,
			TxCounterStoreKey: keys[wasm.StoreKey],
		},
	)
	if err != nil {
		panic(err)
	}

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

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
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates on every end block
func (app *MarsApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application updates at chain initialization
func (app *MarsApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.Codec, genesisState)
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
// Further note on the 2nd point: the distrkeeper's `DistributeFromFeePool`
// function uses bankkeeper's `SendCoinsFromModuleToAccount` instead of
// `SendCoinsFromModuleToModule`. If it had used `FromModuleToModule`
// then we won't need to allow incentives module account to receive funds here.
//
// Forked from: https://github.com/cosmos/gaia/pull/1493
func getBlockedModuleAccountAddrs(app *MarsApp) map[string]bool {
	modAccAddrs := app.ModuleAccountAddrs()

	delete(modAccAddrs, authtypes.NewModuleAddress(authtypes.FeeCollectorName).String())
	delete(modAccAddrs, authtypes.NewModuleAddress(incentivestypes.ModuleName).String())
	delete(modAccAddrs, authtypes.NewModuleAddress(safetytypes.ModuleName).String())

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
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(wasm.ModuleName)

	return paramsKeeper
}

// initializes governance proposal router
func initGovRouter(app *MarsApp) govv1beta1.Router {
	govRouter := govv1beta1.NewRouter()

	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler)
	govRouter.AddRoute(paramsproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper))
	govRouter.AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper))
	govRouter.AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))
	govRouter.AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))
	govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, getEnabledProposals()))

	return govRouter
}

// initialzies IBC router
func initIBCRouter(app *MarsApp) *ibcporttypes.Router {
	// wrap transfer module inside fee middleware
	var transferStack ibcporttypes.IBCModule
	transferStack = ibctransfer.NewIBCModule(app.IBCTransferKeeper)
	transferStack = ibcfee.NewIBCMiddleware(transferStack, app.IBCFeeKeeper)

	// wrap shuttle module in ICA controller middleware then in fee middleware
	var icaControllerStack ibcporttypes.IBCModule
	icaControllerStack = icacontroller.NewIBCMiddleware(nil, app.ICAControllerKeeper) // TODO: replace nil with shuttle module
	icaControllerStack = ibcfee.NewIBCMiddleware(icaControllerStack, app.IBCFeeKeeper)

	// wrap ICA host module inside fee middleware
	var icaHostStack ibcporttypes.IBCModule
	icaHostStack = icahost.NewIBCModule(app.ICAHostKeeper)
	icaHostStack = ibcfee.NewIBCMiddleware(icaHostStack, app.IBCFeeKeeper)

	// add modules to the router
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)
	ibcRouter.AddRoute(icacontrollertypes.SubModuleName, icaControllerStack)
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostStack)

	return ibcRouter
}
