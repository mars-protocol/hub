package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/mars-protocol/hub/v2/x/incentives/types"
)

func newDecFromInt64(i int64) sdk.Dec {
	return sdk.NewDecFromInt(sdk.NewInt(i))
}

// ReleaseBlockReward handles the release of incentives. Returns the total
// amount of block reward released and the list of relevant schedule ids.
//
// `bondedVotes` is a list of {validator address, validator voted on last block
// flag} for all validators in the bonded set.
func (k Keeper) ReleaseBlockReward(ctx sdk.Context, bondedVotes []abci.VoteInfo) (ids []uint64, totalBlockReward sdk.Coins) {
	currentTime := ctx.BlockTime()

	// iterate through all active schedules, sum up all rewards to be released
	// in this block.
	//
	// If an incentives schedule has been fully released, delete it from the
	// store; otherwise, update the released amount and save
	ids = []uint64{}
	totalBlockReward = sdk.NewCoins()
	k.IterateSchedules(ctx, func(schedule types.Schedule) bool {
		blockReward := schedule.GetBlockReward(currentTime)

		if !blockReward.Empty() {
			ids = append(ids, schedule.Id)
			totalBlockReward = totalBlockReward.Add(blockReward...)
		}

		if currentTime.After(schedule.EndTime) {
			k.DeleteSchedule(ctx, schedule.Id)
		} else {
			schedule.ReleasedAmount = schedule.ReleasedAmount.Add(blockReward...)
			k.SetSchedule(ctx, schedule)
		}

		return false
	})

	// exit here if there is no coin to be released
	if totalBlockReward.Empty() {
		return
	}

	// transfer the coins to distribution module account so that they can be
	// distributed
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distrtypes.ModuleName, totalBlockReward)
	if err != nil {
		panic(err)
	}

	// sum up the total voting power voted in the last block
	//
	// NOTE: The following code is copied from cosmos-sdk's distribution module
	// without change.
	// Here the distr module adds up voting power of _all_ validators without
	// checking whether the validator has signed the previous block or not.
	// In other words, there is no "micro-slashing" for missing single blocks.
	// We keep this behavior without change.
	// More on this issue: https://twitter.com/larry0x/status/1588189416257880064
	totalPower := sdk.ZeroDec()
	for _, vote := range bondedVotes {
		totalPower = totalPower.Add(newDecFromInt64(vote.Validator.Power))
	}

	// allocate reward to validator who have signed the previous block, pro-rata
	// to their voting power
	//
	// NOTE: AllocateTokensToValidator emits the `reward` event, so we don't
	// need to emit separate events
	totalBlockRewardDec := sdk.NewDecCoinsFromCoins(totalBlockReward...)
	for _, vote := range bondedVotes {
		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)

		power := newDecFromInt64(vote.Validator.Power)
		reward := totalBlockRewardDec.MulDec(power).QuoDec(totalPower)

		totalPower = totalPower.Sub(power)
		totalBlockRewardDec = totalBlockRewardDec.Sub(reward)

		k.distrKeeper.AllocateTokensToValidator(ctx, validator, reward)
	}

	return ids, totalBlockReward
}
