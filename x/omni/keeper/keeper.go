package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"omni/x/omni/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		authority  string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	authority string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		authority:  authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Read store to get previous block count
func (k Keeper) GetStoreWhitelistedIndex(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.WhitelistedNodeCountStoreKey)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

// Write store the current block count
func (k Keeper) SetStoreWhitelistedIndex(ctx sdk.Context, count uint64) bool {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(count)
	store.Set(types.WhitelistedNodeCountStoreKey, bz)

	return true
}

// Read requested observe data count
func (k Keeper) GetObserveTxStoreCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ObserveTxCountStoreKey)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

// Write requested observe data count
func (k Keeper) SetObserveTxStoreCount(ctx sdk.Context, count uint64) bool {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(count)
	store.Set(types.ObserveTxCountStoreKey, bz)

	return true
}

// Update balance data using observation voting
func (k Keeper) UpdateBalanceData(ctx sdk.Context, round uint64, minConsensusCnt uint64) error {
	observeVotes := k.GetAllObserveVoteByRound(ctx, round)

	// if there is no observe votes, return
	if len(observeVotes) < 1 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "no observations available")
	}

	// total balance voted which is then divided by the voted count
	total := sdk.NewInt(0)
	// total voted
	voted := (int64)(len(observeVotes))
	// Initially set the time
	lastUpdated := observeVotes[0].Timestamp

	// Loop observe votes to calculate the balance
	for _, ov := range observeVotes {
		// calc total sum of the value assuing each of the vote has the same weight
		total = total.Add(ov.Value)

		// If it is old time, then update it with a newer one
		if lastUpdated.Before(ov.Timestamp) {
			lastUpdated = ov.Timestamp
		}
	}

	// if the voted count is smaller than min consensus count - min validators that needed to make consensus
	if voted < (int64)(minConsensusCnt) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "in sufficient observation voted")
	}

	// averaged balance by diving the total balance with the voted count
	// TOTO:
	// 1. Consider weight of each voter
	// 2. Consider 2/3+ voter of the total validators
	// 3. Apply slahsing logic to give penalty to the validators who hasn't attened the balance voting
	avgBalance := total.Quo(sdk.NewInt(voted))

	// initiate the new balance entity
	newBalance := types.Balance{
		Index:       0,
		Balance:     avgBalance,
		LastUpdated: lastUpdated,
	}

	// Update the balance
	k.SetBalance(ctx, newBalance)

	return nil
}
