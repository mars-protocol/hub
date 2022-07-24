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

		maccAddr := k.GetModuleAddress()
		actualTotal := k.bankKeeper.GetAllBalances(ctx, maccAddr)

		// NOTE: the actual amount does not necessarily need to be _exactly_ equal the expected amount.
		// we allow it as long as it's all greater or equal than expected.
		//
		// the reason is- if we assert exact equality, then anyone can cause the invariance to break
		// (and hence halt the chain) by sending some coins to the incentives module account.
		broken := actualTotal.IsAllGTE(expectedTotal)

		msg := sdk.FormatInvariant(
			types.ModuleName,
			"total-unreleased-incentives",
			fmt.Sprintf("\tsum of unreleased incentives: %s\n\tmodule account balances: %s", expectedTotal.String(), actualTotal.String()),
		)

		return msg, broken
	}
}
