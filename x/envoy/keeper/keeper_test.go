package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"
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

// TestKeeperTestSuite runs all the tests within this package.
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {
	// create 3 chains - one hub chain, two outpost chains
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 3)

	suite.hub = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.outpost1 = suite.coordinator.GetChain(ibctesting.GetChainID(2))
	suite.outpost2 = suite.coordinator.GetChain(ibctesting.GetChainID(3))

	suite.path1 = newICAPath(suite.hub, suite.outpost1)
	suite.path2 = newICAPath(suite.hub, suite.outpost2)

	suite.coordinator.SetupConnections(suite.path1)
	suite.coordinator.SetupConnections(suite.path2)
}

func newICAPath(hub, outpost *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(hub, outpost)

	path.EndpointA.ChannelConfig.PortID = icatypes.HostPortID
	path.EndpointB.ChannelConfig.PortID = icatypes.HostPortID
	path.EndpointA.ChannelConfig.Order = ibcchanneltypes.ORDERED
	path.EndpointB.ChannelConfig.Order = ibcchanneltypes.ORDERED

	return path
}

func getMarsApp(chain *ibctesting.TestChain) *marsapp.MarsApp {
	app, ok := chain.App.(*marsapp.MarsApp)
	if !ok {
		panic("not mars app")
	}

	return app
}
