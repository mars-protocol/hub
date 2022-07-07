package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/custom/gov/keeper"
)

// EndBlocker called at the end of every block, processing proposals
func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	// TODO
}
