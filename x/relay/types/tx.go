package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgRegisterAccount{}
)

func (msg MsgRegisterAccount) ValidateBasic() error {
	return nil
}

func (msg MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}
