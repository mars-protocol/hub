package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/incentives/types"
)

// RegisterInvariants registers the incentives module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "total-unreleased-incentives", TotalUnreleasedIncentives(k))
}

// TotalUnreleasedIncentives asserts that the incentives module's coin balances match exactly the total
// amount of unreleased incentives
func TotalUnreleasedIncentives(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		expectedTotal := sdk.NewCoins()
		k.IterateSchedules(ctx, func(schedule types.Schedule) bool {
			expectedTotal = expectedTotal.Add(schedule.TotalAmount.Sub(schedule.ReleasedAmount)...)
			return false
		})

		maccAddr := k.GetModuleAccount(ctx).GetAddress()
		actualTotal := k.bankKeeper.GetAllBalances(ctx, maccAddr)

		broken := !expectedTotal.IsEqual(actualTotal)

		msg := sdk.FormatInvariant(
			types.ModuleName,
			"total-unreleased-incentives",
			fmt.Sprintf("\tsum of unreleased incentives: %s\n\tmodule account balances: %s\n", expectedTotal.String(), actualTotal.String()),
		)

		return msg, broken
	}
}
