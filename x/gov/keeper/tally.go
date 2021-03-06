package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

// DefaultContractAddr is the wasm contract address generated by code ID 1 and instance ID 1.
//
// In other words, the first ever contract to be deployed on this chain will necessarily have this address.
var DefaultContractAddr = wasmkeeper.BuildContractAddress(1, 1)

// Tally iterates over the votes and updates the tally of a proposal based on the voting power of the voters
//
// NOTE: here the voting power of a user is defined as: amount of MARS tokens staked + amount locked in vesting
func (k Keeper) Tally(ctx sdk.Context, proposal govtypes.Proposal) (passes bool, burnDeposits bool, tallyResults govtypes.TallyResult) {
	results := make(map[govtypes.VoteOption]sdk.Dec)
	results[govtypes.OptionYes] = sdk.ZeroDec()
	results[govtypes.OptionAbstain] = sdk.ZeroDec()
	results[govtypes.OptionNo] = sdk.ZeroDec()
	results[govtypes.OptionNoWithVeto] = sdk.ZeroDec()

	// fetch all currently bounded validators
	currValidators := make(map[string]govtypes.ValidatorGovInfo)
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		currValidators[validator.GetOperator().String()] = govtypes.NewValidatorGovInfo(
			validator.GetOperator(),
			validator.GetBondedTokens(),
			validator.GetDelegatorShares(),
			sdk.ZeroDec(),
			govtypes.WeightedVoteOptions{},
		)

		return false
	})

	// fetch all tokens locked in the vesting contract
	//
	// NOTE: for now we simply use the default contract address. later it may be a better idea to use
	// a configurable parameter
	tokensInVesting, totalTokensInVesting := MustGetTokensInVesting(ctx, k.wasmKeeper, DefaultContractAddr)

	// total amount of tokens bonded with validators
	//
	// TODO: does this only include validators in the active set, or also inactive ones?
	totalTokensBonded := k.stakingKeeper.TotalBondedTokens(ctx)

	// total amount of tokens that are eligible to vote in this poll; used to determine quorum
	totalTokens := totalTokensBonded.Add(totalTokensInVesting).ToDec()

	// total amount of tokens that have voted in this poll; used to determine whether the poll reaches
	// quorum and the pass threshold
	totalTokensVoted := sdk.ZeroDec()

	// iterate through votes
	k.IterateVotes(ctx, proposal.ProposalId, func(vote govtypes.Vote) bool {
		voterAddr := sdk.MustAccAddressFromBech32(vote.Voter)

		// if validator, record its vote options in the map
		valAddrStr := sdk.ValAddress(voterAddr).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Options
			currValidators[valAddrStr] = val
		}

		votingPower := sdk.ZeroDec()

		// iterate over all delegations from voter, deduct from any delegated-to validators
		//
		// there is no need to handle the special case that validator address equal to voter address,
		// because voter's voting power will tally again even if there will deduct voter's voting power
		// from validator
		k.stakingKeeper.IterateDelegations(ctx, voterAddr, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
			valAddrStr := delegation.GetValidatorAddr().String()

			if val, ok := currValidators[valAddrStr]; ok {
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				currValidators[valAddrStr] = val

				votingPower = votingPower.Add(delegation.GetShares().MulInt(val.BondedTokens).Quo(val.DelegatorShares))
			}

			return false
		})

		// if the voter has tokens locked in vesting contract, add that to the voting power
		if votingPowerInVesting, ok := tokensInVesting[vote.Voter]; ok {
			votingPower = votingPower.Add(votingPowerInVesting.ToDec())
		}

		incrementTallyResult(votingPower, vote.Options, results, &totalTokensVoted)
		k.deleteVote(ctx, vote.ProposalId, voterAddr)

		return false
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if len(val.Vote) == 0 {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

		incrementTallyResult(votingPower, val.Vote, results, &totalTokensVoted)
	}

	tallyParams := k.GetTallyParams(ctx)
	tallyResults = govtypes.NewTallyResultFromMap(results)

	// if there is no staked coins, the proposal fails
	if k.stakingKeeper.TotalBondedTokens(ctx).IsZero() {
		return false, false, tallyResults
	}

	// if there is not enough quorum of votes, the proposal fails, and deposit burned
	//
	// NOTE: should the deposit really be burned here?
	if totalTokensVoted.Quo(totalTokens).LT(tallyParams.Quorum) {
		return false, true, tallyResults
	}

	// if everyone abstains, proposal fails
	if totalTokensVoted.Sub(results[govtypes.OptionAbstain]).IsZero() {
		return false, false, tallyResults
	}

	// if more than 1/3 of voters veto, proposal fails, and deposit burned
	//
	// NOTE: here 1/3 is defined as 1/3 *of all votes*, including abstaining votes. could it make more
	// sense to instead define it as 1/3 *of all non-abstaining votes*?
	if results[govtypes.OptionNoWithVeto].Quo(totalTokensVoted).GT(tallyParams.VetoThreshold) {
		return false, true, tallyResults
	}

	// if no less than 1/2 of non-abstaining voters vote No, proposal fails
	if results[govtypes.OptionNo].Quo(totalTokensVoted.Sub(results[govtypes.OptionAbstain])).GTE(tallyParams.Threshold) {
		return false, false, tallyResults
	}

	// otherwise, meaning more than 1/2 of non-abstaining voters vote Yes, proposal passes
	return true, false, tallyResults
}

func incrementTallyResult(votingPower sdk.Dec, options []govtypes.WeightedVoteOption, results map[govtypes.VoteOption]sdk.Dec, totalTokensVoted *sdk.Dec) {
	for _, option := range options {
		subPower := votingPower.Mul(option.Weight)
		results[option.Option] = results[option.Option].Add(subPower)
	}

	*totalTokensVoted = totalTokensVoted.Add(votingPower)
}
