package omni

import (
	"omni/x/omni/keeper"
	"omni/x/omni/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the balance
	for _, elem := range genState.BalanceList {
		k.SetBalance(ctx, elem)
	}
	// Set all the whitelist
	for _, elem := range genState.WhitelistList {
		k.SetWhitelist(ctx, elem)
	}
	// Set all the observeVote
	for _, elem := range genState.ObserveVoteList {
		k.SetObserveVote(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, &genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.BalanceList = k.GetAllBalance(ctx)
	genesis.WhitelistList = k.GetAllWhitelist(ctx)
	genesis.ObserveVoteList = k.GetAllObserveVote(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
