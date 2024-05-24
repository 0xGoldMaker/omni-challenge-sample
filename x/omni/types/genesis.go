package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		BalanceList:     []Balance{},
		WhitelistList:   []Whitelist{},
		ObserveVoteList: []ObserveVote{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in balance
	balanceIndexMap := make(map[string]struct{})

	for _, elem := range gs.BalanceList {
		index := string(BalanceKey(elem.Index))
		if _, ok := balanceIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for balance")
		}
		balanceIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in whitelist
	whitelistIndexMap := make(map[string]struct{})

	for _, elem := range gs.WhitelistList {
		index := string(WhitelistKey(elem.Index))
		if _, ok := whitelistIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for whitelist")
		}
		whitelistIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in observeVote
	observeVoteIndexMap := make(map[string]struct{})

	for _, elem := range gs.ObserveVoteList {
		index := string(ObserveVoteKey(elem.Index))
		if _, ok := observeVoteIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for observeVote")
		}
		observeVoteIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
