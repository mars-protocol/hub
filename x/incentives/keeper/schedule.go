package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/mars-protocol/hub/v2/x/incentives/types"
)

// CreateSchedule upon a successful CreateIncentivesScheduleProposal, withdraws
// appropriate amount of funds from the community pool, and initializes a new
// schedule in module store. Returns the new schedule that was created.
func (k Keeper) CreateSchedule(ctx sdk.Context, startTime, endTime time.Time, amount sdk.Coins) (schedule types.Schedule, err error) {
	id := k.IncrementNextScheduleID(ctx)

	schedule = types.Schedule{
		Id:             id,
		StartTime:      startTime,
		EndTime:        endTime,
		TotalAmount:    amount,
		ReleasedAmount: sdk.NewCoins(),
	}

	k.SetSchedule(ctx, schedule)

	maccAddr := k.GetModuleAddress()
	if err := k.distrKeeper.DistributeFromFeePool(ctx, amount, maccAddr); err != nil {
		return types.Schedule{}, types.ErrFailedWithdrawFromCommunityPool.Wrap(err.Error())
	}

	return schedule, nil
}

// TerminateSchedules upon a successful TerminateIncentivesScheduleProposal,
// deletes the schedules specified by the proposal from module store, and
// returns the unreleased funds to the community pool. Returns the funds that
// ware returned.
func (k Keeper) TerminateSchedules(ctx sdk.Context, ids []uint64) (amount sdk.Coins, err error) {
	amount = sdk.NewCoins()

	for _, id := range ids {
		schedule, found := k.GetSchedule(ctx, id)
		if !found {
			return sdk.NewCoins(), sdkerrors.ErrKeyNotFound.Wrapf("incentives schedule with id %d does not exist", id)
		}

		amount = amount.Add(schedule.TotalAmount.Sub(schedule.ReleasedAmount...)...)

		k.DeleteSchedule(ctx, id)
	}

	maccAddr := k.GetModuleAddress()
	if err := k.distrKeeper.FundCommunityPool(ctx, amount, maccAddr); err != nil {
		return sdk.NewCoins(), types.ErrFailedRefundToCommunityPool.Wrap(err.Error())
	}

	return amount, nil
}
