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

	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"
	"github.com/mars-protocol/hub/x/envoy/types"
)

var (
	// sender is the account that will send permissionless messages,
	// i.e. MsgRegisterAccount, in these tests. Here we use a random address.
	sender = "mars1z926ax906k0ycsuckele6x5hh66e2m4m09whw6"

	// authority is the account who can send privileged messages,
	// i.e. MsgSendFunds and MsgSendMessages. Here we use the gov module account.
	authority = authtypes.NewModuleAddress(govtypes.ModuleName)
)

func init() {
	ibctesting.DefaultTestingAppInit = SetupEnvoyTestingApp
}

func SetupEnvoyTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	encCfg := marsapp.MakeEncodingConfig()
	app := marsapptesting.MakeSimpleMockApp()
	return app, marsapp.DefaultGenesisState(encCfg.Codec)
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

func (suite *KeeperTestSuite) getMarsApp() *marsapp.MarsApp {
	app, ok := suite.hub.App.(*marsapp.MarsApp)
	if !ok {
		panic("not a MarsApp")
	}

	return app
}

func (suite *KeeperTestSuite) setTokenBalances(envoy, communityPool sdk.Coins) {
	ctx := suite.hub.GetContext()
	app := suite.getMarsApp()

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
		FeePool: distrtypes.FeePool{
			CommunityPool: sdk.NewDecCoinsFromCoins(communityPool...),
		},
	})
}

func newTransferPath(hub, outpost *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(hub, outpost)

	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointA.ChannelConfig.Order = ibcchanneltypes.UNORDERED
	path.EndpointB.ChannelConfig.Order = ibcchanneltypes.UNORDERED
	path.EndpointA.ChannelConfig.Version = ibctransfertypes.Version

	return path
}
