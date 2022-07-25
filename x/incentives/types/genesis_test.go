package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/incentives/types"
)

func getMockGenesisState() types.GenesisState {
	return types.GenesisState{
		NextScheduleId: 4,
		Schedules: []types.Schedule{
			{
				Id:             2,
				StartTime:      time.Unix(10000, 0),
				EndTime:        time.Unix(20000, 0),
				TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(10000))),
				ReleasedAmount: sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(7500))),
			},
			{

				Id:             3,
				StartTime:      time.Unix(15000, 0),
				EndTime:        time.Unix(25000, 0),
				TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(20000)), sdk.NewCoin("uastro", sdk.NewInt(30000))),
				ReleasedAmount: sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(5000)), sdk.NewCoin("uastro", sdk.NewInt(7500))),
			},
		},
	}
}

func TestNextIdTooSmall(t *testing.T) {
	gs := getMockGenesisState()
	gs.Schedules[1].Id = 5

	require.EqualError(t, gs.Validate(), "incentives schedule id 5 is not smaller than next schedule id 4")
}

func TestDuplicateId(t *testing.T) {
	gs := getMockGenesisState()
	gs.Schedules[1].Id = 2

	require.EqualError(t, gs.Validate(), "incentives schedule has duplicate id 2")
}

func TestEndTimeEarlierThanStart(t *testing.T) {
	gs := getMockGenesisState()
	gs.Schedules[1].EndTime = time.Unix(15000, 0)

	require.EqualError(t, gs.Validate(), "incentives schedule 3 end time is not after start time")
}

func TestZeroTotalAmount(t *testing.T) {
	gs := getMockGenesisState()
	gs.Schedules[1].TotalAmount = sdk.NewCoins()
	gs.Schedules[1].ReleasedAmount = sdk.NewCoins()

	require.EqualError(t, gs.Validate(), "incentives schedule 3 has zero total amount")
}

func TestReleasedAmountGreaterThanTotal(t *testing.T) {
	gs := getMockGenesisState()
	gs.Schedules[1].ReleasedAmount = sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(69420)), sdk.NewCoin("uastro", sdk.NewInt(7500)))

	require.EqualError(t, gs.Validate(), "incentives schedule 3 total amount is not all greater or equal than released amount")
}

func TestValidGenesis(t *testing.T) {
	gs := getMockGenesisState()

	require.NoError(t, gs.Validate())
}
