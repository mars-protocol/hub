package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidProposalAmount        = sdkerrors.Register(ModuleName, 2, "invalid incentives proposal amount")
	ErrInvalidProposalIds           = sdkerrors.Register(ModuleName, 3, "invalid incentives proposal ids")
	ErrInvalidProposalStartEndTimes = sdkerrors.Register(ModuleName, 4, "invalid incentives proposal start and end times")
)
