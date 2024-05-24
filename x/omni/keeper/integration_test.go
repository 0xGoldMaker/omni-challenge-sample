package keeper_test

import (
	simapp "omni/app"
	"time"

	"omni/x/omni/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (suite *KeeperTestSuite) TestIntegrationTesting() {
	// Generate 3 random accounts with 1000000stake balanced
	addrs := simapp.AddTestAddrs(suite.app, suite.ctx, 3, sdk.NewInt(1000000))
	for _, tc := range []struct {
		desc               string
		round              uint64
		observeVotes       []types.ObserveVote
		balance            sdk.Int
		whitelisted        []string
		isWhitelistEnabled bool
		err                error
	}{
		{
			desc:  "Integration testing - Succeed without whitelisting",
			round: 0,
			observeVotes: []types.ObserveVote{
				{
					Index:     uint64(1),
					Voter:     addrs[0].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
				{
					Index:     uint64(2),
					Voter:     addrs[1].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
				{
					Index:     uint64(3),
					Voter:     addrs[2].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
			},
			balance:            sdk.NewInt(5),
			whitelisted:        []string{addrs[0].String(), addrs[1].String(), addrs[2].String()},
			isWhitelistEnabled: true,
			err:                nil,
		},
		{
			desc:  "Integration testing - Succeed without whitelisting",
			round: 0,
			observeVotes: []types.ObserveVote{
				{
					Index:     uint64(1),
					Voter:     addrs[0].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
				{
					Index:     uint64(2),
					Voter:     addrs[1].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
				{
					Index:     uint64(3),
					Voter:     addrs[2].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
			},
			balance:            sdk.NewInt(5),
			whitelisted:        []string{},
			isWhitelistEnabled: false,
			err:                nil,
		},
		{
			desc:  "Integration testing - Failed because of invalid round",
			round: 1,
			observeVotes: []types.ObserveVote{
				{
					Index:     uint64(1),
					Voter:     addrs[0].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
				{
					Index:     uint64(2),
					Voter:     addrs[1].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
				{
					Index:     uint64(3),
					Voter:     addrs[2].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
			},
			balance:            sdk.NewInt(5),
			whitelisted:        []string{addrs[0].String(), addrs[1].String(), addrs[2].String()},
			isWhitelistEnabled: true,
			err:                sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "in sufficient observation voted"),
		},
		{
			desc:  "Integration testing - Failed because of not whitelisted",
			round: 1,
			observeVotes: []types.ObserveVote{
				{
					Index:     uint64(1),
					Voter:     addrs[0].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
				{
					Index:     uint64(2),
					Voter:     addrs[1].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
				{
					Index:     uint64(3),
					Voter:     addrs[2].String(),
					Round:     0,
					Value:     sdk.NewInt(5),
					Timestamp: time.Now(),
				},
			},
			balance:            sdk.NewInt(5),
			whitelisted:        []string{addrs[1].String(), addrs[2].String()},
			isWhitelistEnabled: true,
			err:                sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "in sufficient observation voted"),
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OmniKeeper, suite.ctx

			// Set params
			params := types.DefaultParams()
			params.IsWhitelistEnabled = tc.isWhitelistEnabled
			k.SetParams(ctx, &params)

			// Set whitelisted addresses
			for i, ws := range tc.whitelisted {
				k.SetWhitelist(ctx, types.Whitelist{Index: (uint64)(i), Address: ws})
			}

			// Set observation votes
			for _, ov := range tc.observeVotes {
				k.SetObserveVote(ctx, ov)
			}

			// Update balance using the voted observations
			err := k.UpdateBalanceData(ctx, tc.round, params.MinConsensus)

			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				rst, found := k.GetBalance(ctx,
					(uint64)(0),
				)
				suite.Require().True(found)
				suite.Require().Equal(tc.balance, rst.Balance)
			}
		})
	}
}
