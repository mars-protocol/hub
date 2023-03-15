package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/v2/x/incentives/types"
)

// there is no unit tests for keeper.go
//
// we simply define the mock variables here, while can be used by other test files

var mockGenesisState = types.GenesisState{
	NextScheduleId: 3,
	Schedules:      mockSchedules,
}

var mockSchedules = []types.Schedule{{
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
}}

var mockSchedulesReleased = []types.Schedule{{
	Id:             1,
	StartTime:      time.Unix(10000, 0),
	EndTime:        time.Unix(20000, 0),
	TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(12345)), sdk.NewCoin("uastro", sdk.NewInt(69420))),
	ReleasedAmount: sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(11066)), sdk.NewCoin("uastro", sdk.NewInt(62228))),
}, {
	Id:             2,
	StartTime:      time.Unix(15000, 0),
	EndTime:        time.Unix(30000, 0),
	TotalAmount:    sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(10000))),
	ReleasedAmount: sdk.NewCoins(sdk.NewCoin("umars", sdk.NewInt(2642))),
}}
