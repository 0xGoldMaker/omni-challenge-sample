package keeper

import (
	"context"

	"omni/x/omni/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (k msgServer) Dewhitelisted(goCtx context.Context, msg *types.MsgDewhitelisted) (*types.MsgDewhitelistedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// Check if it is already registered.
	whitelisted := k.GetAllWhitelist(ctx)
	for _, w := range whitelisted {
		if w.Address == msg.Key {
			k.RemoveWhitelist(ctx, w.Index)

			// Emits event
			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.EventTypeDeWhitelistAdded),
					sdk.NewAttribute(sdk.AttributeKeySender, msg.Key),
				),
				sdk.NewEvent(
					types.EventTypeDeWhitelistAdded,
					sdk.NewAttribute(sdk.AttributeKeySender, msg.Key),
				),
			})

			return &types.MsgDewhitelistedResponse{Code: "200", Msg: "Succeed!"}, nil
		}
	}

	return &types.MsgDewhitelistedResponse{Code: "301", Msg: "Invalid address!"}, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "Not existing address")

}
