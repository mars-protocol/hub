package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/mars-protocol/hub/v2/x/gov/keeper"
	"github.com/mars-protocol/hub/v2/x/gov/types"
)

// verify that the test is properly setup
func TestTallyProperSetup(t *testing.T) {
	votingPowers := []VotingPower{
		{Staked: 30_000_000, Vesting: 20_000_000},
		{Staked: 49_000_000, Vesting: 0},
	}

	ctx, app, proposal, valoper, voters := setupTest(t, votingPowers)

	// total staked token amount should be correct
	// 30 from voter[0] + 49 from voter[1] + 1 from valoper
	require.Equal(t, sdk.NewInt(80_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// validator should have been registered
	val, found := app.StakingKeeper.GetValidator(ctx, sdk.ValAddress(valoper))
	require.True(t, found)
	require.Equal(t, sdk.NewInt(80_000_000), val.Tokens)

	// staked token amount for each voter should be correct
	for idx, votingPower := range votingPowers {
		delegation, found := app.StakingKeeper.GetDelegation(ctx, voters[idx], sdk.ValAddress(valoper))
		require.True(t, found)

		staked := delegation.Shares.MulInt(val.Tokens).Quo(val.DelegatorShares).TruncateInt()
		require.Equal(t, sdk.NewInt(votingPower.Staked), staked)
	}

	// vesting token amount for each voter should be correct
	for idx, votingPower := range votingPowers {
		var votingPowerResponse types.VotingPowerResponse

		req, err := json.Marshal(types.QueryMsg{
			VotingPower: &types.VotingPowerQuery{User: voters[idx].String()},
		})
		require.NoError(t, err)

		res, err := app.WasmKeeper.QuerySmart(ctx, keeper.DefaultContractAddr, req)
		require.NoError(t, err)

		err = json.Unmarshal(res, &votingPowerResponse)
		require.NoError(t, err)

		require.Equal(t, sdk.NewInt(votingPower.Vesting), sdk.Int(votingPowerResponse.VotingPower))
	}

	// the proposal should have been created
	_, found = app.GovKeeper.GetProposal(ctx, proposal.Id)
	require.True(t, found)
}

// voters[0] votes with a small voting power; voters[1] with a large voting power does not vote
func TestTallyNoQuorum(t *testing.T) {
	ctx, app, proposal, _, voters := setupTest(t, []VotingPower{
		{Staked: 1_000_000, Vesting: 1_000_000},
		{Staked: 100_000_000, Vesting: 100_000_000},
	})

	// voters[0] votes yes
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, voters[0], govv1.NewNonSplitVoteOption(govv1.OptionYes), ""))

	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
	require.False(t, burnDeposits) // different from native sdk, we don't burn deposit here
	require.Equal(
		t,
		govv1.NewTallyResult(sdk.NewInt(2_000_000), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()),
		tallyResults,
	)
}

// voter[0] has 49 + 50 = 99 voting power, votes abstain
// valoper also votes abstain
// such that all eligible voters vote abstain
func TestTallyOnlyAbstain(t *testing.T) {
	ctx, app, proposal, valoper, voters := setupTest(t, []VotingPower{
		{Staked: 49_000_000, Vesting: 50_000_000},
	})

	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, voters[0], govv1.NewNonSplitVoteOption(govv1.OptionAbstain), ""))
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, valoper, govv1.NewNonSplitVoteOption(govv1.OptionAbstain), ""))

	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
	require.False(t, burnDeposits)
	require.Equal(
		t,
		govv1.NewTallyResult(sdk.ZeroInt(), sdk.NewInt(100_000_000), sdk.ZeroInt(), sdk.ZeroInt()),
		tallyResults,
	)
}

// voter[0] votes veto with 34 voting power
// voter[1] and valoper abstain with their 66 power
// final result: 66 abstain, 34 veto
//
// NOTE: the 1/3 veto threshold refers to 1/3 of *all votes*, including
// abstaining votes
func TestTallyVeto(t *testing.T) {
	ctx, app, proposal, valoper, voters := setupTest(t, []VotingPower{
		{Staked: 0, Vesting: 34_000_000},
		{Staked: 49_000_000, Vesting: 16_000_000},
	})

	// validator abstains
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, valoper, govv1.NewNonSplitVoteOption(govv1.OptionAbstain), ""))

	// voters[0] votes veto
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, voters[0], govv1.NewNonSplitVoteOption(govv1.OptionNoWithVeto), ""))

	// voter[1] abstains
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, voters[1], govv1.NewNonSplitVoteOption(govv1.OptionAbstain), ""))

	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
	require.True(t, burnDeposits)
	require.Equal(
		t,
		govv1.NewTallyResult(sdk.ZeroInt(), sdk.NewInt(66_000_000), sdk.ZeroInt(), sdk.NewInt(34_000_000)),
		tallyResults,
	)
}

// valoper votes no with 1 power
// voter[0] votes no with 50 power
// voter[1] votes yes with 49 power
func TestTallyNo(t *testing.T) {
	ctx, app, proposal, valoper, voters := setupTest(t, []VotingPower{
		{Staked: 25_000_000, Vesting: 25_000_000},
		{Staked: 0, Vesting: 49_000_000},
	})

	// valoper votes no
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, valoper, govv1.NewNonSplitVoteOption(govv1.OptionNo), ""))

	// voters[0] votes no
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, voters[0], govv1.NewNonSplitVoteOption(govv1.OptionNo), ""))

	// voters[1] votes yes
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, voters[1], govv1.NewNonSplitVoteOption(govv1.OptionYes), ""))

	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
	require.False(t, burnDeposits)
	require.Equal(
		t,
		govv1.NewTallyResult(sdk.NewInt(49_000_000), sdk.ZeroInt(), sdk.NewInt(51_000_000), sdk.ZeroInt()),
		tallyResults,
	)
}

// valoper votes yes with 1 power
// voter[0] votes yes with 50 power
// voter[1] votes no with 49 power
func TestTallyYes(t *testing.T) {
	ctx, app, proposal, valoper, voters := setupTest(t, []VotingPower{
		{Staked: 25_000_000, Vesting: 25_000_000},
		{Staked: 0, Vesting: 49_000_000},
	})

	// valoper votes yes
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, valoper, govv1.NewNonSplitVoteOption(govv1.OptionYes), ""))

	// voters[0] votes yes
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, voters[0], govv1.NewNonSplitVoteOption(govv1.OptionYes), ""))

	// voters[1] votes no
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, voters[1], govv1.NewNonSplitVoteOption(govv1.OptionNo), ""))

	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	require.True(t, passes)
	require.False(t, burnDeposits)
	require.Equal(
		t,
		govv1.NewTallyResult(sdk.NewInt(51_000_000), sdk.ZeroInt(), sdk.NewInt(49_000_000), sdk.ZeroInt()),
		tallyResults,
	)
}

// validator has 49 voting power, who votes yes
// voter has 51 total voting power, voting no
// the final result should be 51 no vs 49 yes, proposal fails
func TestTallyValidatorVoteOverride(t *testing.T) {
	ctx, app, proposal, valoper, voters := setupTest(t, []VotingPower{
		{Staked: 30_000_000, Vesting: 21_000_000},
		{Staked: 48_000_000, Vesting: 0},
	})

	// validator votes yes
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, valoper, govv1.NewNonSplitVoteOption(govv1.OptionYes), ""))

	// NOTE: we now delete the votes after tallying, so in order for the 2nd part of this test to work,
	// we have to use a cached context for the 1st part
	cacheCtx, _ := ctx.CacheContext()

	// if voters[0] does not override validator's vote, proposal should pass with 79 yes vs 21 not-voting
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(cacheCtx, proposal)
	require.True(t, passes)
	require.False(t, burnDeposits)
	require.Equal(
		t,
		govv1.NewTallyResult(sdk.NewInt(79_000_000), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()),
		tallyResults,
	)

	// if voters[0] does override validator's vote, proposal should fail with 49 yes vs 51 no
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, voters[0], govv1.NewNonSplitVoteOption(govv1.OptionNo), ""))

	passes, burnDeposits, tallyResults = app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
	require.False(t, burnDeposits)
	require.Equal(
		t,
		govv1.NewTallyResult(sdk.NewInt(49_000_000), sdk.ZeroInt(), sdk.NewInt(51_000_000), sdk.ZeroInt()),
		tallyResults,
	)
}

func TestDeleteVoteAfterTally(t *testing.T) {
	ctx, app, proposal, _, voters := setupTest(t, []VotingPower{{Staked: 1, Vesting: 0}})

	voter := voters[0]

	// the user votes
	app.GovKeeper.SetVote(ctx, govv1.NewVote(proposal.Id, voter, govv1.NewNonSplitVoteOption(govv1.OptionYes), ""))

	// the vote should have been registered
	votes := app.GovKeeper.GetVotes(ctx, proposal.Id)
	require.Equal(t, 1, len(votes))

	_, _, _ = app.GovKeeper.Tally(ctx, proposal)

	// the vote should have been deleted
	votes = app.GovKeeper.GetVotes(ctx, proposal.Id)
	require.Empty(t, votes)
}
