package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrAccountExists          = sdkerrors.Register(ModuleName, 2, "interchain account already exists")
	ErrAccountNotFound        = sdkerrors.Register(ModuleName, 3, "interchain account does not exist")
	ErrConnectionNotFound     = sdkerrors.Register(ModuleName, 4, "connect does not exist")
	ErrUnexpectedChannelOpen  = sdkerrors.Register(ModuleName, 5, "ICA channel open request should not be initialized from the host chain")
	ErrUnexpectedChannelClose = sdkerrors.Register(ModuleName, 6, "ICA channel is not expected to be closed")
)
