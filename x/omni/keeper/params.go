package keeper

import (
	"omni/x/omni/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefix(types.ParamsKey))
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(params)
	if err != nil {
		return err
	}
	store.Set(types.KeyPrefix(types.ParamsKey), bz)

	return nil
}
