package keeper_test

import (
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"

	"github.com/mars-protocol/hub/x/incentives/types"
)

func setupTest(t *testing.T, schedules []types.Schedule) (ctx sdk.Context, app *marsapp.MarsApp) {
	app = marsapptesting.MakeMockApp()
	ctx = app.BaseApp.NewContext(false, tmproto.Header{})

	accts := marsapptesting.MakeRandomAccounts(1)
	staker := accts[0]

	pks := simapp.CreateTestPubKeys(1)
	valPubKey := pks[0]

	// calculate the total mars token amount needed to be given to incentives module account
	totalIncentives := sdk.NewCoins()
	for _, schedule := range schedules {
		totalIncentives = totalIncentives.Add(schedule.TotalAmount...)
	}

	// set mars token balance for the incentives module account
	macc := app.IncentivesKeeper.GetModuleAccount(ctx)
	app.BankKeeper.InitGenesis(
		ctx,
		&banktypes.GenesisState{
			Params: banktypes.Params{
				DefaultSendEnabled: true, // must set this to true so that tokens can be transferred
			},
			Balances: []banktypes.Balance{{
				Address: macc.GetAddress().String(),
				Coins:   totalIncentives,
			}},
		},
	)

	// initialize each validator's accumulated rewards

	return ctx, app
}

func TestReleaseBlockRewardNoActiveSchedule(t *testing.T) {

}
