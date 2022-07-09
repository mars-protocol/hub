package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/custom/distribution/types"
)

// AllocateTokens handles the distribution of fees collected in the previous block
//
// bondedVotes is a list of {validator address, validator voted on last block flag} tuple for all
// validators in the bonded set.
//
// NOTE: this function not super optimized, as `k.authKeeper.GetModuleAccount` and `k.bankKeeper.GetAllBalances`
// are invoked twice both here and in `k.Keeper.AllocateTokens`. however at this time i consider it
// not worth to modify the whole function just to optimize this little bit
func (k Keeper) AllocateTokens(
	ctx sdk.Context, sumPreviousPrecommitPower, totalPreviousPower int64,
	previousProposer sdk.ConsAddress, bondedVotes []abci.VoteInfo,
) {
	// fetch and clear the collected fees for distribution
	//
	// since this is called in BeginBlock, collected fees will be from the previous block, and
	// distributed to the previous proposer
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollected := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	feesCollectedRewardDenom := sdk.NewCoin(types.RewardDenom, feesCollected.AmountOf(types.RewardDenom))
	feesCollectedNonRewardDenom := feesCollected.Sub(sdk.NewCoins(feesCollectedRewardDenom))

	// fees that are NOT in the reward denom go directly to the community pool
	k.FundCommunityPool(ctx, feesCollectedNonRewardDenom, feeCollector.GetAddress())

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSafetyFund,
			sdk.NewAttribute(sdk.AttributeKeyAmount, feesCollectedNonRewardDenom.String()),
		),
	)

	// for fees collected in the reward denom, we simply forward them to the vanilla `AllocateTokens` method
	k.Keeper.AllocateTokens(ctx, sumPreviousPrecommitPower, totalPreviousPower, previousProposer, bondedVotes)
}
