package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgRegisterAccount{}

//------------------------------------------------------------------------------
// MsgRegisterAccount
//------------------------------------------------------------------------------

// ValidateBasic does a sanity check on the provided data
func (m *MsgRegisterAccount) ValidateBasic() error {
	// the authority address must be valid
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return ErrInvalidProposalAuthority.Wrap(err.Error())
	}

	return nil
}

// GetSigners returns the expected signers for the message
func (m *MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	// we have already asserted that the authority address is valid in
	// ValidateBasic, so can ignore the error here
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}
