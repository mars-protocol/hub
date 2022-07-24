package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/incentives/types"
)

func getMockCreateScheduleProposal() types.CreateIncentivesScheduleProposal {
	return types.CreateIncentivesScheduleProposal{
		Title:       "title",
		Description: "description",
		StartTime:   time.Unix(10000, 0),
		EndTime:     time.Unix(20000, 0),
		Amount:      sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(10000))),
	}
}

func getMockTerminateScheduleProposal() types.TerminateIncentivesScheduleProposal {
	return types.TerminateIncentivesScheduleProposal{
		Title:       "title",
		Description: "description",
		Ids:         []uint64{1, 2, 3, 4, 5},
	}
}

func TestValidateCreateScheduleProposal(t *testing.T) {
	p := getMockCreateScheduleProposal()
	p.EndTime = time.Unix(9999, 0)
	require.Error(t, p.ValidateBasic(), types.ErrInvalidProposalStartEndTimes)

	p = getMockCreateScheduleProposal()
	p.Amount = sdk.NewCoins()
	require.Error(t, p.ValidateBasic(), types.ErrInvalidProposalAmount)

	p = getMockCreateScheduleProposal()
	require.NoError(t, p.ValidateBasic())
}

func TestValidateTerminateScheduleProposal(t *testing.T) {
	p := getMockTerminateScheduleProposal()
	p.Ids = []uint64{}
	require.Error(t, p.ValidateBasic(), types.ErrInvalidProposalIds)

	p = getMockTerminateScheduleProposal()
	require.NoError(t, p.ValidateBasic())
}
