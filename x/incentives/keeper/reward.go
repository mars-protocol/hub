package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/mars-protocol/hub/x/incentives/types"
)

func newDecFromInt64(i int64) sdk.Dec {
	return sdk.NewDecFromInt(sdk.NewInt(i))
}

// ReleaseBlockReward handles the release of incentives. Returns the total amount of block reward released
// and the list of relevant schedule ids.
//
// `bondedBotes` is a list of {validator address, validator voted on last block flag} for all validators
// in the bonded set.
func (k Keeper) ReleaseBlockReward(ctx sdk.Context, bondedVotes []abci.VoteInfo) (ids []uint64, totalBlockReward sdk.Coins) {
	currentTime := ctx.BlockTime()

	// iterate through all active schedules, sum up all rewards to be released in this block.
	//
	// If an incentives schedule has been fully released, delete it from the store; otherwise, update
	// the released amount and save
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

	// transfer the coins to distribution module account so that they can be distributed
	k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distrtypes.ModuleName, totalBlockReward)

	// sum up the total voting power voted in the last block
	//
	// TODO: we don't need to check `SignedLastBlock`? the allocate function in distr module doesn't
	totalPower := sdk.ZeroDec()
	for _, vote := range bondedVotes {
		totalPower = totalPower.Add(newDecFromInt64(vote.Validator.Power))
	}

	// allocate reward to validator who have signed the previous block, pro-rate to their voting power
	//
	// NOTE: AllocateTokensToValidator emits the `reward` event, so we don't need to emit separate events
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
