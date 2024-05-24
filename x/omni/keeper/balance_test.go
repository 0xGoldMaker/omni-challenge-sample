package keeper_test

import (
	"strconv"
	"testing"

	keepertest "omni/testutil/keeper"
	"omni/testutil/nullify"
	"omni/x/omni/keeper"
	"omni/x/omni/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNBalance(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Balance {
	items := make([]types.Balance, n)
	for i := range items {
		items[i].Index = (uint64)(i)
		items[i].Balance = sdk.NewInt(1)
		keeper.SetBalance(ctx, items[i])
	}
	return items
}

func TestBalanceGet(t *testing.T) {
	keeper, ctx := keepertest.OmniKeeper(t)
	items := createNBalance(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetBalance(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestBalanceRemove(t *testing.T) {
	keeper, ctx := keepertest.OmniKeeper(t)
	items := createNBalance(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveBalance(ctx,
			item.Index,
		)
		_, found := keeper.GetBalance(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestBalanceGetAll(t *testing.T) {
	keeper, ctx := keepertest.OmniKeeper(t)
	items := createNBalance(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllBalance(ctx)),
	)
}
