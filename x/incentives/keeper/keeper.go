package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/incentives/types"
)

// Keeper is the incentives module's keeper
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey

	authKeeper    types.AccountKeeper
	bankKeeper    types.BankKeeper
	distrKeeper   types.DistrKeeper
	stakingKeeper types.StakingKeeper
}

// NewKeeper creates a new incentives module keeper
func NewKeeper(
	cdc codec.BinaryCodec, storeKey sdk.StoreKey, authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper, distrKeeper types.DistrKeeper, stakingKeeper types.StakingKeeper,
) Keeper {
	// ensure incentives module account is set
	if addr := authKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		authKeeper:    authKeeper,
		bankKeeper:    bankKeeper,
		distrKeeper:   distrKeeper,
		stakingKeeper: stakingKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetModuleAddress returns the incentives module account's address
func (k Keeper) GetModuleAddress(ctx sdk.Context) sdk.AccAddress {
	return k.authKeeper.GetModuleAccount(ctx, types.ModuleName).GetAddress()
}

//--------------------------------------------------------------------------------------------------
// ScheduleId
//--------------------------------------------------------------------------------------------------

// GetNextScheduleId loads the next schedule id if a new schedule is to be created
//
// NOTE: the id should have been initialized in genesis, so it being undefined is a fatal error. we
// have the module panic in this case, instead of returning an error
func (k Keeper) GetNextScheduleId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyNextScheduleId)
	if bz == nil {
		panic("stored next schedule id should not have been nil")
	}

	return sdk.BigEndianToUint64(bz)
}

// SetNextScheduleId sets the next schedule id to the provided value
func (k Keeper) SetNextScheduleId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyNextScheduleId, sdk.Uint64ToBigEndian(id))
}

// IncrementNextScheduleId increases the next id by one, and returns the incremented value
func (k Keeper) IncrementNextScheduleId(ctx sdk.Context) uint64 {
	id := k.GetNextScheduleId(ctx)

	id += 1
	k.SetNextScheduleId(ctx, id)

	return id
}

//--------------------------------------------------------------------------------------------------
// Schedule
//--------------------------------------------------------------------------------------------------

// GetSchedule loads the incentives schedule of the specified id
func (k Keeper) GetSchedule(ctx sdk.Context, id uint64) (schedule types.Schedule, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetScheduleKey(id))
	if bz == nil {
		return schedule, false
	}

	k.cdc.MustUnmarshal(bz, &schedule)

	return schedule, true
}

// SetSchedule saves the provided incentives schedule to store
func (k Keeper) SetSchedule(ctx sdk.Context, schedule types.Schedule) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetScheduleKey(schedule.Id), k.cdc.MustMarshal(&schedule))
}

// IterateSchedules iterates over all active schedules, calling the callback function with the schedule
// info. The iteration stops if the callback returns false.
func (k Keeper) IterateSchedules(ctx sdk.Context, cb func(types.Schedule) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeySchedule)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var schedule types.Schedule
		k.cdc.MustUnmarshal(iterator.Value(), &schedule)

		if cb(schedule) {
			break
		}
	}
}

// DeleteSchedule removes the incentives schedule of the given id from module store
func (k Keeper) DeleteSchedule(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetScheduleKey(id))
}

// GetSchedulePrefixStore returns a prefix store of all schedules
func (k Keeper) GetSchedulePrefixStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.KeySchedule)
}
