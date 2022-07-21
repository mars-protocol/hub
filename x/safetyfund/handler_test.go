package safetyfund_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"

	"github.com/mars-protocol/hub/x/safetyfund"
	"github.com/mars-protocol/hub/x/safetyfund/types"
)

var (
	recipientAddr = marsapptesting.MakeRandomAccounts(1)[0]
	amount        = sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(69420)))
	proposal      = types.NewSafetyFundSpendProposal("test", "description", recipientAddr, amount)
)

func setupTest(t *testing.T, maccBalances sdk.Coins) (ctx sdk.Context, app *marsapp.MarsApp, maccAddr sdk.AccAddress) {
	app = marsapptesting.MakeMockApp()
	ctx = app.BaseApp.NewContext(false, tmproto.Header{})

	maccAddr = app.SafetyFundKeeper.GetModuleAccount(ctx).GetAddress()

	// mint the specified amount of coins to safety fund module account
	app.BankKeeper.InitGenesis(
		ctx,
		&banktypes.GenesisState{
			Params: banktypes.Params{
				DefaultSendEnabled: true, // must set this to true so that tokens can be transferred
			},
			Balances: []banktypes.Balance{{
				Address: maccAddr.String(),
				Coins:   maccBalances,
			}},
		},
	)

	// verify the coins have been successfully minted
	balances := app.BankKeeper.GetAllBalances(ctx, maccAddr)
	require.Equal(t, maccBalances, balances)

	return ctx, app, maccAddr
}

// TestProposalHandlerPassed tests a case where the safety fund module account has a sufficient token
// balance, which should be successfully transferred to the recipient
func TestProposalHandlerPassed(t *testing.T) {
	// set up test by minting a sufficient amount of coins to the module account
	ctx, app, maccAddr := setupTest(t, amount)

	// the proposal should be executed with no error
	hdlr := safetyfund.NewProposalHandler(app.SafetyFundKeeper)
	require.NoError(t, hdlr(ctx, proposal))

	// the module account's balance should have been reduced to zero
	balances := app.BankKeeper.GetAllBalances(ctx, maccAddr)
	require.Equal(t, sdk.NewCoins(), balances)

	// the recipient should have received the correct amount
	balances = app.BankKeeper.GetAllBalances(ctx, recipientAddr)
	require.Equal(t, amount, balances)
}

// TestProposalHandlerFailed tests a case where the safetyfund module account does NOT have a sufficient
// token balance, which should result in an error when executing the proposal
func TestProposalHandlerFailed(t *testing.T) {
	maccBalances := sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(42069)))

	// set up test by minting an insufficient amount of coins to the module account
	ctx, app, maccAddr := setupTest(t, maccBalances)

	// attempt to execute the proposal without first giving safety fund module account any coin; should fail
	hdlr := safetyfund.NewProposalHandler(app.SafetyFundKeeper)
	require.Error(t, hdlr(ctx, proposal))

	// the module account's balance should NOT have changed
	balances := app.BankKeeper.GetAllBalances(ctx, maccAddr)
	require.Equal(t, maccBalances, balances)

	// the recipient should have NOT received any coin
	balances = app.BankKeeper.GetAllBalances(ctx, recipientAddr)
	require.Equal(t, sdk.NewCoins(), balances)
}
