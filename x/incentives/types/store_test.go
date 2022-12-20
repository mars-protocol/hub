package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/incentives/types"
)

var mockSchedule = types.Schedule{
	Id:             1,
	StartTime:      time.Unix(10000, 0),
	EndTime:        time.Unix(20000, 0),
	TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(12345)), sdk.NewCoin("uastro", sdk.NewInt(69420))),
	ReleasedAmount: sdk.NewCoins(),
}

func TestGetBlockRewardBeforeStart(t *testing.T) {
	blockReward := mockSchedule.GetBlockReward(time.Unix(5000, 0))
	require.Empty(t, blockReward)
}

func TestGetBlockRewardAfterEnd(t *testing.T) {
	blockReward := mockSchedule.GetBlockReward(time.Unix(42069, 0))
	expected := mockSchedule.TotalAmount.Sub(mockSchedule.ReleasedAmount...)
	require.Equal(t, expected, blockReward)
}

func TestGetBlockRewardBetweenStartAndEnd(t *testing.T) {
	// NOTE: the default number of decimals used by sdk.Dec is 18
	// umars:  12345 * 1e18 * 3333 / 10000 / 1e18 = 4114
	// uastro: 69420 * 1e18 * 3333 / 10000 / 1e18 = 23137
	blockReward := mockSchedule.GetBlockReward(time.Unix(13333, 0))
	expected := sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(4114)), sdk.NewCoin("uastro", sdk.NewInt(23137)))
	require.Equal(t, expected, blockReward)

	// next, try if the already released amount is properly subtracted
	mockSchedule.ReleasedAmount = blockReward

	// umars:  12345 * 1e18 * 8964 / 10000 / 1e18 - 4114  = 6952
	// uastro: 69420 * 1e18 * 8964 / 10000 / 1e18 - 23137 = 39091
	blockReward = mockSchedule.GetBlockReward(time.Unix(18964, 0))
	expected = sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(6952)), sdk.NewCoin("uastro", sdk.NewInt(39091)))
	require.Equal(t, expected, blockReward)
}
