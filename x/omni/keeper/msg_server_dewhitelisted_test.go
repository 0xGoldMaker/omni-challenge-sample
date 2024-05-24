package keeper_test

import (
	simapp "omni/app"
	"omni/x/omni/keeper"
	"omni/x/omni/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func (suite *KeeperTestSuite) TestDewhitelistGovernanceMsgServerUpdate() {
	// Generate 2 random accounts with 1000000stake balanced
	addrs := simapp.AddTestAddrs(suite.app, suite.ctx, 2, sdk.NewInt(1000000))
	govAddress := sdk.AccAddress(address.Module("gov"))
	govAuthority := govAddress.String()

	for _, tc := range []struct {
		desc        string
		request     *types.MsgDewhitelisted
		whitelisted []string
		err         error
	}{
		{
			desc: "Fail as the authority is not gov module",
			request: &types.MsgDewhitelisted{
				Authority: addrs[0].String(),
				Key:       addrs[1].String(),
			},
			whitelisted: []string{},
			err:         sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", govAuthority, addrs[0].String()),
		},
		{
			desc: "Succeed as the authority is gov module",
			request: &types.MsgDewhitelisted{
				Authority: govAuthority,
				Key:       addrs[0].String(),
			},
			whitelisted: []string{addrs[0].String()},
			err:         nil,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OmniKeeper, suite.ctx

			// Set whitelisted addresses
			for i, ws := range tc.whitelisted {
				suite.app.OmniKeeper.SetWhitelist(ctx, types.Whitelist{Index: (uint64)(i + 1), Address: ws})
			}

			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)
			_, err := srv.Dewhitelisted(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)

				_, found := k.GetWhitelist(ctx,
					(uint64)(1),
				)
				suite.Require().False(found)
			}
		})
	}
}
