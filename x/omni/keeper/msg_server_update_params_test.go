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

func (suite *KeeperTestSuite) TestUpdateParamsGovernanceMsgServerUpdate() {
	// Generate 2 random accounts with 1000000stake balanced
	addrs := simapp.AddTestAddrs(suite.app, suite.ctx, 2, sdk.NewInt(1000000))
	govAddress := sdk.AccAddress(address.Module("gov"))
	govAuthority := govAddress.String()

	for _, tc := range []struct {
		desc        string
		request     *types.MsgUpdateParams
		whitelisted []string
		err         error
	}{
		{
			desc: "Succees with right params",
			request: &types.MsgUpdateParams{
				Authority: govAuthority,
				Params: &types.Params{
					NumEpochs:          15,
					MinConsensus:       5,
					CurRound:           1,
					IsWhitelistEnabled: true,
				},
			},
			err: nil,
		},
		{
			desc: "Fail as the authority is not gov module",
			request: &types.MsgUpdateParams{
				Authority: addrs[0].String(),
				Params: &types.Params{
					NumEpochs:          15,
					MinConsensus:       5,
					CurRound:           1,
					IsWhitelistEnabled: true,
				},
			},
			err: sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", govAuthority, addrs[0].String()),
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OmniKeeper, suite.ctx

			// Set params
			params := types.DefaultParams()
			k.SetParams(ctx, &params)

			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)
			_, err := srv.UpdateParams(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)

				params = k.GetParams(ctx)

				suite.Require().Equal(tc.request.Params.NumEpochs, params.NumEpochs)
				suite.Require().Equal(tc.request.Params.MinConsensus, params.MinConsensus)
				suite.Require().Equal(tc.request.Params.IsWhitelistEnabled, params.IsWhitelistEnabled)
			}
		})
	}
}
