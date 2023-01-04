package keeper_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"

	"github.com/mars-protocol/hub/x/gov/keeper"
	"github.com/mars-protocol/hub/x/gov/testdata"
	"github.com/mars-protocol/hub/x/gov/types"
)

var mockSchedule = &types.Schedule{
	StartTime: 10000,
	Cliff:     0,
	Duration:  1,
}

// VotingPower defines the composition of a mock account's voting power
type VotingPower struct {
	Staked  int64
	Vesting int64
}

func setupTest(t *testing.T, votingPowers []VotingPower) (ctx sdk.Context, app *marsapp.MarsApp, proposal govv1.Proposal, valoper sdk.AccAddress, voters []sdk.AccAddress) {
	accts := marsapptesting.MakeRandomAccounts(len(votingPowers) + 2)
	deployer := accts[0]
	valoper = accts[1]
	voters = accts[2:]

	// calculate the sum of tokens staked and in vesting
	totalStaked := sdk.ZeroInt()
	totalVesting := sdk.ZeroInt()
	for _, votingPower := range votingPowers {
		totalStaked = totalStaked.Add(sdk.NewInt(votingPower.Staked))
		totalVesting = totalVesting.Add(sdk.NewInt(votingPower.Vesting))
	}

	// set mars token balance for deployer and voters
	balances := []banktypes.Balance{{
		Address: deployer.String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(marsapp.BondDenom, totalVesting)),
	}}
	for idx, votingPower := range votingPowers {
		balances = append(balances, banktypes.Balance{
			Address: voters[idx].String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(marsapp.BondDenom, sdk.NewInt(votingPower.Staked))),
		})
	}

	app = marsapptesting.MakeMockApp(accts, balances, []sdk.AccAddress{valoper}, sdk.NewCoins())
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	// register voter accounts at the auth module
	for _, voter := range voters {
		app.AccountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(voter))
	}

	// voters make delegations
	for idx, votingPower := range votingPowers {
		if votingPower.Staked > 0 {
			val, found := app.StakingKeeper.GetValidator(ctx, sdk.ValAddress(valoper))
			require.True(t, found)

			_, err := app.StakingKeeper.Delegate(
				ctx,
				voters[idx],
				sdk.NewInt(votingPower.Staked),
				stakingtypes.Unbonded,
				val,
				true, // true means it's a delegation, not a redelegation
			)
			require.NoError(t, err)
		}
	}

	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(app.WasmKeeper)

	// store vesting contract code
	codeID, _, err := contractKeeper.Create(ctx, deployer, testdata.VestingWasm, nil)
	require.NoError(t, err)

	// instantiate vesting contract
	instantiateMsg, err := json.Marshal(&types.InstantiateMsg{
		Owner:          deployer.String(),
		UnlockSchedule: mockSchedule,
	})
	require.NoError(t, err)

	contractAddr, _, err := contractKeeper.Instantiate(
		ctx,
		codeID,
		deployer,
		nil,
		instantiateMsg,
		"mars/vesting",
		sdk.NewCoins(),
	)
	require.NoError(t, err)

	// create a vesting positions for voters
	for idx, votingPower := range votingPowers {
		if votingPower.Vesting > 0 {
			executeMsg, err := json.Marshal(&types.ExecuteMsg{
				CreatePosition: &types.CreatePosition{
					User:         voters[idx].String(),
					VestSchedule: mockSchedule,
				},
			})
			require.NoError(t, err)

			_, err = contractKeeper.Execute(
				ctx,
				contractAddr,
				deployer,
				executeMsg,
				sdk.NewCoins(sdk.NewCoin(marsapp.BondDenom, sdk.NewInt(votingPower.Vesting))),
			)
			require.NoError(t, err)
		}
	}

	// create a governance proposal
	//
	// typically it requires a minimum deposit to make the proposal enter voting
	// period, but here we forcibly set the status as StatusVotingPeriod.
	//
	// typically we require the proposal's metadata to conform to a schema, but
	// it's not necessary here as we're not creating the proposal through the
	// msgServer.
	proposal, err = govv1.NewProposal([]sdk.Msg{}, 1, "", time.Now(), time.Now())
	proposal.Status = govv1.StatusVotingPeriod
	require.NoError(t, err)

	app.GovKeeper.SetProposal(ctx, proposal)

	return ctx, app, proposal, valoper, voters
}

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
