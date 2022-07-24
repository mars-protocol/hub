package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	marsapptesting "github.com/mars-protocol/hub/app/testing"

	"github.com/mars-protocol/hub/x/incentives/keeper"
	"github.com/mars-protocol/hub/x/incentives/types"
)

func TestUnreleasedIncentivesInvariant(t *testing.T) {
	app := marsapptesting.MakeMockApp()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.IncentivesKeeper.SetSchedule(
		ctx,
		types.Schedule{
			Id:             1,
			StartTime:      time.Unix(10000, 0),
			EndTime:        time.Unix(20000, 0),
			TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(12345)), sdk.NewCoin("uastro", sdk.NewInt(69420))),
			ReleasedAmount: sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(11066)), sdk.NewCoin("uastro", sdk.NewInt(62228))),
		},
	)
	app.IncentivesKeeper.SetSchedule(
		ctx,
		types.Schedule{
			Id:             2,
			StartTime:      time.Unix(15000, 0),
			EndTime:        time.Unix(30000, 0),
			TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(10000))),
			ReleasedAmount: sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(2642))),
		},
	)

	invariant := keeper.TotalUnreleasedIncentives(app.IncentivesKeeper)

	// set incorrect balances for the incentives module account
	maccAddr := app.IncentivesKeeper.GetModuleAddress(ctx)
	app.BankKeeper.InitGenesis(
		ctx,
		&banktypes.GenesisState{
			Balances: []banktypes.Balance{{
				Address: maccAddr.String(),
				Coins:   sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(123)), sdk.NewCoin("uastro", sdk.NewInt(456))),
			}},
		},
	)

	msg, broken := invariant(ctx)
	require.True(t, broken)
	require.Equal(
		t,
		`incentives: total-unreleased-incentives invariant
	sum of unreleased incentives: 7192uastro,8637umars
	module account balances: 456uastro,123umars
`,
		msg,
	)

	// set the correct balances for the incentives module account
	app.BankKeeper.InitGenesis(
		ctx,
		&banktypes.GenesisState{
			Balances: []banktypes.Balance{{
				Address: maccAddr.String(),
				Coins:   sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(8637)), sdk.NewCoin("uastro", sdk.NewInt(7192))),
			}},
		},
	)

	_, broken = invariant(ctx)
	require.False(t, broken)
}
