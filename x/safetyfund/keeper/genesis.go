package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/safetyfund/types"
)

// InitGenesis initializes the safetyfund module's storage according to the provided genesis state
//
// NOTE: we call `GetModuleAccount` instead of `SetModuleAccount` because the "get" function automatically
// sets the module account if it doesn't exist
func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	k.authKeeper.GetModuleAccount(ctx, types.ModuleAccountName)
}

// ExportGenesis returns a genesis state for a given context and keeper
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{}
}
