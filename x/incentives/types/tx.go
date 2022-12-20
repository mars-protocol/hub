package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateSchedule{}
	_ sdk.Msg = &MsgTerminateSchedules{}
)

//------------------------------------------------------------------------------
// MsgCreateSchedule
//------------------------------------------------------------------------------

// ValidateBasic does a sanity check on the provided data
func (m *MsgCreateSchedule) ValidateBasic() error {
	// the authority address must be valid
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return ErrInvalidProposalAuthority.Wrap(err.Error())
	}

	// The start time must be earlier (strictly less than) the end time
	if !m.StartTime.Before(m.EndTime) {
		return ErrInvalidProposalStartEndTimes
	}

	// the coins must be valid (unique denoms, non-zero amount, and sorted
	// alphabetically)
	if m.Amount.Empty() {
		return ErrInvalidProposalAmount
	}

	return nil
}

// GetSigners returns the expected signers for the message
func (m *MsgCreateSchedule) GetSigners() []sdk.AccAddress {
	// we have already asserted that the authority address is valid in
	// ValidateBasic, so can ignore the error here
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

//------------------------------------------------------------------------------
// MsgTerminateSchedules
//------------------------------------------------------------------------------

// ValidateBasic does a sanity check on the provided data
func (m *MsgTerminateSchedules) ValidateBasic() error {
	// the authority address must be valid
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return ErrInvalidProposalAuthority.Wrap(err.Error())
	}

	// there must be at least one schedule id to be terminated
	if len(m.Ids) == 0 {
		return ErrInvalidProposalIds
	}

	return nil
}

// GetSigners returns the expected signers for the message
func (m *MsgTerminateSchedules) GetSigners() []sdk.AccAddress {
	// we have already asserted that the authority address is valid in
	// ValidateBasic, so can ignore the error here
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}
