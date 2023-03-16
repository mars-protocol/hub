package keeper_test

import (
	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"

	wasm "github.com/CosmWasm/wasmd/x/wasm"

	marsutils "github.com/mars-protocol/hub/v2/utils"
	"github.com/mars-protocol/hub/v2/x/envoy/keeper"
	"github.com/mars-protocol/hub/v2/x/envoy/types"
)

var (
	// initial coin balance of the envoy module
	envoyInitBalance = sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(200)))

	// initial balance of the community pool
	communityPoolInitBalance = sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(300)), sdk.NewCoin("uastro", sdk.NewInt(500)))

	// some random messages for use in testing SendMessages
	mockMessages = []proto.Message{
		&govv1.MsgVote{
			Voter:      "interchainAccountAddress",
			ProposalId: 69420,
			Option:     govv1.OptionNoWithVeto,
			Metadata:   "lol",
		},
		&wasm.MsgExecuteContract{
			Sender:   "interchainAccountAddress",
			Contract: "contractAddress",
			Msg:      []byte(`{"detonate_nuclear_bomb":{}}`),
			Funds:    sdk.NewCoins(),
		},
	}
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
			events := ctx.EventManager().Events()

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Len(events, 2)
				suite.Require().Equal(ibcchanneltypes.EventTypeChannelOpenInit, events[0].Type)
				suite.Require().Equal(sdk.EventTypeMessage, events[1].Type)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
				suite.Require().Len(events, 0)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSendFunds() {
	testCases := []struct {
		name      string
		authority string
		channelID string
		amount    sdk.Coins
		expPass   bool
	}{
		{
			"success - empty amount",
			authority.String(),
			suite.path1.EndpointA.ChannelID,
			sdk.NewCoins(),
			true,
		},
		{
			"success - envoy module has a sufficient balance",
			authority.String(),
			suite.path1.EndpointA.ChannelID,
			sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(123))),
			true,
		},
		{
			"success - envoy module does not have a sufficient balance, but community pool does",
			authority.String(),
			suite.path1.EndpointA.ChannelID,
			sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(250)), sdk.NewCoin("uastro", sdk.NewInt(69))),
			true,
		},
		{
			"fail - even community pool does not have a sufficient balance",
			authority.String(),
			suite.path1.EndpointA.ChannelID,
			sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(88888888))),
			false,
		},
		{
			"fail - sender is not authority",
			sender,
			suite.path1.EndpointA.ChannelID,
			sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(123))),
			false,
		},
		{
			"fail - non-existent channel",
			authority.String(),
			"channel-42069",
			sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(123))),
			false,
		},
		{
			"fail - no interchain account found on the connection",
			authority.String(),
			suite.path2.EndpointA.ChannelID, // to outpost2, which doesn't have an ICA registered
			sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(123))),
			false,
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
				suite.Require().NotZero(events)

				envoyBalance := app.BankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(types.ModuleName))
				expEnvoyBalance := marsutils.SaturateSub(envoyInitBalance, tc.amount)
				suite.Require().True(envoyBalance.IsEqual(expEnvoyBalance))

				communityPoolBalance := app.BankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(distrtypes.ModuleName))
				expCommunityPoolBalance := marsutils.SaturateSub(communityPoolInitBalance, marsutils.SaturateSub(tc.amount, envoyInitBalance))
				suite.Require().True(communityPoolBalance.IsEqual(expCommunityPoolBalance))
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
				suite.Require().Len(events, 0)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSendMessages() {
	testCases := []struct {
		name         string
		authority    string
		connectionID string
		messages     []proto.Message
		expPass      bool
	}{
		{
			"success",
			authority.String(),
			suite.path1.EndpointA.ConnectionID,
			mockMessages,
			true,
		},
		{
			"fail - no message",
			authority.String(),
			suite.path1.EndpointA.ConnectionID,
			[]proto.Message{},
			false,
		},
		{
			"fail - sender is not authority",
			sender,
			suite.path1.EndpointA.ConnectionID,
			mockMessages,
			false,
		},
		{
			"fail - non-existent connection",
			authority.String(),
			"connection-88888",
			mockMessages,
			false,
		},
		{
			"fail - no interchain account found on the connection",
			authority.String(),
			suite.path2.EndpointA.ConnectionID, // to outpost2, which doesn't have an ICA registered
			mockMessages,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			// register an interchain account on outpost 1
			registerInterchainAccount(suite.path1, owner.String())

			ctx := suite.hub.GetContext()
			app := getMarsApp(suite.hub)
			msgServer := keeper.NewMsgServerImpl(app.EnvoyKeeper)

			anys := []*codectypes.Any{}
			for _, protoMsg := range tc.messages {
				any, err := codectypes.NewAnyWithValue(protoMsg)
				suite.Require().NoError(err)

				anys = append(anys, any)
			}

			msg := &types.MsgSendMessages{
				Authority:    tc.authority,
				ConnectionId: tc.connectionID,
				Messages:     anys,
			}
			res, err := msgServer.SendMessages(sdk.WrapSDKContext(ctx), msg)
			events := ctx.EventManager().Events()

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Len(events, 2)
				suite.Require().Equal(ibcchanneltypes.EventTypeSendPacket, events[0].Type)
				suite.Require().Equal(sdk.EventTypeMessage, events[1].Type)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
				suite.Require().Len(events, 0)
			}
		})
	}
}
