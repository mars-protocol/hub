package types

import "cosmossdk.io/errors"

var (
	ErrInvalidProposalAmount    = errors.Register(ModuleName, 2, "invalid safety fund spend proposal amount")
	ErrInvalidProposalAuthority = errors.Register(ModuleName, 3, "invalid safety fund spend proposal authority")
	ErrInvalidProposalRecipient = errors.Register(ModuleName, 4, "invalid safety fund spend proposal recipient")
)
