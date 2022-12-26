package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrAccountExists            = sdkerrors.Register(ModuleName, 2, "duplicate interchain account")
	ErrInvalidProposalAuthority = sdkerrors.Register(ModuleName, 3, "invalid shuttle proposal authority")
)
