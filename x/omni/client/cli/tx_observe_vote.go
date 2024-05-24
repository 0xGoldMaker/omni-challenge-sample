package cli

import (
	"strconv"
	"time"

	"omni/x/omni/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdObserveVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "observe-vote [round] [value]",
		Short: "Broadcast message observe-vote",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRound := args[0]
			argValue := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			round, err := strconv.ParseUint(argRound, 10, 64)
			if err != nil {
				return err
			}
			value, ok := sdk.NewIntFromString(argValue)
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Invalid balance")
			}

			timestamp := time.Now()
			msg := types.NewMsgObserveVote(
				clientCtx.GetFromAddress().String(),
				round,
				value,
				timestamp,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
