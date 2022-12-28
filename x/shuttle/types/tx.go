package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
)

var (
	_ sdk.Msg = &MsgRegisterAccount{}
	_ sdk.Msg = &MsgSendFunds{}
	_ sdk.Msg = &MsgSendMessages{}

	// IMPORTANT: must implement this interface so that the GetCachedValue
	// method will work.
	//
	// docs:
	// https://docs.cosmos.network/main/core/encoding#interface-encoding-and-usage-of-any
	//
	// example in gov v1:
	// https://github.com/cosmos/cosmos-sdk/blob/v0.46.7/x/gov/types/v1/msgs.go#L97
	_ codectypes.UnpackInterfacesMessage = MsgSendMessages{}
)

//------------------------------------------------------------------------------
// MsgRegisterAccount
//------------------------------------------------------------------------------

func (m *MsgRegisterAccount) ValidateBasic() error {
	// the authority address must be valid
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrap(ErrInvalidProposalAuthority, err.Error())
	}

	return nil
}

func (m *MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	// we have already asserted that the authority address is valid in
	// ValidateBasic, so can ignore the error here
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

//------------------------------------------------------------------------------
// MsgSendFunds
//------------------------------------------------------------------------------

func (m *MsgSendFunds) ValidateBasic() error {
	// the authority address must be valid
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrap(ErrInvalidProposalAuthority, err.Error())
	}

	// the coins amount must be valid
	if err := m.Amount.Validate(); err != nil {
		return sdkerrors.Wrap(ErrInvalidProposalAmount, err.Error())
	}

	return nil
}

func (m *MsgSendFunds) GetSigners() []sdk.AccAddress {
	// we have already asserted that the authority address is valid in
	// ValidateBasic, so can ignore the error here
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

//------------------------------------------------------------------------------
// MsgSendMessages
//------------------------------------------------------------------------------

func (m *MsgSendMessages) ValidateBasic() error {
	// the authority address must be valid
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrap(ErrInvalidProposalAuthority, err.Error())
	}

	// the messages must each implement the sdk.Msg interface
	msgs, err := sdktx.GetMsgs(m.Messages, sdk.MsgTypeURL(m))
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidProposalMsg, err.Error())
	}

	// there must be at least one message
	if len(msgs) < 1 {
		return sdkerrors.Wrap(ErrInvalidProposalMsg, "proposal must contain at least one message")
	}

	// ideally, we want to check each message:
	//
	//  1. is valid (run msg.ValidateBasic)
	//  2. has only one signer
	//  3. this one signer is the interchain account
	//
	// unfortunately, these are not possible:
	//
	//  - for 1 and 2, the signer addresses has the host chain's bech prefix,
	//    this would cause ValidateBasic and GetSigners to fail, despite the
	//    message is perfectly valid.
	//  - for 3, this is a stateful check (we need to query the ICA's address)
	//    while in ValidateBasic we can only do stateless checks.
	return nil
}

func (m *MsgSendMessages) GetSigners() []sdk.AccAddress {
	// we have already asserted that the authority address is valid in
	// ValidateBasic, so can ignore the error here
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func (m MsgSendMessages) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return sdktx.UnpackInterfaces(unpacker, m.Messages)
}
