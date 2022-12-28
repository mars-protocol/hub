package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrAccountExists            = sdkerrors.Register(ModuleName, 2, "interchain account already exists")
	ErrInvalidProposalAmount    = sdkerrors.Register(ModuleName, 3, "invalid shuttle module proposal amount")
	ErrInvalidProposalAuthority = sdkerrors.Register(ModuleName, 4, "invalid shuttle module proposal authority")
	ErrInvalidProposalMsg       = sdkerrors.Register(ModuleName, 5, "invalid shuttle module proposal messages")
	ErrMultihopUnsupported      = sdkerrors.Register(ModuleName, 6, "multihop channels are not supported")
)
