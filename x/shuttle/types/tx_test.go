package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

var (
	testAuthority    = authtypes.NewModuleAddress(types.ModuleName)
	testSender, _    = sdk.AccAddressFromBech32("cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs")
	testConnectionId = "connection_0"
	testChannelId    = "channel_0"
	testValidMsg, _  = codectypes.NewAnyWithValue(govv1.NewMsgVote(testAuthority, 1, govv1.OptionYes, ""))
	testInvalidMsg   = &codectypes.Any{TypeUrl: "/test.MsgInvalidTest", Value: []byte{}}
)

func TestValidateBasic(t *testing.T) {
	testCases := []struct {
		name    string
		msg     sdk.Msg
		expPass bool
	}{
		{
			"MsgRegisterAccount - success",
			&types.MsgRegisterAccount{
				Sender:       testSender.String(),
				ConnectionId: testConnectionId,
			},
			true,
		},
		{
			"MsgRegisterAccount - owner address is empty",
			&types.MsgRegisterAccount{
				Sender:       "",
				ConnectionId: testConnectionId,
			},
			false,
		},
		{
			"MsgRegisterAccount - owner address is invalid",
			&types.MsgRegisterAccount{
				Sender:       "larry",
				ConnectionId: testConnectionId,
			},
			false,
		},
		{
			"MsgSendFunds - success",
			&types.MsgSendFunds{
				Authority: testAuthority.String(),
				ChannelId: testChannelId,
				Amount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(12345))),
			},
			true,
		},
		{
			"MsgSendFunds - coin amount is empty",
			&types.MsgSendFunds{
				Authority: testAuthority.String(),
				ChannelId: testChannelId,
				Amount:    sdk.NewCoins(),
			},
			false,
		},
		{
			"MsgSendFunds - coin amount is invalid",
			&types.MsgSendFunds{
				Authority: testAuthority.String(),
				ChannelId: testChannelId,
				// denoms not sorted alphabetically
				Amount: []sdk.Coin{
					sdk.NewCoin("umars", sdk.NewInt(12345)),
					sdk.NewCoin("uastro", sdk.NewInt(23456)),
				},
			},
			false,
		},
		{
			"MsgSendMessages - success",
			&types.MsgSendMessages{
				Authority:    testAuthority.String(),
				ConnectionId: testConnectionId,
				Messages:     []*codectypes.Any{testValidMsg},
			},
			true,
		},
		{
			"MsgSendMessages - messages is empty",
			&types.MsgSendMessages{
				Authority:    testAuthority.String(),
				ConnectionId: testConnectionId,
				Messages:     []*codectypes.Any{},
			},
			false,
		},
		{
			"MsgSendMessages - message does not implement sdk.Msg interface",
			&types.MsgSendMessages{
				Authority:    testAuthority.String(),
				ConnectionId: testConnectionId,
				Messages:     []*codectypes.Any{testInvalidMsg},
			},
			false,
		},
	}

	for _, tc := range testCases {
		err := tc.msg.ValidateBasic()

		if tc.expPass {
			require.NoError(t, err, "expect success but failed: name = %s", tc.name)
		} else {
			require.Error(t, err, "expect error but succeeded: name = %s", tc.name)
		}
	}
}

func TestGetSigners(t *testing.T) {
	testCases := []struct {
		name      string
		msg       sdk.Msg
		expSigner sdk.AccAddress
	}{
		{
			"MsgRegisterAccount",
			&types.MsgRegisterAccount{
				Sender:       testSender.String(),
				ConnectionId: testConnectionId,
			},
			testSender,
		},
		{
			"MsgSendFunds",
			&types.MsgSendFunds{
				Authority: testAuthority.String(),
				ChannelId: testChannelId,
				Amount:    sdk.NewCoins(),
			},
			testAuthority,
		},
		{
			"MsgSendMessages",
			&types.MsgSendMessages{
				Authority:    testAuthority.String(),
				ConnectionId: testConnectionId,
				Messages:     []*codectypes.Any{},
			},
			testAuthority,
		},
	}

	for _, tc := range testCases {
		require.Equal(t, []sdk.AccAddress{tc.expSigner}, tc.msg.GetSigners(), "incorrect sender: name = %s", tc.name)
	}
}
