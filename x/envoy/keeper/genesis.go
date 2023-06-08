package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/v2/x/envoy/types"
)

// InitGenesis initializes the envoy module's storage according to the
// provided genesis state.
//
// NOTE: we call `GetModuleAccount` instead of `SetModuleAccount` because the
// "get" function automatically sets the module account if it doesn't exist.
func (k Keeper) InitGenesis(ctx sdk.Context, _ *types.GenesisState) {
	// set module account
	k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// ExportGenesis returns a genesis state for a given context and keeper.
func (k Keeper) ExportGenesis(_ sdk.Context) *types.GenesisState {
	return &types.GenesisState{}
}
