package keeper

import (
	"context"

	"omni/x/omni/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Check if the voter is whitelisted
func (k msgServer) CheckIfWhitelisted(ctx sdk.Context, voter string) bool {
	// Get all whitelisted nodes
	whitelisted := k.GetAllWhitelist(ctx)

	for _, w := range whitelisted {
		// If the address is listed in whitelisted.
		if w.Address == voter {
			return true
		}
	}

	return false
}

// Get observe voted for the current voter with roun if he has already voted
func (k msgServer) GetObserveVoted(ctx sdk.Context, voter string, round uint64) (types.ObserveVote, bool) {
	observeVotes := k.GetAllObserveVoteByRound(ctx, round)
	for _, ov := range observeVotes {
		if ov.Voter == voter {
			return ov, true
		}
	}

	return types.ObserveVote{}, false
}

// Process Observe vote
func (k msgServer) ObserveVote(goCtx context.Context, msg *types.MsgObserveVote) (*types.MsgObserveVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get params
	params := k.Keeper.GetParams(ctx)

	// Check if it is the vote for the current round
	if params.CurRound != msg.Round {
		return &types.MsgObserveVoteResponse{Code: "501", Msg: "Invalid round!"}, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid round")
	}

	// If whitelisting enabled then, check if it's a valid voter.
	if params.IsWhitelistEnabled && !k.CheckIfWhitelisted(ctx, msg.Creator) {
		return &types.MsgObserveVoteResponse{Code: "501", Msg: "Invalid request!"}, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid request")
	}

	// If if the proposer already voted, then replace it with the new one
	observeVote, found := k.GetObserveVoted(ctx, msg.Creator, msg.Round)
	if found {
		// Replace it with the new data
		observeVote.Timestamp = msg.Timestamp
		observeVote.Value = msg.Value

		// Save new observe vote data
		k.SetObserveVote(ctx, observeVote)
	} else {
		// Get next observation vote index
		n := k.GetObserveTxStoreCount(ctx) + 1

		// Initiate observe vote data
		observeVote = types.ObserveVote{
			Index:     n,
			Voter:     msg.Creator,
			Round:     msg.Round,
			Value:     msg.Value,
			Timestamp: msg.Timestamp,
		}

		// Save new observe vote data
		k.SetObserveVote(ctx, observeVote)

		// Update Requested Vote Count
		k.SetObserveTxStoreCount(ctx, n)
	}

	// Emits event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventTypeObservationVoted),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
		sdk.NewEvent(
			types.EventTypeObservationVoted,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgObserveVoteResponse{Code: "200", Msg: "Observation vote Succeed!"}, nil
}
