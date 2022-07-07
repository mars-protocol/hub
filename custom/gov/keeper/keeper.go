package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/mars-protocol/hub/custom/gov/types"
)

// Keeper defines the custom governance module Keeper
//
// NOTE: Keeper wraps the vanilla gov keeper to inherit most of its function. However, we include an
// additional dependency, the wasm keeper, which is needed for our custom vote tallying logic
type Keeper struct {
	govkeeper.Keeper

	wasmKeeper types.WasmKeeper
}

// NewKeeper returns a custom gov keeper
//
// NOTE: compared to the vanilla gov keeper's constructor function, here we require an additional
// wasm keeper, which is needed for our custom vote tallying logic
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace govtypes.ParamSubspace,
	authKeeper govtypes.AccountKeeper, bankKeeper govtypes.BankKeeper, sk govtypes.StakingKeeper,
	wasmKeeper types.WasmKeeper, rtr govtypes.Router,
) Keeper {
	return Keeper{
		Keeper:     govkeeper.NewKeeper(cdc, key, paramSpace, authKeeper, bankKeeper, sk, rtr),
		wasmKeeper: wasmKeeper,
	}
}
