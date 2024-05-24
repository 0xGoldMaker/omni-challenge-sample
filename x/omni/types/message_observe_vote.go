package types

import (
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgObserveVote = "observe_vote"

var _ sdk.Msg = &MsgObserveVote{}

func NewMsgObserveVote(creator string, round uint64, value sdk.Int, timestamp time.Time) *MsgObserveVote {
	return &MsgObserveVote{
		Creator:   creator,
		Round:     round,
		Value:     value,
		Timestamp: timestamp,
	}
}

func (msg *MsgObserveVote) Route() string {
	return RouterKey
}

func (msg *MsgObserveVote) Type() string {
	return TypeMsgObserveVote
}

func (msg *MsgObserveVote) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgObserveVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgObserveVote) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
