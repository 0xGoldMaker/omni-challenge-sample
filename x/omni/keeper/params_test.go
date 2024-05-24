package keeper_test

import (
	"testing"

	testkeeper "omni/testutil/keeper"
	"omni/x/omni/types"

	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.OmniKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, &params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
