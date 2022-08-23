package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/keeper"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

// Keeper is the module's keeper
type Keeper struct {
	cdc        codec.Codec
	paramSpace paramtypes.Subspace

	accountKeeper       authkeeper.AccountKeeper
	icaControllerKeeper icacontrollerkeeper.Keeper
	scopedKeeper        capabilitykeeper.ScopedKeeper
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.Codec, paramSpace paramtypes.Subspace, accountKeeper authkeeper.AccountKeeper,
	icaControllerKeeper icacontrollerkeeper.Keeper, scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {
	// ensure incentives module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{cdc, paramSpace, accountKeeper, icaControllerKeeper, scopedKeeper}
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

// GetParams returns the module's parameters
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the module's parameters
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
