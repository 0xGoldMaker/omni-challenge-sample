package omni

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"omni/testutil/sample"
	omnisimulation "omni/x/omni/simulation"
	"omni/x/omni/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = omnisimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgWhitelisted = "op_weight_msg_whitelisted"
	// TODO: Determine the simulation weight value
	defaultWeightMsgWhitelisted int = 100

	opWeightMsgDewhitelisted = "op_weight_msg_dewhitelisted"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDewhitelisted int = 100

	opWeightMsgObserveVote = "op_weight_msg_observe_vote"
	// TODO: Determine the simulation weight value
	defaultWeightMsgObserveVote int = 100

	opWeightMsgUpdateParams = "op_weight_msg_update_params"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateParams int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	omniGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&omniGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgWhitelisted int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgWhitelisted, &weightMsgWhitelisted, nil,
		func(_ *rand.Rand) {
			weightMsgWhitelisted = defaultWeightMsgWhitelisted
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgWhitelisted,
		omnisimulation.SimulateMsgWhitelisted(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDewhitelisted int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDewhitelisted, &weightMsgDewhitelisted, nil,
		func(_ *rand.Rand) {
			weightMsgDewhitelisted = defaultWeightMsgDewhitelisted
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDewhitelisted,
		omnisimulation.SimulateMsgDewhitelisted(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgObserveVote int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgObserveVote, &weightMsgObserveVote, nil,
		func(_ *rand.Rand) {
			weightMsgObserveVote = defaultWeightMsgObserveVote
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgObserveVote,
		omnisimulation.SimulateMsgObserveVote(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateParams int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateParams, &weightMsgUpdateParams, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateParams = defaultWeightMsgUpdateParams
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateParams,
		omnisimulation.SimulateMsgUpdateParams(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgWhitelisted,
			defaultWeightMsgWhitelisted,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				omnisimulation.SimulateMsgWhitelisted(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDewhitelisted,
			defaultWeightMsgDewhitelisted,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				omnisimulation.SimulateMsgDewhitelisted(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgObserveVote,
			defaultWeightMsgObserveVote,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				omnisimulation.SimulateMsgObserveVote(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateParams,
			defaultWeightMsgUpdateParams,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				omnisimulation.SimulateMsgUpdateParams(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
