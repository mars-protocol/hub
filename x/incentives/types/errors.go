package types

import "cosmossdk.io/errors"

var (
	ErrFailedRefundToCommunityPool     = errors.Register(ModuleName, 2, "failed to return funds to community pool")
	ErrFailedWithdrawFromCommunityPool = errors.Register(ModuleName, 3, "failed to withdraw funds from community pool")
	ErrInvalidProposalAmount           = errors.Register(ModuleName, 4, "invalid incentives proposal amount")
	ErrInvalidProposalAuthority        = errors.Register(ModuleName, 5, "invalid incentives proposal authority")
	ErrInvalidProposalIds              = errors.Register(ModuleName, 6, "invalid incentives proposal ids")
	ErrInvalidProposalStartEndTimes    = errors.Register(ModuleName, 7, "invalid incentives proposal start and end times")
)
