package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"
	"github.com/mars-protocol/hub/x/envoy/types"
)

var (
	// owner of the interchain account; should be the envoy module address
	owner = authtypes.NewModuleAddress(types.ModuleName)

	// sender is the account that will send permissionless messages,
	// i.e. MsgRegisterAccount, in these tests. Here we use a random address.
	sender = "mars1z926ax906k0ycsuckele6x5hh66e2m4m09whw6"

	// authority is the account who can send privileged messages,
	// i.e. MsgSendFunds and MsgSendMessages. Here we use the gov module account.
	authority = authtypes.NewModuleAddress(govtypes.ModuleName)
)

func init() {
	ibctesting.DefaultTestingAppInit = func() (ibctesting.TestingApp, map[string]json.RawMessage) {
		encCfg := marsapp.MakeEncodingConfig()
		app := marsapptesting.MakeSimpleMockApp()
		return app, marsapp.DefaultGenesisState(encCfg.Codec)
	}
}

// KeeperTestSuite is a testing suite
type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	hub      *ibctesting.TestChain
	outpost1 *ibctesting.TestChain
	outpost2 *ibctesting.TestChain
	path1    *ibctesting.Path
	path2    *ibctesting.Path
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	// create 3 chains - one hub chain, two outpost chains
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 3)

	suite.hub = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.outpost1 = suite.coordinator.GetChain(ibctesting.GetChainID(2))
	suite.outpost2 = suite.coordinator.GetChain(ibctesting.GetChainID(3))

	suite.path1 = newTransferPath(suite.hub, suite.outpost1)
	suite.path2 = newTransferPath(suite.hub, suite.outpost2)

	suite.coordinator.SetupConnections(suite.path1)
	suite.coordinator.SetupConnections(suite.path2)

	suite.coordinator.CreateTransferChannels(suite.path1)
	suite.coordinator.CreateTransferChannels(suite.path2)
}

func setTokenBalances(chain *ibctesting.TestChain, envoy, communityPool sdk.Coins) {
	ctx := chain.GetContext()
	app := getMarsApp(chain)

	distrAddr := authtypes.NewModuleAddress(distrtypes.ModuleName)
	envoyAddr := authtypes.NewModuleAddress(types.ModuleName)

	app.BankKeeper.InitGenesis(ctx, &banktypes.GenesisState{
		Params: banktypes.Params{
			DefaultSendEnabled: true,
		},
		Balances: []banktypes.Balance{
			{
				Address: distrAddr.String(),
				Coins:   communityPool,
			},
			{
				Address: envoyAddr.String(),
				Coins:   communityPool,
			},
		},
	})

	app.DistrKeeper.InitGenesis(ctx, distrtypes.GenesisState{
		Params: distrtypes.DefaultParams(),
		FeePool: distrtypes.FeePool{
			CommunityPool: sdk.NewDecCoinsFromCoins(communityPool...),
		},
	})
}

func registerInterchainAccount(path *ibctesting.Path, owner string) {
	controllerCtx := path.EndpointA.Chain.GetContext()
	controllerApp := getMarsApp(path.EndpointA.Chain)

	channelSequence := controllerApp.IBCKeeper.ChannelKeeper.GetNextChannelSequence(controllerCtx)

	if err := controllerApp.ICAControllerKeeper.RegisterInterchainAccount(controllerCtx, path.EndpointA.ConnectionID, owner, ""); err != nil {
		panic(err)
	}

	// commit state changes for proof verification
	path.EndpointA.Chain.NextBlock()

	// create ICA path, update channel ids
	icaPath := updateToICAPath(path, owner)
	icaPath.EndpointA.ChannelID = ibcchanneltypes.FormatChannelIdentifier(channelSequence)
	icaPath.EndpointB.ChannelID = ""

	if err := icaPath.EndpointB.ChanOpenTry(); err != nil {
		panic(err)
	}

	if err := icaPath.EndpointA.ChanOpenAck(); err != nil {
		panic(err)
	}

	if err := icaPath.EndpointB.ChanOpenConfirm(); err != nil {
		panic(err)
	}
}

func getMarsApp(chain *ibctesting.TestChain) *marsapp.MarsApp {
	app, ok := chain.App.(*marsapp.MarsApp)
	if !ok {
		panic("not a MarsApp")
	}

	return app
}

func newTransferPath(hub, outpost *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(hub, outpost)

	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointA.ChannelConfig.Order = ibcchanneltypes.UNORDERED
	path.EndpointB.ChannelConfig.Order = ibcchanneltypes.UNORDERED
	path.EndpointA.ChannelConfig.Version = ibctransfertypes.Version
	path.EndpointB.ChannelConfig.Version = ibctransfertypes.Version

	return path
}

func updateToICAPath(path *ibctesting.Path, owner string) *ibctesting.Path {
	controllerPortID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		panic(err)
	}

	version := icatypes.NewDefaultMetadataString(path.EndpointA.ConnectionID, path.EndpointB.ConnectionID)

	path.EndpointA.ChannelConfig.PortID = controllerPortID
	path.EndpointB.ChannelConfig.PortID = icatypes.HostPortID
	path.EndpointA.ChannelConfig.Order = ibcchanneltypes.ORDERED
	path.EndpointB.ChannelConfig.Order = ibcchanneltypes.ORDERED
	path.EndpointA.ChannelConfig.Version = version
	path.EndpointB.ChannelConfig.Version = version

	return path
}
