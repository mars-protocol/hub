package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrFailedRefundToCommunityPool     = sdkerrors.Register(ModuleName, 2, "failed to return funds to community pool")
	ErrFailedWithdrawFromCommunityPool = sdkerrors.Register(ModuleName, 3, "failed to withdraw funds from community pool")
	ErrInvalidProposalAmount           = sdkerrors.Register(ModuleName, 4, "invalid incentives proposal amount")
	ErrInvalidProposalAuthority        = sdkerrors.Register(ModuleName, 5, "invalid incentives proposal authority")
	ErrInvalidProposalIds              = sdkerrors.Register(ModuleName, 6, "invalid incentives proposal ids")
	ErrInvalidProposalStartEndTimes    = sdkerrors.Register(ModuleName, 7, "invalid incentives proposal start and end times")
)
