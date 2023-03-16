package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/v2/x/safety/types"
)

// Keeper is the module's keeper
type Keeper struct {
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	authority     string
}

// NewKeeper creates a new Keeper instance
func NewKeeper(accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, authority string) Keeper {
	// ensure the module account is set
	if accountKeeper.GetModuleAddress(types.ModuleName) == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{accountKeeper, bankKeeper, authority}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetModuleAddress returns the safety fund module account's address
func (k Keeper) GetModuleAddress() sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// GetBalances returns the amount of coins available in the safety fund
func (k Keeper) GetBalances(ctx sdk.Context) sdk.Coins {
	return k.bankKeeper.GetAllBalances(ctx, k.GetModuleAddress())
}

// ReleaseFund releases coins from the safety fund to the specified recipient
func (k Keeper) ReleaseFund(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, amount)
}
