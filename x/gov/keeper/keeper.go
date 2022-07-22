package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// Keeper defines the custom governance module Keeper
//
// NOTE: Keeper wraps the vanilla gov keeper to inherit most of its functions. However, we include an
// additional dependency, the wasm keeper, which is needed for our custom vote tallying logic
type Keeper struct {
	govkeeper.Keeper

	storeKey sdk.StoreKey

	stakingKeeper govtypes.StakingKeeper // gov keeper has `sk` as a private field; we can't access it when tallying
	wasmKeeper    wasmtypes.ViewKeeper
}

// NewKeeper returns a custom gov keeper
//
// NOTE: compared to the vanilla gov keeper's constructor function, here we require an additional
// wasm keeper, which is needed for our custom vote tallying logic
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace govtypes.ParamSubspace,
	authKeeper govtypes.AccountKeeper, bankKeeper govtypes.BankKeeper, stakingKeeper govtypes.StakingKeeper,
	wasmKeeper wasmtypes.ViewKeeper, rtr govtypes.Router,
) Keeper {
	return Keeper{
		Keeper:        govkeeper.NewKeeper(cdc, key, paramSpace, authKeeper, bankKeeper, stakingKeeper, rtr),
		storeKey:      key,
		stakingKeeper: stakingKeeper,
		wasmKeeper:    wasmKeeper,
	}
}

// deleteVote deletes a vote from a given proposalID and voter from the store
//
// NOTE: the vanilla gov module does not make the `deleteVote` function public, so in order to delete
// votes, we need to redefine the function here
func (k Keeper) deleteVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(govtypes.VoteKey(proposalID, voterAddr))
}
