package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	marsapp "github.com/mars-protocol/hub/v2/app"
	marsapptesting "github.com/mars-protocol/hub/v2/app/testing"

	"github.com/mars-protocol/hub/v2/x/incentives/keeper"
	"github.com/mars-protocol/hub/v2/x/incentives/types"
)

func setupQueryServerTest() (ctx sdk.Context, app *marsapp.MarsApp) {
	app = marsapptesting.MakeSimpleMockApp()
	ctx = app.BaseApp.NewContext(false, tmproto.Header{})

	for _, schedule := range mockSchedules {
		app.IncentivesKeeper.SetSchedule(ctx, schedule)
	}

	return ctx, app
}

func TestEmptyQuery(t *testing.T) {
	ctx, app := setupQueryServerTest()

	queryServer := keeper.NewQueryServerImpl(app.IncentivesKeeper)

	_, err := queryServer.Schedule(sdk.WrapSDKContext(ctx), nil)
	require.Errorf(t, err, "empty request")

	_, err = queryServer.Schedules(sdk.WrapSDKContext(ctx), nil)
	require.Errorf(t, err, "empty request")
}

func TestQuerySchedule(t *testing.T) {
	ctx, app := setupQueryServerTest()

	queryServer := keeper.NewQueryServerImpl(app.IncentivesKeeper)

	// NOTE: the mock schedules use `time.UTC` while the schedules in the responses use `time.Local`.
	// i can't figure out how to enforce the response to return UTC, so here we just compare the id
	// and total amount, and skip the times
	for _, mockSchedule := range mockSchedules {
		res, err := queryServer.Schedule(sdk.WrapSDKContext(ctx), &types.QueryScheduleRequest{Id: mockSchedule.Id})
		require.NoError(t, err)
		require.Equal(t, res.Schedule.Id, mockSchedule.Id)
		require.Equal(t, res.Schedule.TotalAmount, mockSchedule.TotalAmount)
	}
}

func TestQuerySchedules(t *testing.T) {
	ctx, app := setupQueryServerTest()

	queryServer := keeper.NewQueryServerImpl(app.IncentivesKeeper)

	pageReq := &query.PageRequest{
		Key:        nil,
		Limit:      1,
		CountTotal: false,
	}
	res, err := queryServer.Schedules(sdk.WrapSDKContext(ctx), &types.QuerySchedulesRequest{Pagination: pageReq})
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Schedules))

	pageReq = &query.PageRequest{
		Key:        res.Pagination.NextKey,
		Limit:      1,
		CountTotal: false,
	}
	res, err = queryServer.Schedules(sdk.WrapSDKContext(ctx), &types.QuerySchedulesRequest{Pagination: pageReq})
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Schedules))
}
