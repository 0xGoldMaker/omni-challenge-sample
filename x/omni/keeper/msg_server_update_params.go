package keeper

import (
	"context"

	"omni/x/omni/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// get params
	params := k.Keeper.GetParams(ctx)

	// set param value except current round
	params.IsWhitelistEnabled = msg.Params.IsWhitelistEnabled
	params.MinConsensus = msg.Params.MinConsensus
	params.NumEpochs = msg.Params.NumEpochs
	params.ContractAddress = msg.Params.ContractAddress

	// update params
	k.Keeper.SetParams(ctx, &params)

	return &types.MsgUpdateParamsResponse{}, nil
}
