package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"omni/x/omni/types"
)

func (k Keeper) ObserveVoteAll(goCtx context.Context, req *types.QueryAllObserveVoteRequest) (*types.QueryAllObserveVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var observeVotes []types.ObserveVote
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	observeVoteStore := prefix.NewStore(store, types.KeyPrefix(types.ObserveVoteKeyPrefix))

	pageRes, err := query.Paginate(observeVoteStore, req.Pagination, func(key []byte, value []byte) error {
		var observeVote types.ObserveVote
		if err := k.cdc.Unmarshal(value, &observeVote); err != nil {
			return err
		}

		observeVotes = append(observeVotes, observeVote)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllObserveVoteResponse{ObserveVote: observeVotes, Pagination: pageRes}, nil
}

func (k Keeper) ObserveVote(goCtx context.Context, req *types.QueryGetObserveVoteRequest) (*types.QueryGetObserveVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetObserveVote(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetObserveVoteResponse{ObserveVote: val}, nil
}
