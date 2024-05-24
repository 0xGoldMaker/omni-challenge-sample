package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"omni/x/omni/types"
)

func CmdListObserveVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-observe-vote",
		Short: "list all observe-vote",
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

			params := &types.QueryAllObserveVoteRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ObserveVoteAll(cmd.Context(), params)
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

func CmdShowObserveVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-observe-vote [index]",
		Short: "shows a observe-vote",
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

			params := &types.QueryGetObserveVoteRequest{
				Index: index,
			}

			res, err := queryClient.ObserveVote(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
