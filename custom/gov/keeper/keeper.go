package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// Keeper defines the custom governance module Keeper
//
// NOTE: Keeper wraps the vanilla gov keeper to inherit most of its functions. However, we include an
// additional dependency, the wasm keeper, which is needed for our custom vote tallying logic
type Keeper struct {
	govkeeper.Keeper

	stakingKeeper govtypes.StakingKeeper // gov keeper has `sk` as a private field; we can't access it when tallying
	wasmKeeper    wasmtypes.ViewKeeper
}

// NewKeeper returns a custom gov keeper
//
// NOTE: compared to the vanilla gov keeper's constructor function, here we require an additional
// wasm keeper, which is needed for our custom vote tallying logic
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace govtypes.ParamSubspace,
	authKeeper govtypes.AccountKeeper, bankKeeper govtypes.BankKeeper, stakingKeeper govtypes.StakingKeeper,
	wasmKeeper wasmtypes.ViewKeeper, rtr govv1beta1.Router,
) Keeper {
	return Keeper{
		Keeper:        govkeeper.NewKeeper(cdc, key, paramSpace, authKeeper, bankKeeper, stakingKeeper, rtr),
		stakingKeeper: stakingKeeper,
		wasmKeeper:    wasmKeeper,
	}
}
