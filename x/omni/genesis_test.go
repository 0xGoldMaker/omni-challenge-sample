package omni_test

import (
	"testing"

	keepertest "omni/testutil/keeper"
	"omni/testutil/nullify"
	"omni/x/omni"
	"omni/x/omni/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		BalanceList: []types.Balance{
			{
				Index: 0,
			},
			{
				Index: 1,
			},
		},
		WhitelistList: []types.Whitelist{
			{
				Index: 0,
			},
			{
				Index: 1,
			},
		},
		ObserveVoteList: []types.ObserveVote{
			{
				Index: 0,
			},
			{
				Index: 1,
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.OmniKeeper(t)
	omni.InitGenesis(ctx, *k, genesisState)
	got := omni.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.BalanceList, got.BalanceList)
	require.ElementsMatch(t, genesisState.WhitelistList, got.WhitelistList)
	require.ElementsMatch(t, genesisState.ObserveVoteList, got.ObserveVoteList)
	// this line is used by starport scaffolding # genesis/test/assert
}
