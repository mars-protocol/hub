package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidProposalAuthority = sdkerrors.Register(ModuleName, 2, "invalid shuttle proposal authority")
)
