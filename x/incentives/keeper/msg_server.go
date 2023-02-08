package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	marsutils "github.com/mars-protocol/hub/utils"

	"github.com/mars-protocol/hub/x/incentives/types"
)

type msgServer struct{ k Keeper }

// NewMsgServerImpl creates an implementation of the `MsgServer` interface for
// the given keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{k}
}

func (ms msgServer) CreateSchedule(goCtx context.Context, req *types.MsgCreateSchedule) (*types.MsgCreateScheduleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if req.Authority != ms.k.authority {
		return nil, govtypes.ErrInvalidSigner.Wrapf("expected %s got %s", ms.k.authority, req.Authority)
	}

	schedule, err := ms.k.CreateSchedule(ctx, req.StartTime, req.EndTime, req.Amount)
	if err != nil {
		return nil, err
	}

	ms.k.Logger(ctx).Info(
		"incentives schedule created",
		"id", schedule.Id,
		"amount", schedule.TotalAmount.String(),
		"startTime", schedule.StartTime.String(),
		"endTime", schedule.EndTime.String(),
	)

	return &types.MsgCreateScheduleResponse{}, nil
}

func (ms msgServer) TerminateSchedules(goCtx context.Context, req *types.MsgTerminateSchedules) (*types.MsgTerminateSchedulesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if req.Authority != ms.k.authority {
		return nil, govtypes.ErrInvalidSigner.Wrapf("expected %s got %s", ms.k.authority, req.Authority)
	}

	amount, err := ms.k.TerminateSchedules(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	ms.k.Logger(ctx).Info(
		"incentives schedule terminated",
		"ids", marsutils.UintArrayToString(req.Ids, ","),
		"refundedAmount", amount.String(),
	)

	return &types.MsgTerminateSchedulesResponse{RefundedAmount: amount}, nil
}
