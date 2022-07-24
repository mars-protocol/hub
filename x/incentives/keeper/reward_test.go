package keeper_test

import (
	"testing"
	"time"

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

	"github.com/mars-protocol/hub/x/incentives/types"
)

//--------------------------------------------------------------------------------------------------
// Test Suite
//--------------------------------------------------------------------------------------------------

type testSuite struct {
	t   *testing.T
	ctx sdk.Context
	app *marsapp.MarsApp

	delAddr     sdk.AccAddress
	valAddr     sdk.ValAddress
	valconsAddr sdk.ConsAddress
}

func (suite *testSuite) setBlockHeight(height int64) {
	suite.ctx = suite.ctx.WithBlockHeight(height)
}

func (suite *testSuite) setBlockTime(sec int64) {
	suite.ctx = suite.ctx.WithBlockTime(time.Unix(sec, 0))
}

func (suite *testSuite) releaseBlockReward() (ids []uint64, totalBlockReward sdk.Coins) {
	mockBondedVotes := []abci.VoteInfo{{
		Validator: abci.Validator{
			Address: suite.valconsAddr,
			Power:   10,
		},
		SignedLastBlock: true,
	}}

	return suite.app.IncentivesKeeper.ReleaseBlockReward(suite.ctx, mockBondedVotes)
}

func (suite *testSuite) calculateDelegationReward() sdk.DecCoins {
	// cache the context; we don't want to write changes
	ctx, _ := suite.ctx.CacheContext()

	// query the validator
	val, found := suite.app.StakingKeeper.GetValidator(ctx, suite.valAddr)
	require.True(suite.t, found)
	suite.t.Log("val:", val)

	// query the delegation
	del, found := suite.app.StakingKeeper.GetDelegation(ctx, suite.delAddr, suite.valAddr)
	require.True(suite.t, found)
	suite.t.Log("del:", del)

	// increment the validator ending period
	endingPeriod := suite.app.DistrKeeper.IncrementValidatorPeriod(ctx, val)
	suite.t.Log("endingPeriod:", endingPeriod)

	return suite.app.DistrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
}

//--------------------------------------------------------------------------------------------------
// Test Setup
//--------------------------------------------------------------------------------------------------

func setupRewardTest(t *testing.T, schedules []types.Schedule) *testSuite {
	app := marsapptesting.MakeMockApp()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	accts := marsapptesting.MakeRandomAccounts(2)
	delAddr := accts[0]
	valAddr := sdk.ValAddress(accts[1])

	pks := simapp.CreateTestPubKeys(1)
	valPubKey := pks[0]
	valconsAddr := sdk.ConsAddress(valPubKey.Address())

	//----------------------------------------
	// set up auth module

	// register accounts at the auth module
	for _, acct := range accts {
		app.AccountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(acct))
	}

	//----------------------------------------
	// set up bank module

	// at the same time, calculate the total mars token amount needed to be given to incentives module account
	totalIncentives := sdk.NewCoins()
	for _, schedule := range schedules {
		totalIncentives = totalIncentives.Add(schedule.TotalAmount...)
	}

	// set mars token balance for the staker and the incentives module account
	maccAddr := app.IncentivesKeeper.GetModuleAddress(ctx)
	app.BankKeeper.InitGenesis(
		ctx,
		&banktypes.GenesisState{
			Params: banktypes.Params{
				DefaultSendEnabled: true, // must set this to true so that tokens can be transferred
			},
			Balances: []banktypes.Balance{{
				Address: maccAddr.String(),
				Coins:   totalIncentives,
			}, {
				Address: delAddr.String(),
				Coins:   sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(10))),
			}},
		},
	)

	//----------------------------------------
	// set up staking module

	// set bond denom to `umars`
	stakingParams := app.StakingKeeper.GetParams(ctx)
	stakingParams.BondDenom = "umars"
	app.StakingKeeper.SetParams(ctx, stakingParams)

	// create validator. for simplicity, we choose zero commission
	val, err := stakingtypes.NewValidator(valAddr, valPubKey, stakingtypes.Description{})
	val.Status = stakingtypes.Bonded
	val.Commission.CommissionRates = stakingtypes.NewCommissionRates(sdk.NewDec(0), sdk.NewDec(0), sdk.NewDec(0))
	require.NoError(t, err)
	require.True(t, val.IsBonded())

	// save validator info in staking module store
	app.StakingKeeper.SetValidator(ctx, val)
	app.StakingKeeper.SetValidatorByConsAddr(ctx, val)
	app.StakingKeeper.SetValidatorByPowerIndex(ctx, val)
	app.StakingKeeper.AfterValidatorCreated(ctx, val.GetOperator()) // required to initialize distr keeper properly

	// user makes delegation to validator
	newShares, err := app.StakingKeeper.Delegate(
		ctx,
		delAddr,
		sdk.NewInt(10),
		stakingtypes.Unbonded,
		val,
		true, // true means it's a delegation, not a redelegation
	)
	require.NoError(t, err)
	require.True(t, newShares.GT(sdk.ZeroDec()))

	//----------------------------------------
	// set up distr module

	// initialize parameters of the distr module
	app.DistrKeeper.SetParams(ctx, distrtypes.Params{
		CommunityTax:        sdk.ZeroDec(),
		BaseProposerReward:  sdk.ZeroDec(),
		BonusProposerReward: sdk.ZeroDec(),
	})

	//----------------------------------------
	// set up incentives module

	// save incentives schedules
	for _, schedule := range schedules {
		app.IncentivesKeeper.SetSchedule(ctx, schedule)
	}

	return &testSuite{t, ctx, app, delAddr, valAddr, valconsAddr}
}

//--------------------------------------------------------------------------------------------------
// Tests
//--------------------------------------------------------------------------------------------------

func TestNoActiveSchedule(t *testing.T) {
	suite := setupRewardTest(t, []types.Schedule{})

	ids, blockReward := suite.releaseBlockReward()
	require.Empty(t, ids)
	require.Equal(t, sdk.NewCoins(), blockReward)
}

func TestBeforeStartTime(t *testing.T) {
	suite := setupRewardTest(
		t,
		[]types.Schedule{{
			Id:             1,
			StartTime:      time.Unix(10000, 0),
			EndTime:        time.Unix(20000, 0),
			TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(12345)), sdk.NewCoin("uastro", sdk.NewInt(69420))),
			ReleasedAmount: sdk.NewCoins(),
		}},
	)

	// set time to 1 sec before the schedule starsuite. no token should be released
	suite.setBlockTime(9999)

	ids, blockReward := suite.releaseBlockReward()
	require.Empty(t, ids)
	require.Equal(t, sdk.NewCoins(), blockReward)
}

func TestTwoActiveSchedules(t *testing.T) {
	suite := setupRewardTest(
		t,
		[]types.Schedule{{
			Id:             1,
			StartTime:      time.Unix(10000, 0),
			EndTime:        time.Unix(20000, 0),
			TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(12345)), sdk.NewCoin("uastro", sdk.NewInt(69420))),
			ReleasedAmount: sdk.NewCoins(),
		}, {
			Id:             2,
			StartTime:      time.Unix(15000, 0),
			EndTime:        time.Unix(30000, 0),
			TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(10000))),
			ReleasedAmount: sdk.NewCoins(),
		}},
	)

	//----------------------------------------
	// part 1

	// set time to 13333
	// schedule 1 should release 4114 umars + 23137 uastro (see reward_test.go for calculation)
	// schedule 2 should release nothing
	suite.setBlockHeight(1)
	suite.setBlockTime(13333)
	expectedBlockReward := sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(4114)), sdk.NewCoin("uastro", sdk.NewInt(23137)))

	ids, blockReward := suite.releaseBlockReward()
	require.Equal(t, []uint64{1}, ids)
	require.Equal(t, expectedBlockReward, blockReward)

	// expected delegation reward should be equal to this block's block reward
	//
	// NOTE: need to advance to the next block in order for the reward to register (!!!)
	suite.setBlockHeight(2)
	expectedDelReward := sdk.NewDecCoinsFromCoins(expectedBlockReward...)

	delegationReward := suite.calculateDelegationReward()
	require.Equal(t, expectedDelReward, delegationReward)

	//----------------------------------------
	// part 2

	// set time to 18964
	// schedule 1 should release 6952 umars + 39091 uastro (see reward_test.go for calculation)
	// schedule 2 should release 10000 * 1e18 * 3964 / 15000 / 1e18 = 2642 umars
	// total: 9594 umars + 39091 uastro
	suite.setBlockTime(18964)
	expectedBlockReward = sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(9594)), sdk.NewCoin("uastro", sdk.NewInt(39091)))

	ids, blockReward = suite.releaseBlockReward()
	require.Equal(t, []uint64{1, 2}, ids)
	require.Equal(t, expectedBlockReward, blockReward)

	// expected delegation reward should be the sum of the previous two
	suite.setBlockHeight(3)
	expectedDelReward = expectedDelReward.Add(sdk.NewDecCoinsFromCoins(expectedBlockReward...)...)

	delegationReward = suite.calculateDelegationReward()
	require.Equal(t, expectedDelReward, delegationReward)
}

func TestDeleteEndedSchedules(t *testing.T) {
	suite := setupRewardTest(
		t,
		[]types.Schedule{{
			Id:             1,
			StartTime:      time.Unix(10000, 0),
			EndTime:        time.Unix(20000, 0),
			TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(12345)), sdk.NewCoin("uastro", sdk.NewInt(69420))),
			ReleasedAmount: sdk.NewCoins(),
		}, {
			Id:             2,
			StartTime:      time.Unix(15000, 0),
			EndTime:        time.Unix(30000, 0),
			TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(10000))),
			ReleasedAmount: sdk.NewCoins(),
		}},
	)

	ctx, keeper := suite.ctx, &suite.app.IncentivesKeeper

	suite.setBlockHeight(1)
	suite.setBlockTime(20001)

	_, _ = suite.releaseBlockReward()

	// schedule 1 should have been deleted
	_, found := keeper.GetSchedule(ctx, 1)
	require.False(t, found)

	// schedule 2 should NOT have been deleted
	_, found = keeper.GetSchedule(ctx, 2)
	require.True(t, found)
}
