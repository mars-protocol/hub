package incentives_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"

	"github.com/mars-protocol/hub/x/incentives"
	"github.com/mars-protocol/hub/x/incentives/types"
)

var mockSchedule = types.Schedule{
	Id:             1,
	StartTime:      time.Unix(10000, 0),
	EndTime:        time.Unix(20000, 0),
	TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(12345)), sdk.NewCoin("uastro", sdk.NewInt(69420))),
	ReleasedAmount: sdk.Coins(nil),
}

func setupHandlerTest() (ctx sdk.Context, app *marsapp.MarsApp) {
	app = marsapptesting.MakeMockApp()
	ctx = app.BaseApp.NewContext(false, tmproto.Header{})

	maccAddr := app.IncentivesKeeper.GetModuleAddress()

	app.BankKeeper.InitGenesis(
		ctx,
		&banktypes.GenesisState{
			Params: banktypes.Params{
				DefaultSendEnabled: true, // must set this to true so that tokens can be transferred
			},
			Balances: []banktypes.Balance{{
				Address: maccAddr.String(),
				Coins:   mockSchedule.TotalAmount,
			}},
		},
	)

	app.IncentivesKeeper.InitGenesis(
		ctx,
		&types.GenesisState{
			NextScheduleId: 1,
			Schedules:      []types.Schedule{},
		},
	)

	return ctx, app
}

func TestCreateScheduleProposalPassed(t *testing.T) {
	ctx, app := setupHandlerTest()

	hdlr := incentives.NewProposalHandler(app.IncentivesKeeper)
	proposal := &types.CreateIncentivesScheduleProposal{
		Title:       "title",
		Description: "description",
		StartTime:   mockSchedule.StartTime,
		EndTime:     mockSchedule.EndTime,
		Amount:      mockSchedule.TotalAmount,
	}
	require.NoError(t, hdlr(ctx, proposal))

	_, found := app.IncentivesKeeper.GetSchedule(ctx, 1)
	require.True(t, found)
}

func TestTerminateSchedulesProposalPassed(t *testing.T) {
	ctx, app := setupHandlerTest()

	app.IncentivesKeeper.SetSchedule(ctx, mockSchedule)

	hdlr := incentives.NewProposalHandler(app.IncentivesKeeper)
	proposal := &types.TerminateIncentivesSchedulesProposal{
		Title:       "title",
		Description: "description",
		Ids:         []uint64{1},
	}
	require.NoError(t, hdlr(ctx, proposal))

	_, found := app.IncentivesKeeper.GetSchedule(ctx, 1)
	require.False(t, found)
}
