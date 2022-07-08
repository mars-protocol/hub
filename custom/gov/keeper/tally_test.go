package keeper_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	simapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"

	"github.com/mars-protocol/hub/custom/gov/keeper"
	"github.com/mars-protocol/hub/custom/gov/testdata"
	"github.com/mars-protocol/hub/custom/gov/types"
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

func setupTest(t *testing.T, votingPowers []VotingPower) (ctx sdk.Context, app *marsapp.MarsApp, proposal govtypes.Proposal, valoper sdk.AccAddress, voters []sdk.AccAddress) {
	app = marsapptesting.MakeMockApp()
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	accts := marsapptesting.MakeRandomAccounts(len(votingPowers) + 2)
	deployer := accts[0]
	valoper = accts[1]
	voters = accts[2:]

	pks := simapp.CreateTestPubKeys(1)
	valPubKey := pks[0]

	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(app.WasmKeeper)

	// calculate the sum of tokens staked and in vesting
	totalStaked := sdk.ZeroInt()
	totalVesting := sdk.ZeroInt()
	for _, votingPower := range votingPowers {
		totalStaked = totalStaked.Add(sdk.NewInt(votingPower.Staked))
		totalVesting = totalVesting.Add(sdk.NewInt(votingPower.Vesting))
	}

	// register voter accounts at the auth module
	for _, voter := range voters {
		app.AccountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(voter))
	}

	// set mars token balance for deployer and voters
	balances := []banktypes.Balance{{
		Address: deployer.String(),
		Coins:   sdk.NewCoins(sdk.NewCoin("umars", totalVesting)),
	}}
	for idx, votingPower := range votingPowers {
		balances = append(balances, banktypes.Balance{
			Address: voters[idx].String(),
			Coins:   sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(votingPower.Staked))),
		})
	}

	app.BankKeeper.InitGenesis(
		ctx,
		&banktypes.GenesisState{
			Params: banktypes.Params{
				DefaultSendEnabled: true, // must set this to true so that tokens can be transferred
			},
			Balances: balances,
		},
	)

	// set bond denom to `umars`
	stakingParams := app.StakingKeeper.GetParams(ctx)
	stakingParams.BondDenom = "umars"
	app.StakingKeeper.SetParams(ctx, stakingParams)

	// create validator
	// NOTE: the validator's status must be set to as bonded
	val, err := stakingtypes.NewValidator(sdk.ValAddress(valoper), valPubKey, stakingtypes.Description{})
	val.Status = stakingtypes.Bonded
	require.NoError(t, err)
	require.True(t, val.IsBonded())

	app.StakingKeeper.SetValidator(ctx, val)
	app.StakingKeeper.SetValidatorByConsAddr(ctx, val)
	app.StakingKeeper.SetValidatorByPowerIndex(ctx, val)
	app.StakingKeeper.AfterValidatorCreated(ctx, val.GetOperator()) // required to initialize distr keeper properly

	// voters make delegations
	for idx, votingPower := range votingPowers {
		if votingPower.Staked > 0 {
			val, found := app.StakingKeeper.GetValidator(ctx, sdk.ValAddress(valoper))
			require.True(t, found)

			_, err = app.StakingKeeper.Delegate(
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

	// store vesting contract code
	codeID, err := contractKeeper.Create(ctx, deployer, testdata.VestingWasm, nil)
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
				sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(votingPower.Vesting))),
			)
			require.NoError(t, err)
		}

	}

	// create a governance proposal
	proposal, err = app.GovKeeper.SubmitProposal(ctx, govtypes.NewTextProposal("mock title", "mock description"))
	require.NoError(t, err)

	return ctx, app, proposal, valoper, voters
}

// verify that the test is properly setup
func TestTallyProperSetup(t *testing.T) {
	votingPowers := []VotingPower{
		{Staked: 30, Vesting: 21},
		{Staked: 49, Vesting: 0},
	}

	ctx, app, proposal, valoper, voters := setupTest(t, votingPowers)

	// total staked token amount should be correct
	require.Equal(t, sdk.NewInt(79), app.StakingKeeper.TotalBondedTokens(ctx))

	// validator should have been registered
	val, found := app.StakingKeeper.GetValidator(ctx, sdk.ValAddress(valoper))
	require.True(t, found)
	require.Equal(t, sdk.NewInt(79), val.Tokens)

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
	_, found = app.GovKeeper.GetProposal(ctx, proposal.ProposalId)
	require.True(t, found)
}

// voters[0] votes with a small voting power; voters[1] with a large voting power does not vote
func TestTallyNoQuorum(t *testing.T) {
	ctx, app, proposal, _, voters := setupTest(t, []VotingPower{
		{Staked: 1, Vesting: 1},
		{Staked: 100, Vesting: 100},
	})

	// voters[0] votes yes
	app.GovKeeper.SetVote(ctx, govtypes.NewVote(proposal.ProposalId, voters[0], govtypes.NewNonSplitVoteOption(govtypes.OptionYes)))

	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
	require.True(t, burnDeposits)
	require.Equal(
		t,
		govtypes.NewTallyResult(sdk.NewInt(2), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()),
		tallyResults,
	)
}

func TestTallyOnlyAbstain(t *testing.T) {
	ctx, app, proposal, _, voters := setupTest(t, []VotingPower{{Staked: 50, Vesting: 50}})

	app.GovKeeper.SetVote(ctx, govtypes.NewVote(proposal.ProposalId, voters[0], govtypes.NewNonSplitVoteOption(govtypes.OptionAbstain)))

	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
	require.False(t, burnDeposits)
	require.Equal(
		t,
		govtypes.NewTallyResult(sdk.ZeroInt(), sdk.NewInt(100), sdk.ZeroInt(), sdk.ZeroInt()),
		tallyResults,
	)
}

func TestTallyVeto(t *testing.T) {
	ctx, app, proposal, valoper, voters := setupTest(t, []VotingPower{
		{Staked: 0, Vesting: 34},
		{Staked: 50, Vesting: 16}, // NOTE: validator only votes with the 50 staked tokens, not the 16 vesting tokens
	})

	// validator abstains
	app.GovKeeper.SetVote(ctx, govtypes.NewVote(proposal.ProposalId, valoper, govtypes.NewNonSplitVoteOption(govtypes.OptionAbstain)))

	// voters[0] votes veto
	app.GovKeeper.SetVote(ctx, govtypes.NewVote(proposal.ProposalId, voters[0], govtypes.NewNonSplitVoteOption(govtypes.OptionNoWithVeto)))

	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
	require.True(t, burnDeposits)
	require.Equal(
		t,
		govtypes.NewTallyResult(sdk.ZeroInt(), sdk.NewInt(50), sdk.ZeroInt(), sdk.NewInt(34)),
		tallyResults,
	)
}

func TestTallyNo(t *testing.T) {
	ctx, app, proposal, _, voters := setupTest(t, []VotingPower{
		{Staked: 25, Vesting: 26},
		{Staked: 0, Vesting: 49},
	})

	// voters[0] votes no
	app.GovKeeper.SetVote(ctx, govtypes.NewVote(proposal.ProposalId, voters[0], govtypes.NewNonSplitVoteOption(govtypes.OptionNo)))

	// voters[1] votes yes
	app.GovKeeper.SetVote(ctx, govtypes.NewVote(proposal.ProposalId, voters[1], govtypes.NewNonSplitVoteOption(govtypes.OptionYes)))

	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
	require.False(t, burnDeposits)
	require.Equal(
		t,
		govtypes.NewTallyResult(sdk.NewInt(49), sdk.ZeroInt(), sdk.NewInt(51), sdk.ZeroInt()),
		tallyResults,
	)
}

func TestTallyYes(t *testing.T) {
	ctx, app, proposal, _, voters := setupTest(t, []VotingPower{
		{Staked: 25, Vesting: 26},
		{Staked: 0, Vesting: 49},
	})

	// voters[0] votes yes
	app.GovKeeper.SetVote(ctx, govtypes.NewVote(proposal.ProposalId, voters[0], govtypes.NewNonSplitVoteOption(govtypes.OptionYes)))

	// voters[1] votes no
	app.GovKeeper.SetVote(ctx, govtypes.NewVote(proposal.ProposalId, voters[1], govtypes.NewNonSplitVoteOption(govtypes.OptionNo)))

	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	require.True(t, passes)
	require.False(t, burnDeposits)
	require.Equal(
		t,
		govtypes.NewTallyResult(sdk.NewInt(51), sdk.ZeroInt(), sdk.NewInt(49), sdk.ZeroInt()),
		tallyResults,
	)
}

// validator has 49 voting power, who votes yes
// voter has 51 total voting power, voting no
// the final result should be 51 no vs 49 yes, proposal fails
func TestTallyValidatorVoteOverride(t *testing.T) {
	ctx, app, proposal, valoper, voters := setupTest(t, []VotingPower{
		{Staked: 30, Vesting: 21},
		{Staked: 49, Vesting: 0},
	})

	// validator votes yes
	app.GovKeeper.SetVote(ctx, govtypes.NewVote(proposal.ProposalId, valoper, govtypes.NewNonSplitVoteOption(govtypes.OptionYes)))

	// if voters[0] does not override validator's vote, proposal should pass with 79 yes vs 21 not-voting
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	require.True(t, passes)
	require.False(t, burnDeposits)
	require.Equal(
		t,
		govtypes.NewTallyResult(sdk.NewInt(79), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()),
		tallyResults,
	)

	// if voters[0] does override validator's vote, proposal should fail with 49 yes vs 51 no
	app.GovKeeper.SetVote(ctx, govtypes.NewVote(proposal.ProposalId, voters[0], govtypes.NewNonSplitVoteOption(govtypes.OptionNo)))

	passes, burnDeposits, tallyResults = app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
	require.False(t, burnDeposits)
	require.Equal(
		t,
		govtypes.NewTallyResult(sdk.NewInt(49), sdk.ZeroInt(), sdk.NewInt(51), sdk.ZeroInt()),
		tallyResults,
	)
}
