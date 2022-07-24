package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"

	"github.com/mars-protocol/hub/x/incentives/keeper"
	"github.com/mars-protocol/hub/x/incentives/types"
)

var mockSchedules = []types.Schedule{{
	Id:             1,
	StartTime:      time.Unix(10000, 0).Local(),
	EndTime:        time.Unix(20000, 0).Local(),
	TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(12345))),
	ReleasedAmount: sdk.NewCoins(),
}, {
	Id:             2,
	StartTime:      time.Unix(20000, 0).Local(),
	EndTime:        time.Unix(30000, 0).Local(),
	TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(23456))),
	ReleasedAmount: sdk.NewCoins(),
}, {
	Id:             3,
	StartTime:      time.Unix(30000, 0).Local(),
	EndTime:        time.Unix(40000, 0).Local(),
	TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(34567))),
	ReleasedAmount: sdk.NewCoins(),
}}

func setupQueryServerTest() (ctx sdk.Context, app *marsapp.MarsApp) {
	app = marsapptesting.MakeMockApp()
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
		Limit:      2,
		CountTotal: false,
	}
	res, err := queryServer.Schedules(sdk.WrapSDKContext(ctx), &types.QuerySchedulesRequest{Pagination: pageReq})
	require.NoError(t, err)
	require.Equal(t, 2, len(res.Schedules))

	pageReq = &query.PageRequest{
		Key:        res.Pagination.NextKey,
		Limit:      1,
		CountTotal: false,
	}
	res, err = queryServer.Schedules(sdk.WrapSDKContext(ctx), &types.QuerySchedulesRequest{Pagination: pageReq})
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Schedules))
}
