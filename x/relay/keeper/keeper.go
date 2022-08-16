package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	"github.com/mars-protocol/hub/x/relay/types"
)

// Keeper is the relay module's keeper
type Keeper struct {
	accountKeeper       types.AccountKeeper
	icaControllerKeeper types.ICAControllerKeeper
	scopedKeeper        capabilitykeeper.ScopedKeeper
}

// NewKeeper creates a new relay Keeper instance
func NewKeeper(accountKeeper types.AccountKeeper, icaControllerKeeper types.ICAControllerKeeper, scopedKeeper capabilitykeeper.ScopedKeeper) Keeper {
	// ensure incentives module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{accountKeeper, icaControllerKeeper, scopedKeeper}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetModuleAddress returns the incentives module account's address
func (k Keeper) GetModuleAddress() sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// ClaimCapability claims the channel capability passed via the OnOpenChanInit callback
func (k *Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}
