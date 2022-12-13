package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/incentives/types"
)

// InitGenesis initializes the incentives module's storage according to the provided genesis state
//
// NOTE: we call `GetModuleAccount` instead of `SetModuleAccount` because the "get" function automatically
// sets the module account if it doesn't exist
func (k Keeper) InitGenesis(ctx sdk.Context, gs *types.GenesisState) {
	// set module account
	k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)

	// set incentives schedules
	for _, schedule := range gs.Schedules {
		k.SetSchedule(ctx, schedule)
	}

	// set next schedule id
	k.SetNextScheduleID(ctx, gs.NextScheduleId)
}

// ExportGenesis returns a genesis state for a given context and keeper
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	nextScheduleID := k.GetNextScheduleID(ctx)

	schedules := []types.Schedule{}
	k.IterateSchedules(ctx, func(schedule types.Schedule) bool {
		schedules = append(schedules, schedule)
		return false
	})

	return &types.GenesisState{
		NextScheduleId: nextScheduleID,
		Schedules:      schedules,
	}
}
