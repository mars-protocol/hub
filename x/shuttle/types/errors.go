package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrAccountExists            = sdkerrors.Register(ModuleName, 2, "interchain account already exists")
	ErrAccountNotFound          = sdkerrors.Register(ModuleName, 3, "interchain account not found")
	ErrChannelNotFound          = sdkerrors.Register(ModuleName, 4, "channel not found")
	ErrInvalidProposalAuthority = sdkerrors.Register(ModuleName, 5, "invalid shuttle module proposal authority")
	ErrMultihopUnsupported      = sdkerrors.Register(ModuleName, 6, "multihop channels are not supported")
)
