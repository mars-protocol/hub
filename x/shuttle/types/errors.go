package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrAccountExists   = sdkerrors.Register(ModuleName, 2, "interchain account already exists")
	ErrAccountNotFound = sdkerrors.Register(ModuleName, 3, "interchain account does not exist")
)
