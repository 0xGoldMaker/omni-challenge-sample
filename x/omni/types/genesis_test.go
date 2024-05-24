package types_test

import (
	"testing"
	time "time"

	"omni/x/omni/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				BalanceList: []types.Balance{
					{
						Index:       (uint64)(0),
						Balance:     sdk.NewInt(1),
						LastUpdated: time.Now(),
					},
					{
						Index:       (uint64)(1),
						Balance:     sdk.NewInt(1),
						LastUpdated: time.Now(),
					},
				},
				WhitelistList: []types.Whitelist{
					{
						Index:   (uint64)(0),
						Address: "",
					},
					{
						Index:   (uint64)(1),
						Address: "",
					},
				},
				ObserveVoteList: []types.ObserveVote{
					{
						Index:     (uint64)(0),
						Voter:     "",
						Round:     (uint64)(0),
						Value:     sdk.NewInt(1),
						Timestamp: time.Now(),
					},
					{
						Index:     (uint64)(1),
						Voter:     "",
						Round:     (uint64)(0),
						Value:     sdk.NewInt(1),
						Timestamp: time.Now(),
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated balance",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				BalanceList: []types.Balance{
					{
						Index:       (uint64)(0),
						Balance:     sdk.NewInt(1),
						LastUpdated: time.Now(),
					},
					{
						Index:       (uint64)(0),
						Balance:     sdk.NewInt(2),
						LastUpdated: time.Now(),
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated whitelist",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				WhitelistList: []types.Whitelist{
					{
						Index:   (uint64)(0),
						Address: "",
					},
					{
						Index:   (uint64)(0),
						Address: "",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated observeVote",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ObserveVoteList: []types.ObserveVote{
					{
						Index:     (uint64)(0),
						Voter:     "",
						Round:     (uint64)(0),
						Value:     sdk.NewInt(1),
						Timestamp: time.Now(),
					},
					{
						Index:     (uint64)(0),
						Voter:     "",
						Round:     (uint64)(0),
						Value:     sdk.NewInt(1),
						Timestamp: time.Now(),
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
