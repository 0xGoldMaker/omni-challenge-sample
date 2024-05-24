package omni

import (
	"omni/x/omni/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker
// Check all requests from observers and update blockchain states
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// If it is genesis block, return
	if ctx.BlockHeight() < 1 {
		return
	}

	// Get params
	params := k.GetParams(ctx)

	// If it is the block height for updating balance
	if ctx.BlockHeight()%int64(params.NumEpochs) == 0 {
		// Update balance data using observation voted
		k.UpdateBalanceData(ctx, params.CurRound, params.MinConsensus)

		// Increase current round
		params.CurRound++

		// Update param by increasing current round
		k.SetParams(ctx, &params)
	}
}
