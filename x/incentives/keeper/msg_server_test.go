package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	marsapp "github.com/mars-protocol/hub/v2/app"
	marsapptesting "github.com/mars-protocol/hub/v2/app/testing"

	"github.com/mars-protocol/hub/v2/x/incentives/keeper"
	"github.com/mars-protocol/hub/v2/x/incentives/types"
)

var mockSchedule = types.Schedule{
	Id:             1,
	StartTime:      time.Unix(10000, 0),
	EndTime:        time.Unix(20000, 0),
	TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(12345)), sdk.NewCoin("uastro", sdk.NewInt(69420))),
	ReleasedAmount: sdk.Coins(nil),
}

const (
	govModuleAccount    = "mars10d07y265gmmuvt4z0w9aw880jnsr700j8l2urg"
	notGovModuleAccount = "mars1z926ax906k0ycsuckele6x5hh66e2m4m09whw6"
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount("mars", "marspub")
}

func setupMsgServerTest() (ctx sdk.Context, app *marsapp.MarsApp) {
	accts := marsapptesting.MakeRandomAccounts(1)
	maccAddr := authtypes.NewModuleAddress(types.ModuleName)

	// we give sufficient token amounts to both the community pool and
	// incentives module account, so that we can both create or terminate
	// schedules.
	app = marsapptesting.MakeMockApp(
		accts,
		[]banktypes.Balance{{
			Address: maccAddr.String(),
			Coins:   mockSchedule.TotalAmount,
		}},
		accts,
		mockSchedule.TotalAmount,
	)
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Time: time.Unix(16667, 0)})

	return ctx, app
}

func TestCreateScheduleProposalPassed(t *testing.T) {
	ctx, app := setupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.IncentivesKeeper)
	req := &types.MsgCreateSchedule{
		Authority: govModuleAccount,
		StartTime: mockSchedule.StartTime,
		EndTime:   mockSchedule.EndTime,
		Amount:    mockSchedule.TotalAmount,
	}
	_, err := msgServer.CreateSchedule(ctx, req)
	require.NoError(t, err)

	_, found := app.IncentivesKeeper.GetSchedule(ctx, 1)
	require.True(t, found)
}

func TestTerminateSchedulesProposalPassed(t *testing.T) {
	ctx, app := setupMsgServerTest()

	app.IncentivesKeeper.SetSchedule(ctx, mockSchedule)

	msgServer := keeper.NewMsgServerImpl(app.IncentivesKeeper)
	req := &types.MsgTerminateSchedules{
		Authority: govModuleAccount,
		Ids:       []uint64{1},
	}
	_, err := msgServer.TerminateSchedules(ctx, req)
	require.NoError(t, err)

	_, found := app.IncentivesKeeper.GetSchedule(ctx, 1)
	require.False(t, found)
}

func TestNotAuthority(t *testing.T) {
	ctx, app := setupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.IncentivesKeeper)
	req := &types.MsgCreateSchedule{
		Authority: notGovModuleAccount,
		StartTime: mockSchedule.StartTime,
		EndTime:   mockSchedule.EndTime,
		Amount:    mockSchedule.TotalAmount,
	}
	_, err := msgServer.CreateSchedule(ctx, req)
	require.Error(t, err, govtypes.ErrInvalidSigner)
}
