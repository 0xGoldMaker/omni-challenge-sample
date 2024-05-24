package keeper

import (
	"omni/x/omni/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetWhitelist set a specific whitelist in the store from its index
func (k Keeper) SetWhitelist(ctx sdk.Context, whitelist types.Whitelist) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhitelistKeyPrefix))
	b := k.cdc.MustMarshal(&whitelist)
	store.Set(types.WhitelistKey(
		whitelist.Index,
	), b)
}

// GetWhitelist returns a whitelist from its index
func (k Keeper) GetWhitelist(
	ctx sdk.Context,
	index uint64,

) (val types.Whitelist, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhitelistKeyPrefix))

	b := store.Get(types.WhitelistKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveWhitelist removes a whitelist from the store
func (k Keeper) RemoveWhitelist(
	ctx sdk.Context,
	index uint64,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhitelistKeyPrefix))
	store.Delete(types.WhitelistKey(
		index,
	))
}

// GetAllWhitelist returns all whitelist
func (k Keeper) GetAllWhitelist(ctx sdk.Context) (list []types.Whitelist) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhitelistKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Whitelist
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
