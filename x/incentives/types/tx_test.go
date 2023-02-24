package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/incentives/types"
)

const govModuleAccount = "mars10d07y265gmmuvt4z0w9aw880jnsr700j8l2urg"

var (
	mockMsgCreateSchedule = types.MsgCreateSchedule{
		Authority: govModuleAccount,
		StartTime: time.Unix(10000, 0),
		EndTime:   time.Unix(20000, 0),
		Amount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(10000))),
	}

	mockMsgTerminateSchedules = types.MsgTerminateSchedules{
		Authority: govModuleAccount,
		Ids:       []uint64{1, 2, 3, 4, 5},
	}
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount("mars", "marspub")
}

func TestValidateCreateScheduleProposal(t *testing.T) {
	var msg types.MsgCreateSchedule

	testCases := []struct {
		name     string
		malleate func()
		expError error
	}{
		{
			"succeed",
			func() {},
			nil,
		},
		{
			"fail - end time is earlier than start time",
			func() {
				msg.EndTime = time.Unix(9999, 0)
			},
			types.ErrInvalidProposalStartEndTimes,
		},
		{
			"fail - amount is empty",
			func() {
				msg.Amount = sdk.NewCoins()
			},
			types.ErrInvalidProposalAmount,
		},
		{
			"fail - amount contains zero coin",
			func() {
				msg.Amount = []sdk.Coin{{Denom: "umars", Amount: sdk.NewInt(0)}}
			},
			types.ErrInvalidProposalAmount,
		},
		{
			"fail - amount contains negative coin",
			func() {
				msg.Amount = []sdk.Coin{{Denom: "umars", Amount: sdk.NewInt(-1)}}
			},
			types.ErrInvalidProposalAmount,
		},
		{
			"fail - coins are out of order",
			func() {
				msg.Amount = []sdk.Coin{sdk.NewCoin("uastro", sdk.NewInt(12345)), sdk.NewCoin("uastro", sdk.NewInt(12345))}
			},
			types.ErrInvalidProposalAmount,
		},
	}

	for _, tc := range testCases {
		msg = mockMsgCreateSchedule
		tc.malleate()

		if tc.expError != nil {
			require.Error(t, msg.ValidateBasic(), tc.expError.Error(), tc.name)
		} else {
			require.NoError(t, msg.ValidateBasic(), tc.name)
		}
	}
}

func TestValidateTerminateScheduleProposal(t *testing.T) {
	var msg types.MsgTerminateSchedules

	testCases := []struct {
		name     string
		malleate func()
		expError error
	}{
		{
			"succeed",
			func() {},
			nil,
		},
		{
			"fail - ids are empty",
			func() {
				msg.Ids = []uint64{}
			},
			types.ErrInvalidProposalIds,
		},
	}

	for _, tc := range testCases {
		msg = mockMsgTerminateSchedules
		tc.malleate()

		if tc.expError != nil {
			require.Error(t, msg.ValidateBasic(), tc.expError.Error(), tc.name)
		} else {
			require.NoError(t, msg.ValidateBasic(), tc.name)
		}
	}
}
