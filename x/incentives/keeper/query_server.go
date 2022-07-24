package keeper

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	query "github.com/cosmos/cosmos-sdk/types/query"

	"github.com/mars-protocol/hub/x/incentives/types"
)

type queryServer struct{ k Keeper }

// NewQueryServerImpl creates an implementation of the `QueryServer` interface for the given keeper
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{k}
}

func (qs queryServer) Schedule(goCtx context.Context, req *types.QueryScheduleRequest) (*types.QueryScheduleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	schedule, found := qs.k.GetSchedule(ctx, req.Id)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "incentives schedule not found for id %d", req.Id)
	}

	return &types.QueryScheduleResponse{Schedule: schedule}, nil
}

func (qs queryServer) Schedules(goCtx context.Context, req *types.QuerySchedulesRequest) (*types.QuerySchedulesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	schedules := []types.Schedule{}

	pageRes, err := query.Paginate(qs.k.GetSchedulePrefixStore(ctx), req.Pagination, func(_, value []byte) error {
		var schedule types.Schedule
		fmt.Println("value:", value)
		err := qs.k.cdc.Unmarshal(value, &schedule)
		if err != nil {
			return err
		}

		schedules = append(schedules, schedule)

		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QuerySchedulesResponse{Schedules: schedules, Pagination: pageRes}, nil
}
