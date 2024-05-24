package keeper_test

import (
	"strconv"
	"time"

	simapp "omni/app"
	"omni/x/omni/keeper"
	"omni/x/omni/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func (suite *KeeperTestSuite) TestObserveVoteMsgServerUpdate() {
	// Generate 2 random accounts with 1000000stake balanced
	addrs := simapp.AddTestAddrs(suite.app, suite.ctx, 2, sdk.NewInt(1000000))

	for _, tc := range []struct {
		desc                 string
		request              *types.MsgObserveVote
		isWhitelistEnabled   bool
		whitelistedAddresses []string
		err                  error
	}{
		{
			desc: "Succeed without whitelisting",
			request: &types.MsgObserveVote{
				Creator:   addrs[0].String(),
				Round:     0,
				Value:     sdk.NewInt(1),
				Timestamp: time.Now(),
			},
			isWhitelistEnabled:   false,
			whitelistedAddresses: []string{},
			err:                  nil,
		},
		{
			desc: "Succeed with whitelisting",
			request: &types.MsgObserveVote{
				Creator:   addrs[0].String(),
				Round:     0,
				Value:     sdk.NewInt(1),
				Timestamp: time.Now(),
			},
			isWhitelistEnabled:   true,
			whitelistedAddresses: []string{addrs[0].String(), addrs[1].String()},
			err:                  nil,
		},
		{
			desc: "Fail without whitelisting because of currnet round",
			request: &types.MsgObserveVote{
				Creator:   addrs[0].String(),
				Round:     1,
				Value:     sdk.NewInt(1),
				Timestamp: time.Now(),
			},
			isWhitelistEnabled:   false,
			whitelistedAddresses: []string{},
			err:                  sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid round"),
		},
		{
			desc: "Fail with whitelisting but not having any whitelisted",
			request: &types.MsgObserveVote{
				Creator:   addrs[0].String(),
				Round:     1,
				Value:     sdk.NewInt(1),
				Timestamp: time.Now(),
			},
			isWhitelistEnabled:   true,
			whitelistedAddresses: []string{},
			err:                  sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid request"),
		},
		{
			desc: "Fail with whitelisting but tried with noted whitelisted",
			request: &types.MsgObserveVote{
				Creator:   addrs[0].String(),
				Round:     1,
				Value:     sdk.NewInt(1),
				Timestamp: time.Now(),
			},
			isWhitelistEnabled:   true,
			whitelistedAddresses: []string{addrs[1].String()},
			err:                  sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid request"),
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OmniKeeper, suite.ctx

			// Set params
			params := types.DefaultParams()
			params.IsWhitelistEnabled = tc.isWhitelistEnabled
			suite.app.OmniKeeper.SetParams(ctx, &params)

			// Set whitelisted addresses
			for i, ws := range tc.whitelistedAddresses {
				suite.app.OmniKeeper.SetWhitelist(ctx, types.Whitelist{Index: (uint64)(i), Address: ws})
			}

			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)
			_, err := srv.ObserveVote(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				rst, found := k.GetObserveVote(ctx,
					(uint64)(1),
				)
				suite.Require().True(found)
				suite.Require().Equal(tc.request.Creator, rst.Voter)
			}
		})
	}
}
