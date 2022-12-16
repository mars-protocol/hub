package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	marsapptesting "github.com/mars-protocol/hub/app/testing"
)

func TestCreateSchedule(t *testing.T) {
	app := marsapptesting.MakeMockApp()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	accts := marsapptesting.MakeRandomAccounts(1)
	deployer := accts[0]

	// assume we already have mockSchedule[0] active; a successful gov proposal about to add mockSchedule[1]
	app.IncentivesKeeper.SetNextScheduleID(ctx, 2)
	app.IncentivesKeeper.SetSchedule(ctx, mockSchedules[0])

	// assign appropriate amount of coins to the deployer and incentives module account
	maccAddr := app.IncentivesKeeper.GetModuleAddress()
	app.BankKeeper.InitGenesis(
		ctx,
		&banktypes.GenesisState{
			Balances: []banktypes.Balance{{
				Address: deployer.String(),
				Coins:   mockSchedules[1].TotalAmount,
			}, {
				Address: maccAddr.String(),
				Coins:   mockSchedules[0].TotalAmount,
			}},
		},
	)

	// the deployer funds the community pool
	app.DistrKeeper.FundCommunityPool(ctx, mockSchedules[1].TotalAmount, deployer)

	// create the schedule upon a successful governance proposal
	_, err := app.IncentivesKeeper.CreateSchedule(
		ctx,
		mockSchedules[1].StartTime,
		mockSchedules[1].EndTime,
		mockSchedules[1].TotalAmount,
	)
	require.NoError(t, err)

	// next schedule id should have been updated
	nextScheduleId := app.IncentivesKeeper.GetNextScheduleID(ctx)
	require.Equal(t, uint64(3), nextScheduleId)

	// the new schedule should have been saved
	schedule, found := app.IncentivesKeeper.GetSchedule(ctx, 2)
	require.True(t, found)
	require.Equal(t, mockSchedules[1].Id, schedule.Id)
	require.Equal(t, mockSchedules[1].TotalAmount, schedule.TotalAmount)

	// the incentives module account should have been funded
	balances := app.BankKeeper.GetAllBalances(ctx, maccAddr)
	expectedBalances := mockSchedules[0].TotalAmount.Add(mockSchedules[1].TotalAmount...)
	require.Equal(t, expectedBalances, balances)

	// the distribution module account should have been deducted balances
	balances = app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName))
	require.Equal(t, sdk.NewCoins(), balances)

	// the fee pool should have been properly updated
	feePool := app.DistrKeeper.GetFeePool(ctx)
	require.Equal(t, sdk.DecCoins(nil), feePool.CommunityPool)
}

func TestTerminateSchedule(t *testing.T) {
	app := marsapptesting.MakeMockApp()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.IncentivesKeeper.SetNextScheduleID(ctx, 3)
	for _, mockSchedule := range mockSchedulesReleased {
		app.IncentivesKeeper.SetSchedule(ctx, mockSchedule)
	}

	maccAddr := app.IncentivesKeeper.GetModuleAddress()
	amount := sdk.NewCoins()
	for _, mockSchedule := range mockSchedulesReleased {
		amount = amount.Add(mockSchedule.TotalAmount...).Sub(mockSchedule.ReleasedAmount...)
	}

	app.BankKeeper.InitGenesis(
		ctx,
		&banktypes.GenesisState{
			Balances: []banktypes.Balance{{
				Address: maccAddr.String(),
				Coins:   amount,
			}},
		},
	)

	// terminate the schedules upon a successful governance proposal
	_, err := app.IncentivesKeeper.TerminateSchedules(ctx, []uint64{1, 2})
	require.NoError(t, err)

	// next schedule id should have not been changed
	nextScheduleId := app.IncentivesKeeper.GetNextScheduleID(ctx)
	require.Equal(t, uint64(3), nextScheduleId)

	// the two schedules should have been deleted
	_, found := app.IncentivesKeeper.GetSchedule(ctx, 1)
	require.False(t, found)
	_, found = app.IncentivesKeeper.GetSchedule(ctx, 2)
	require.False(t, found)

	// the incentives module account should have been deducted balance
	balances := app.BankKeeper.GetAllBalances(ctx, maccAddr)
	require.Equal(t, sdk.NewCoins(), balances)

	// the distribution module account should have been funded
	balances = app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName))
	require.Equal(t, amount, balances)

	// the fee pool should have been properly updated
	feePool := app.DistrKeeper.GetFeePool(ctx)
	require.Equal(t, sdk.NewDecCoinsFromCoins(amount...), feePool.CommunityPool)
}
