package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWhitelisted = "whitelisted"

var _ sdk.Msg = &MsgWhitelisted{}

func NewMsgWhitelisted(authority string, key string) *MsgWhitelisted {
	return &MsgWhitelisted{
		Authority: authority,
		Key:       key,
	}
}

func (msg *MsgWhitelisted) Route() string {
	return RouterKey
}

func (msg *MsgWhitelisted) Type() string {
	return TypeMsgWhitelisted
}

func (msg *MsgWhitelisted) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWhitelisted) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWhitelisted) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
