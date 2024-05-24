package keeper

import (
	"context"

	"omni/x/omni/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) WhitelistAll(goCtx context.Context, req *types.QueryAllWhitelistRequest) (*types.QueryAllWhitelistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var whitelists []types.Whitelist
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	whitelistStore := prefix.NewStore(store, types.KeyPrefix(types.WhitelistKeyPrefix))

	pageRes, err := query.Paginate(whitelistStore, req.Pagination, func(key []byte, value []byte) error {
		var whitelist types.Whitelist
		if err := k.cdc.Unmarshal(value, &whitelist); err != nil {
			return err
		}

		whitelists = append(whitelists, whitelist)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllWhitelistResponse{Whitelist: whitelists, Pagination: pageRes}, nil
}

func (k Keeper) Whitelist(goCtx context.Context, req *types.QueryGetWhitelistRequest) (*types.QueryGetWhitelistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetWhitelist(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetWhitelistResponse{Whitelist: val}, nil
}
