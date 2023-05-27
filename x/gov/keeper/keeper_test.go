// There is no test in this file.
// We have the setup script used by other tests in this package here.

package keeper_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	marsapp "github.com/mars-protocol/hub/v2/app"
	marsapptesting "github.com/mars-protocol/hub/v2/app/testing"

	"github.com/mars-protocol/hub/v2/x/gov/testdata"
	"github.com/mars-protocol/hub/v2/x/gov/types"
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
