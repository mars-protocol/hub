package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/incentives/types"
)

const govModuleAccount = "mars10d07y265gmmuvt4z0w9aw880jnsr700j8l2urg"

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount("mars", "marspub")
}

func getMockMsgCreateSchedule() types.MsgCreateSchedule {
	return types.MsgCreateSchedule{
		Authority: govModuleAccount,
		StartTime: time.Unix(10000, 0),
		EndTime:   time.Unix(20000, 0),
		Amount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(10000))),
	}
}

func getMockTerminateScheduleProposal() types.MsgTerminateSchedules {
	return types.MsgTerminateSchedules{
		Authority: govModuleAccount,
		Ids:       []uint64{1, 2, 3, 4, 5},
	}
}

func TestValidateCreateScheduleProposal(t *testing.T) {
	p := getMockMsgCreateSchedule()
	p.EndTime = time.Unix(9999, 0)
	require.Error(t, p.ValidateBasic(), types.ErrInvalidProposalStartEndTimes)

	p = getMockMsgCreateSchedule()
	p.Amount = sdk.NewCoins()
	require.Error(t, p.ValidateBasic(), types.ErrInvalidProposalAmount)

	p = getMockMsgCreateSchedule()
	require.NoError(t, p.ValidateBasic())
}

func TestValidateTerminateScheduleProposal(t *testing.T) {
	p := getMockTerminateScheduleProposal()
	p.Ids = []uint64{}
	require.Error(t, p.ValidateBasic(), types.ErrInvalidProposalIds)

	p = getMockTerminateScheduleProposal()
	require.NoError(t, p.ValidateBasic())
}
