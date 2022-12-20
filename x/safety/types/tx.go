package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSafetyFundSpend{}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgSafetyFundSpend) ValidateBasic() error {
	// the authority address must be valid
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return ErrInvalidProposalAuthority.Wrap(err.Error())
	}

	// the recipeint address must be valid
	if _, err := sdk.AccAddressFromBech32(m.Recipient); err != nil {
		return ErrInvalidProposalRecipient.Wrap(err.Error())
	}

	// the coins must be valid (unique denoms, non-zero amount, and sorted
	// alphabetically)
	if !m.Amount.IsValid() {
		return ErrInvalidProposalAmount
	}

	return nil
}

// GetSigners returns the expected signers for the message
func (m *MsgSafetyFundSpend) GetSigners() []sdk.AccAddress {
	// we have already asserted that the authority address is valid in
	// ValidateBasic, so can ignore the error here
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}
