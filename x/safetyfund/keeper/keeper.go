package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/mars-protocol/hub/x/safetyfund/types"
)

// Keeper is the safetyfund module's keeper
type Keeper struct {
	authKeeper types.AccountKeeper
	bankKeeper types.BankKeeper
}

// NewKeeper creates a new safetyfund Keeper instance
func NewKeeper(authKeeper types.AccountKeeper, bankKeeper types.BankKeeper) Keeper {
	// ensure the module account is set
	if authKeeper.GetModuleAddress(types.ModuleName) == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{authKeeper, bankKeeper}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("")
}

// GetModuleAccount returns the safety fund module account
func (k Keeper) GetModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetBalances returns the amount of coins available in the safety fund
func (k Keeper) GetBalances(ctx sdk.Context) sdk.Coins {
	return k.bankKeeper.GetAllBalances(ctx, k.GetModuleAccount(ctx).GetAddress())
}

// ReleaseFund releases coins from the safety fund to the specified recipient
func (k Keeper) ReleaseFund(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, amount)
}
