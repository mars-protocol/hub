package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidProposalAmount    = sdkerrors.Register(ModuleName, 2, "invalid safety fund spend proposal amount")
	ErrInvalidProposalRecipient = sdkerrors.Register(ModuleName, 3, "invalid safety fund spend proposal recipient")
)
