package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	marsapp "github.com/mars-protocol/hub/v2/app"
	marsapptesting "github.com/mars-protocol/hub/v2/app/testing"

	"github.com/mars-protocol/hub/v2/x/incentives/types"
)

//--------------------------------------------------------------------------------------------------
// Test Suite
//--------------------------------------------------------------------------------------------------

// In this test, for simplicity, we assume
// - There is only one validator who has self-bond, and no other delegator
// - community tax rate is zero
type testSuite struct {
	t           *testing.T
	ctx         sdk.Context
	app         *marsapp.MarsApp
	validator   sdk.AccAddress
	valConsAddr sdk.ConsAddress
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
			Address: suite.valConsAddr,
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
	val, found := suite.app.StakingKeeper.GetValidator(ctx, sdk.ValAddress(suite.validator))
	require.True(suite.t, found)
	suite.t.Log("val:", val)

	// query the delegation
	del, found := suite.app.StakingKeeper.GetDelegation(ctx, suite.validator, sdk.ValAddress(suite.validator))
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
	accts := marsapptesting.MakeRandomAccounts(1)
	validator := accts[0]
	maccAddr := authtypes.NewModuleAddress(types.ModuleName)

	// calculate the total mars token amount needed to be given to incentives
	// module account
	totalIncentives := sdk.NewCoins()
	for _, schedule := range schedules {
		totalIncentives = totalIncentives.Add(schedule.TotalAmount...)
	}

	app := marsapptesting.MakeMockApp(
		accts,
		[]banktypes.Balance{{
			Address: maccAddr.String(),
			Coins:   totalIncentives,
		}},
		accts,
		sdk.NewCoins(),
	)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// set community tax rate to zero for convenience of this test
	app.DistrKeeper.SetParams(ctx, distrtypes.Params{
		CommunityTax:        sdk.ZeroDec(),
		BaseProposerReward:  sdk.ZeroDec(),
		BonusProposerReward: sdk.ZeroDec(),
	})

	// save incentives schedules
	for _, schedule := range schedules {
		app.IncentivesKeeper.SetSchedule(ctx, schedule)
	}

	// query the validator's consensus address
	val, found := app.StakingKeeper.GetValidator(ctx, sdk.ValAddress(validator))
	if !found {
		panic("Validator with address not found")
	}

	valCondAddr, err := val.GetConsAddr()
	if err != nil {
		panic(err)
	}

	return &testSuite{t, ctx, app, validator, valCondAddr}
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
	suite := setupRewardTest(t, mockSchedules)

	// set time to 1 sec before the schedule starsuite. no token should be released
	suite.setBlockTime(9999)

	ids, blockReward := suite.releaseBlockReward()
	require.Empty(t, ids)
	require.Equal(t, sdk.NewCoins(), blockReward)
}

func TestTwoActiveSchedules(t *testing.T) {
	suite := setupRewardTest(t, mockSchedules)

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
	suite := setupRewardTest(t, mockSchedules)

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
