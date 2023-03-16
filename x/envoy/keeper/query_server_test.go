package keeper_test

import (
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"

	"github.com/mars-protocol/hub/v2/x/envoy/keeper"
	"github.com/mars-protocol/hub/v2/x/envoy/types"
)

func (suite *KeeperTestSuite) TestQueryAccount() {
	suite.SetupTest()

	registerInterchainAccount(suite.path1, owner.String())

	ctx := suite.hub.GetContext()
	app := getMarsApp(suite.hub)
	queryServer := keeper.NewQueryServerImpl(app.EnvoyKeeper)

	// query account at outpost 1 - should succeed
	res, err := queryServer.Account(ctx, &types.QueryAccountRequest{
		ConnectionId: suite.path1.EndpointA.ConnectionID,
	})
	suite.Require().NoError(err)

	// set address as empty since we don't care about its particular value
	res.Account.Address = ""
	suite.Require().Equal(composeAccountInfoFromPath(suite.path1), res.Account)

	// query account at outpost 2 - should fail
	res, err = queryServer.Account(ctx, &types.QueryAccountRequest{
		ConnectionId: suite.path2.EndpointA.ConnectionID,
	})
	suite.Require().Error(err)
	suite.Require().Nil(res)
}

func (suite *KeeperTestSuite) TestQueryAccounts() {
	suite.SetupTest()

	registerInterchainAccount(suite.path1, owner.String())
	registerInterchainAccount(suite.path2, owner.String())

	ctx := suite.hub.GetContext()
	app := getMarsApp(suite.hub)
	queryServer := keeper.NewQueryServerImpl(app.EnvoyKeeper)

	res, err := queryServer.Accounts(ctx, &types.QueryAccountsRequest{})
	suite.Require().NoError(err)

	res.Accounts[0].Address = ""
	suite.Require().Equal(composeAccountInfoFromPath(suite.path1), res.Accounts[0])

	res.Accounts[1].Address = ""
	suite.Require().Equal(composeAccountInfoFromPath(suite.path2), res.Accounts[1])
}

func composeAccountInfoFromPath(path *ibctesting.Path) *types.AccountInfo {
	return &types.AccountInfo{
		Controller: &types.ChainInfo{
			ClientId:     path.EndpointA.ClientID,
			ConnectionId: path.EndpointA.ConnectionID,
			PortId:       portID,
			ChannelId:    path.EndpointA.ChannelID,
		},
		Host: &types.ChainInfo{
			ClientId:     path.EndpointB.ClientID,
			ConnectionId: path.EndpointB.ConnectionID,
			PortId:       icatypes.HostPortID,
			ChannelId:    path.EndpointB.ChannelID,
		},
	}
}
