package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected interface for the auth module keeper
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
}

// ICAControllerKeeper defines the expected interface for the ICA controller keeper
type ICAControllerKeeper interface {
	RegisterInterchainAccount(ctx sdk.Context, connectionID, owner string) error
	GetInterchainAccountAddress(ctx sdk.Context, connectionID, portID string) (string, bool)
}
