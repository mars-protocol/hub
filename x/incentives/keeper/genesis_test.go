package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"
)

func setupGenesisTest() (ctx sdk.Context, app *marsapp.MarsApp) {
	app = marsapptesting.MakeSimpleMockApp()
	ctx = app.BaseApp.NewContext(false, tmproto.Header{})

	app.IncentivesKeeper.InitGenesis(ctx, &mockGenesisState)

	return ctx, app
}

func TestInitGenesis(t *testing.T) {
	ctx, app := setupGenesisTest()

	nextScheduleId := app.IncentivesKeeper.GetNextScheduleID(ctx)
	require.Equal(t, uint64(3), nextScheduleId)

	for _, mockSchedule := range mockSchedules {
		schedule, found := app.IncentivesKeeper.GetSchedule(ctx, mockSchedule.Id)
		require.True(t, found)
		require.Equal(t, mockSchedule.Id, schedule.Id)
		require.Equal(t, mockSchedule.TotalAmount, schedule.TotalAmount)
	}
}

func TestExportGenesis(t *testing.T) {
	ctx, app := setupGenesisTest()

	exported := app.IncentivesKeeper.ExportGenesis(ctx)
	require.Equal(t, mockGenesisState.NextScheduleId, exported.NextScheduleId)
	for idx := range mockSchedules {
		require.Equal(t, mockSchedules[idx].Id, exported.Schedules[idx].Id)
		require.Equal(t, mockSchedules[idx].TotalAmount, exported.Schedules[idx].TotalAmount)
	}
}
