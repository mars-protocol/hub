package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/keeper"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

// Keeper is the shuttle module's keeper.
type Keeper struct {
	accountKeeper       authkeeper.AccountKeeper
	scopedKeeper        capabilitykeeper.ScopedKeeper
	icaControllerKeeper icacontrollerkeeper.Keeper

	authority string
}

// NewKeeper creates a new shuttle module keeper.
func NewKeeper(
	accountKeeper authkeeper.AccountKeeper, scopedKeeper capabilitykeeper.ScopedKeeper,
	icaControllerKeeper icacontrollerkeeper.Keeper, authority string,
) Keeper {
	// ensure shuttle module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		accountKeeper:       accountKeeper,
		scopedKeeper:        scopedKeeper,
		icaControllerKeeper: icaControllerKeeper,
		authority:           authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetModuleAddress returns the shuttle module account's address
func (k Keeper) GetModuleAddress() sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}
