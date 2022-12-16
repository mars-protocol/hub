package incentives

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	marsutils "github.com/mars-protocol/hub/utils"

	"github.com/mars-protocol/hub/x/incentives/keeper"
	"github.com/mars-protocol/hub/x/incentives/types"
)

// NewMsgHandler creates a new handler for messages
func NewMsgHandler() sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized incentives message type: %T", msg)
	}
}

// NewProposalHandler creates a new handler for governance proposals
func NewProposalHandler(k keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.CreateIncentivesScheduleProposal:
			return handleCreateIncentivesScheduleProposal(ctx, k, c)
		case *types.TerminateIncentivesSchedulesProposal:
			return handleTerminateIncentivesSchedulesProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized incentives proposal content type: %T", c)
		}
	}
}

func handleCreateIncentivesScheduleProposal(ctx sdk.Context, k keeper.Keeper, p *types.CreateIncentivesScheduleProposal) error {
	schedule, err := k.CreateSchedule(ctx, p.StartTime, p.EndTime, p.Amount)
	if err != nil {
		return nil
	}

	logger := k.Logger(ctx)
	logger.Info("created a new incentives schedule", "id", schedule.Id, "amount", schedule.TotalAmount)

	return nil
}

func handleTerminateIncentivesSchedulesProposal(ctx sdk.Context, k keeper.Keeper, p *types.TerminateIncentivesSchedulesProposal) error {
	amount, err := k.TerminateSchedules(ctx, p.Ids)
	if err != nil {
		return nil
	}

	logger := k.Logger(ctx)
	logger.Info("terminated incentives schedules", "ids", marsutils.UintArrayToString(p.Ids, ","), "amount", amount.String())

	return nil
}
