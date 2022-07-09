package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"
)

func setupTest(t *testing.T) (ctx sdk.Context, app *marsapp.MarsApp, valcons sdk.ConsAddress, valoper sdk.ValAddress, user sdk.AccAddress) {
	app = marsapptesting.MakeMockApp()
	ctx = app.BaseApp.NewContext(false, tmproto.Header{})

	accts := marsapptesting.MakeRandomAccounts(3)
	valoper = sdk.ValAddress(accts[0])
	user = accts[1]
	feeSender := accts[2]

	pks := simapp.CreateTestPubKeys(1)
	valPubKey := pks[0]

	fees := sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(10000)), sdk.NewCoin("ibc/1234ABCD", sdk.NewInt(20000)))

	// register accounts at the auth module
	for _, acct := range accts {
		app.AccountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(acct))
	}

	// set coin balances for user and fee sender
	app.BankKeeper.InitGenesis(
		ctx,
		&banktypes.GenesisState{
			Params: banktypes.Params{
				DefaultSendEnabled: true,
			},
			Balances: []banktypes.Balance{{
				Address: user.String(),
				Coins:   sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(10))),
			}, {
				Address: feeSender.String(),
				Coins:   fees,
			}},
		},
	)

	// set bond denom to `umars`
	stakingParams := app.StakingKeeper.GetParams(ctx)
	stakingParams.BondDenom = "umars"
	app.StakingKeeper.SetParams(ctx, stakingParams)

	// create validator with bonded status and 20% commission
	val, err := stakingtypes.NewValidator(sdk.ValAddress(valoper), valPubKey, stakingtypes.Description{})
	val.Status = stakingtypes.Bonded
	val.Commission.CommissionRates = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(2, 1), sdk.NewDecWithPrec(2, 1), sdk.NewDec(0))
	require.NoError(t, err)
	require.True(t, val.IsBonded())

	app.StakingKeeper.SetValidator(ctx, val)
	app.StakingKeeper.SetValidatorByConsAddr(ctx, val)
	app.StakingKeeper.SetValidatorByPowerIndex(ctx, val)
	app.StakingKeeper.AfterValidatorCreated(ctx, val.GetOperator()) // required to initialize distr keeper properly

	// user makes delegation to validator
	newShares, err := app.StakingKeeper.Delegate(
		ctx,
		user,
		sdk.NewInt(10),
		stakingtypes.Unbonded,
		val,
		true, // true means it's a delegation, not a redelegation
	)
	require.NoError(t, err)
	require.True(t, newShares.GT(sdk.ZeroDec()))

	// initialize parameters of the distr module
	app.DistrKeeper.SetParams(ctx, distrtypes.Params{
		CommunityTax:        sdk.NewDecWithPrec(3, 1), // 30%
		BaseProposerReward:  sdk.ZeroDec(),
		BonusProposerReward: sdk.ZeroDec(), // set to zero for simplicity
	})

	// fee sender sends coins to the fee collector module account
	// this simulates the IBC module receives protocol revenue from outposts, and forward them to the fee collector
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, feeSender, authtypes.FeeCollectorName, fees)

	return ctx, app, sdk.ConsAddress(valPubKey.Address()), valoper, user
}

func calculateDelegationReward(t *testing.T, ctx sdk.Context, app *marsapp.MarsApp, valAddr sdk.ValAddress, delAddr sdk.AccAddress) sdk.DecCoins {
	// cache, we don't want to write changes
	ctx, _ = ctx.CacheContext()

	// query the validator
	val, found := app.StakingKeeper.GetValidator(ctx, valAddr)
	require.True(t, found)

	// query the delegation
	del, found := app.StakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	require.True(t, found)

	// increment the validator ending period
	endingPeriod := app.DistrKeeper.IncrementValidatorPeriod(ctx, val)

	return app.DistrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
}

func TestAllocateTokens(t *testing.T) {
	ctx, app, valcons, valoper, user := setupTest(t)

	app.DistrKeeper.AllocateTokens(
		ctx,
		10,
		10,
		valcons,
		[]abci.VoteInfo{{
			Validator: abci.Validator{
				Address: valcons,
				Power:   10,
			},
			SignedLastBlock: true,
		}},
	)

	// assert that distr module keeps track of the correct community pool balances
	// full amount of the 20000 ibc tokens should have been donated to the community pool
	// among the 10000 umars tokens, 30% or 3000 should have been donated to the community pool
	communityPool := app.DistrKeeper.GetFeePoolCommunityCoins(ctx)
	require.Equal(
		t,
		sdk.NewDecCoins(sdk.NewDecCoin("ibc/1234ABCD", sdk.NewInt(20000)), sdk.NewDecCoin("umars", sdk.NewInt(3000))),
		communityPool,
	)

	// assert that the actual balances of the community pool is correct
	// since no commission or reward has been withdrawn, the umars amount should be the whole 10000
	distrModAcctAddr := app.AccountKeeper.GetModuleAccount(ctx, distrtypes.ModuleName)
	distrModAcctBalances := app.BankKeeper.GetAllBalances(ctx, distrModAcctAddr.GetAddress())
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewCoin("ibc/1234ABCD", sdk.NewInt(20000)), sdk.NewCoin("umars", sdk.NewInt(10000))),
		distrModAcctBalances,
	)

	// amoung the 7000 umars after community tax, 20% or 1400 goes to validator commissions
	commission := app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valoper)
	require.Equal(t, sdk.NewDecCoins(sdk.NewDecCoin("umars", sdk.NewInt(1400))), commission.Commission)

	// the rest, or 5600 umars goes to staker's reward
	// NOTE: need to advance to the next block in order for the reward to register
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1) // advance 1 block
	reward := calculateDelegationReward(t, ctx, app, valoper, user)
	require.Equal(t, sdk.NewDecCoins(sdk.NewDecCoin("umars", sdk.NewInt(5600))), reward)
}
