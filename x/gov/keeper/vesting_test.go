package keeper_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"

	"github.com/mars-protocol/hub/x/gov/keeper"
	"github.com/mars-protocol/hub/x/gov/testdata"
	"github.com/mars-protocol/hub/x/gov/types"
)

func TestQueryVotingPowers(t *testing.T) {
	// generate random addresses
	accts := marsapptesting.MakeRandomAccounts(4)
	validator := accts[0]
	deployer := accts[1]
	voters := accts[2:]

	// create mock app and context
	app := marsapptesting.MakeMockApp(
		accts,
		[]banktypes.Balance{{
			Address: deployer.String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(marsapp.BondDenom, sdk.NewInt(50000000))),
		}},
		[]sdk.AccAddress{validator},
		sdk.NewCoins(),
	)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Unix(10000, 0)}) // block time is required for testing

	// take the wasm keeper from the app
	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(app.WasmKeeper)

	// store vesting contract code
	codeID, _, err := contractKeeper.Create(ctx, deployer, testdata.VestingWasm, nil)
	require.NoError(t, err)

	// instantiate vesting contract
	instantiateMsg, err := json.Marshal(&types.InstantiateMsg{
		Owner: deployer.String(),
		UnlockSchedule: &types.Schedule{
			StartTime: 10000,
			Cliff:     0,
			Duration:  1,
		},
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

	// create vesting position for voters[0] with 30_000_000 umars
	executeMsg, err := json.Marshal(&types.ExecuteMsg{
		CreatePosition: &types.CreatePosition{
			User: voters[0].String(),
			VestSchedule: &types.Schedule{
				StartTime: 10000,
				Cliff:     0,
				Duration:  20000,
			},
		},
	})
	require.NoError(t, err)

	_, err = contractKeeper.Execute(
		ctx,
		contractAddr,
		deployer,
		executeMsg,
		sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(30000000))),
	)
	require.NoError(t, err)

	// create vesting position for voters[1] with 20_000_000 umars
	executeMsg, err = json.Marshal(&types.ExecuteMsg{
		CreatePosition: &types.CreatePosition{
			User: voters[1].String(),
			VestSchedule: &types.Schedule{
				StartTime: 0,
				Cliff:     0,
				Duration:  20000,
			},
		},
	})
	require.NoError(t, err)

	_, err = contractKeeper.Execute(
		ctx,
		contractAddr,
		deployer,
		executeMsg,
		sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(20000000))),
	)
	require.NoError(t, err)

	// voters should have 50_000_000 umars locked in vesting combined
	tokensInVesting, totalTokensInVesting, err := keeper.GetTokensInVesting(ctx, app.WasmKeeper, contractAddr)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(50000000), totalTokensInVesting)
	require.Equal(t, sdk.NewInt(30000000), tokensInVesting[voters[0].String()])
	require.Equal(t, sdk.NewInt(20000000), tokensInVesting[voters[1].String()])

	// set time to 20000
	ctx = ctx.WithBlockTime(time.Unix(20000, 0))

	// voters[0] is able to withdraw half of their vested tokens, i.e. 15_000_000 umars
	executeMsg, err = json.Marshal(&types.ExecuteMsg{
		Withdraw: &types.Withdraw{},
	})
	require.NoError(t, err)

	_, err = contractKeeper.Execute(
		ctx,
		contractAddr,
		voters[0],
		executeMsg,
		sdk.NewCoins(),
	)
	require.NoError(t, err)

	tokensInVesting, totalTokensInVesting, err = keeper.GetTokensInVesting(ctx, app.WasmKeeper, contractAddr)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(35000000), totalTokensInVesting)
	require.Equal(t, sdk.NewInt(15000000), tokensInVesting[voters[0].String()])
	require.Equal(t, sdk.NewInt(20000000), tokensInVesting[voters[1].String()])
}
