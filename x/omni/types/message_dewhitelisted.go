package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDewhitelisted = "dewhitelisted"

var _ sdk.Msg = &MsgDewhitelisted{}

func NewMsgDewhitelisted(authority string, key string) *MsgDewhitelisted {
	return &MsgDewhitelisted{
		Authority: authority,
		Key:       key,
	}
}

func (msg *MsgDewhitelisted) Route() string {
	return RouterKey
}

func (msg *MsgDewhitelisted) Type() string {
	return TypeMsgDewhitelisted
}

func (msg *MsgDewhitelisted) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDewhitelisted) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDewhitelisted) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
