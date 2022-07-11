package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper define the custom distribution module keeper
//
// NOTE: Keeper wraps the vanilla distr keeper to inherit most of its functions. However, we replace
// the fee distribution logic with our own implementation
type Keeper struct {
	distrkeeper.Keeper

	// the vanilla distr keeper does not make public the auth and bank keepers, so we have to include
	// them here
	authKeeper distrtypes.AccountKeeper
	bankKeeper distrtypes.BankKeeper

	// same with the fee collector name, have to include it here
	feeCollectorName string
}

// NewKeeper returns a custom distr keeper
//
// NOTE: Keeper wraps the vanilla distr keeper to inherit most of its functions. However, we replace
// the fee distribution logic with our own implementation
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	ak distrtypes.AccountKeeper, bk distrtypes.BankKeeper, sk distrtypes.StakingKeeper,
	feeCollectorName string, blockedAddrs map[string]bool,
) Keeper {
	return Keeper{
		Keeper:           distrkeeper.NewKeeper(cdc, key, paramSpace, ak, bk, sk, feeCollectorName),
		authKeeper:       ak,
		bankKeeper:       bk,
		feeCollectorName: feeCollectorName,
	}
}
