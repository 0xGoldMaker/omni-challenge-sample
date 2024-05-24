package keeper

import (
	"context"

	"omni/x/omni/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (k msgServer) Whitelisted(goCtx context.Context, msg *types.MsgWhitelisted) (*types.MsgWhitelistedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// Check if it is already registered.
	whitelisted := k.GetAllWhitelist(ctx)
	for _, w := range whitelisted {
		if w.Address == msg.Key {
			return &types.MsgWhitelistedResponse{Code: "301", Msg: "Already registered!"}, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "Already registered")
		}
	}

	// Get Last saved index
	nIndex := k.GetStoreWhitelistedIndex(ctx) + 1

	// Create a new instance of Superadmin
	node := types.Whitelist{
		Index:   nIndex,
		Address: msg.Key,
	}

	// set whitelisted node
	k.SetWhitelist(ctx, node)

	// Update last saved index
	k.SetStoreWhitelistedIndex(ctx, nIndex)

	// Emits event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventTypeWhitelistAdded),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Key),
		),
		sdk.NewEvent(
			types.EventTypeWhitelistAdded,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Key),
		),
	})

	return &types.MsgWhitelistedResponse{Code: "200", Msg: "Succeed!"}, nil
}
