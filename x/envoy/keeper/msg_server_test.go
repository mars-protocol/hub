package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"

	"github.com/mars-protocol/hub/x/envoy/keeper"
	"github.com/mars-protocol/hub/x/envoy/types"
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
				app := getMarsApp(suite.hub)

				_, portID, err := app.EnvoyKeeper.GetOwnerAndPortID()
				suite.Require().NoError(err)

				app.IBCKeeper.PortKeeper.BindPort(suite.hub.GetContext(), portID)
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
				app := getMarsApp(suite.hub)

				_, portID, _ := app.EnvoyKeeper.GetOwnerAndPortID()

				channel := ibcchanneltypes.NewChannel(
					ibcchanneltypes.OPEN,
					ibcchanneltypes.ORDERED,
					ibcchanneltypes.NewCounterparty(suite.path1.EndpointB.ChannelConfig.PortID, suite.path1.EndpointB.ChannelID),
					[]string{suite.path1.EndpointA.ConnectionID},
					suite.path1.EndpointA.ChannelConfig.Version,
				)

				app.IBCKeeper.ChannelKeeper.SetChannel(suite.hub.GetContext(), portID, ibctesting.FirstChannelID, channel)
				app.ICAControllerKeeper.SetActiveChannelID(suite.hub.GetContext(), ibctesting.FirstConnectionID, portID, ibctesting.FirstChannelID)
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
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}
