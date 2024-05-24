package cli

import (
	"errors"
	"strconv"

	"omni/x/omni/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [num epoch] [min consensus] [is whitelist enabled] [contract address]",
		Short: "Broadcast message update-params",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argNumEpoch := args[0]
			argMinConsensus := args[1]
			argIsWhitelistEnabled := args[2]
			argContractAddress := args[3]

			numEpoch, err := strconv.ParseUint(argNumEpoch, 10, 64)
			if err != nil {
				return err
			}

			minConsensus, err := strconv.ParseUint(argMinConsensus, 10, 64)
			if err != nil {
				return err
			}

			is_whitelist_enabled, err := strconv.ParseBool(argIsWhitelistEnabled)
			if err != nil {
				return err
			}

			// Check is eth address
			if !eth.IsHexAddress(argContractAddress) {
				return errors.New("invalid address")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			summary, err := cmd.Flags().GetString(cli.FlagSummary)
			if err != nil {
				return err
			}

			metadata, err := cmd.Flags().GetString(cli.FlagMetadata)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			govAddress := sdk.AccAddress(address.Module(govtypes.ModuleName))
			msg := types.NewMsgUpdateParams(
				govAddress.String(),
				types.Params{
					// number of block epoch in refreshing balance
					NumEpochs: numEpoch,
					// minimum number of validators voted
					MinConsensus: minConsensus,
					// currnet consensus round
					CurRound: 0,
					// if whitelisting enabled
					IsWhitelistEnabled: is_whitelist_enabled,
					// smart contract address
					ContractAddress: argContractAddress,
				},
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			govMsg, err := v1.NewMsgSubmitProposal([]sdk.Msg{msg}, deposit, signer.String(), metadata, title, summary)
			if err != nil {
				return err
			}

			err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), govMsg)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagSummary, "", "summary of proposal")
	cmd.Flags().String(cli.FlagMetadata, "", "metadata of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")

	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagSummary)
	_ = cmd.MarkFlagRequired(cli.FlagMetadata)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
