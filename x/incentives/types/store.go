package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func durationToSecondsDec(d time.Duration) sdk.Dec {
	return sdk.NewDecFromIntWithPrec(sdk.NewInt(d.Nanoseconds()), 9)
}

// GetBlockReward calculates the reward to be releaed given a time:
//   - if the current time is before the start time, no coin is to be released
//   - if the current time is after the end time, all coins are to be released
//   - if the current time is betweeen the start and end times, coins are to be
//     released linearly
func (s Schedule) GetBlockReward(currentTime time.Time) sdk.Coins {
	if s.StartTime.After(currentTime) {
		return sdk.NewCoins()
	}

	if currentTime.After(s.EndTime) {
		return s.TotalAmount.Sub(s.ReleasedAmount...)
	}

	timeTotal := durationToSecondsDec(s.EndTime.Sub(s.StartTime))
	timeElapsed := durationToSecondsDec(currentTime.Sub(s.StartTime))

	blockRewardDec := sdk.NewDecCoinsFromCoins(s.TotalAmount...).MulDec(timeElapsed).QuoDec(timeTotal)
	blockReward, _ := blockRewardDec.TruncateDecimal()

	return blockReward.Sub(s.ReleasedAmount...)
}
