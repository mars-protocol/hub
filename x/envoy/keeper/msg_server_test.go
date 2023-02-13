package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"

	"github.com/mars-protocol/hub/x/envoy/keeper"
	"github.com/mars-protocol/hub/x/envoy/types"
)

var (
	envoyInitBalance         = sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(200)))
	communityPoolInitBalance = sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(300)), sdk.NewCoin("uastro", sdk.NewInt(500)))
)

func (suite *KeeperTestSuite) TestRegisterAccount() {
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success - first account",
			func() {},
			true,
		},

		// in this test case, we make sure that envoy module can own interchain
		// accounts on multiple chains with the same portID.
		//
		// specifically, we register an account on outpost 2 first, then attempt
		// to register on outpost 1.
		{
			"success - second account",
			func() {
				ctx := suite.hub.GetContext()
				app := getMarsApp(suite.hub)
				msgServer := keeper.NewMsgServerImpl(app.EnvoyKeeper)

				// register an account on outpost 2
				msg := &types.MsgRegisterAccount{
					Sender:       sender,
					ConnectionId: suite.path2.EndpointA.ConnectionID,
				}
				res, err := msgServer.RegisterAccount(sdk.WrapSDKContext(ctx), msg)

				// registring the account on outpost 2 should be successful
				suite.Require().NoError(err)
				suite.Require().NotNil(res)

				// next, registering an account on outpost 1 should be successful as well
			},
			true,
		},

		// this test case refers to the case where a module other than the ICA
		// controller module has claimed the capability of the ICA controller port.
		//
		// this should result in the following error:
		//
		// ```plain
		// another module has claimed capability for and bound port with portID:
		// icacontroller-cosmos1s3fjkvr0yk2c0smyh4esrcyp893atwz0p4yr2j: port is already bound
		// ```
		{
			"failure - port is already bound",
			func() {
				ctx := suite.hub.GetContext()
				app := getMarsApp(suite.hub)

				_, portID, err := app.EnvoyKeeper.GetOwnerAndPortID()
				suite.Require().NoError(err)

				app.IBCKeeper.PortKeeper.BindPort(ctx, portID)
			},
			false,
		},

		// this test case refers to the case where there is already an interchain
		// account for the chosen connection.
		//
		// this should result in the following error:
		//
		// ```plain
		// existing active channel channel-0 for portID
		// icacontroller-cosmos1s3fjkvr0yk2c0smyh4esrcyp893atwz0p4yr2j
		// on connection connection-0: active channel already set for this owner
		// ```
		{
			"failure - channel is already active",
			func() {
				ctx := suite.hub.GetContext()
				app := getMarsApp(suite.hub)

				_, portID, _ := app.EnvoyKeeper.GetOwnerAndPortID()

				channel := ibcchanneltypes.NewChannel(
					ibcchanneltypes.OPEN,
					ibcchanneltypes.ORDERED,
					ibcchanneltypes.NewCounterparty(suite.path1.EndpointB.ChannelConfig.PortID, suite.path1.EndpointB.ChannelID),
					[]string{suite.path1.EndpointA.ConnectionID},
					suite.path1.EndpointA.ChannelConfig.Version,
				)

				app.IBCKeeper.ChannelKeeper.SetChannel(ctx, portID, ibctesting.FirstChannelID, channel)
				app.ICAControllerKeeper.SetActiveChannelID(ctx, ibctesting.FirstConnectionID, portID, ibctesting.FirstChannelID)
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			// IMPORTANT: must do GetContext *after* SetupTest
			ctx := suite.hub.GetContext()
			app := getMarsApp(suite.hub)
			msgServer := keeper.NewMsgServerImpl(app.EnvoyKeeper)

			// mallete mutates test data
			tc.malleate()

			msg := &types.MsgRegisterAccount{
				Sender:       sender,
				ConnectionId: suite.path1.EndpointA.ConnectionID, // should be the connection id on the hub chain
			}
			res, err := msgServer.RegisterAccount(sdk.WrapSDKContext(ctx), msg)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				// TODO: verify the channel is properly registered
				// TODO: verify events
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSendFunds() {
	testCases := []struct {
		name                   string
		authority              string
		channelID              string
		amount                 sdk.Coins
		expPass                bool
		expEnvoyRemain         sdk.Coins
		expCommunityPoolRemain sdk.Coins
	}{
		{
			"success - envoy module has enough coins",
			authority.String(),
			suite.path1.EndpointA.ChannelID, // to outpost1
			sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(123))),
			true,
			sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(200-123))),
			communityPoolInitBalance,
		},

		// success - envoy module does not have enough coins

		// fail - commmunity pool does not have enough coins

		// fail - not authority

		// fail - amount is empty

		// fail - channel not found

		{
			"fail - no interchain account found on the connection",
			authority.String(),
			suite.path2.EndpointA.ChannelID, // to outpost2, which doesn't have an ICA registered
			sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(123))),
			false,
			sdk.NewCoins(),
			sdk.NewCoins(),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			// set initial balances for envoy module and community pool
			setTokenBalances(suite.hub, envoyInitBalance, communityPoolInitBalance)

			// register an interchain account on outpost 1
			registerInterchainAccount(suite.path1, owner.String())

			ctx := suite.hub.GetContext()
			app := getMarsApp(suite.hub)
			msgServer := keeper.NewMsgServerImpl(app.EnvoyKeeper)

			msg := &types.MsgSendFunds{
				Authority: tc.authority,
				ChannelId: tc.channelID,
				Amount:    tc.amount,
			}
			res, err := msgServer.SendFunds(sdk.WrapSDKContext(ctx), msg)
			events := ctx.EventManager().Events()

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				// TODO: verify events
				// TODO: verify token balances after msg execution
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
				suite.Require().Len(events, 0)
			}
		})
	}
}
