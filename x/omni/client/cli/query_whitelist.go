package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"omni/x/omni/types"
)

func CmdListWhitelist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-whitelist",
		Short: "list all whitelist",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllWhitelistRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.WhitelistAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowWhitelist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-whitelist [index]",
		Short: "shows a whitelist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := args[0]

			index, err := strconv.ParseUint(argIndex, 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetWhitelistRequest{
				Index: index,
			}

			res, err := queryClient.Whitelist(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
