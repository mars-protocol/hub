package incentives

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	marsutils "github.com/mars-protocol/hub/utils"

	"github.com/mars-protocol/hub/x/incentives/keeper"
	"github.com/mars-protocol/hub/x/incentives/types"
)

// BeginBlocker distributes block rewards to validators who have signed the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	ids, totalBlockReward := k.ReleaseBlockReward(ctx, req.LastCommitInfo.Votes)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIncentivesReleased,
			sdk.NewAttribute(types.AttributeKeySchedules, marsutils.UintArrayToString(ids, ",")),
			sdk.NewAttribute(sdk.AttributeKeyAmount, totalBlockReward.String()),
		),
	)
}
