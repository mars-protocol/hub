package types

import "cosmossdk.io/errors"

var (
	ErrInvalidProposalAmount    = errors.Register(ModuleName, 2, "invalid envoy module proposal amount")
	ErrInvalidProposalAuthority = errors.Register(ModuleName, 3, "invalid envoy module proposal authority")
	ErrInvalidProposalMsg       = errors.Register(ModuleName, 4, "invalid envoy module proposal messages")
	ErrMultihopUnsupported      = errors.Register(ModuleName, 5, "multihop channels are not supported")
	ErrUnauthorized             = errors.Register(ModuleName, 6, "unauthorized")
)
