package omniclient

import (
	"fmt"

	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	flag "github.com/spf13/pflag"

	"omni/observer/logging/log"

	stypes "github.com/cosmos/cosmos-sdk/types"
)

type (
	// TxID is a string that can uniquely represent a transaction on different
	// block chain
	TxID string
	// TxIDs is a slice of TxID
	TxIDs []TxID
)

// Broadcast Broadcasts tx to Omni
func (b *OmniBridge) Broadcast(msgs ...stypes.Msg) (TxID, error) {
	b.broadcastLock.Lock()
	defer b.broadcastLock.Unlock()

	noTxID := TxID("")
	accountNumber, seqNumber, err := b.GetAccountNumberAndSequenceNumber()
	if err != nil {
		b.logger.Error().Err(err).Msg("fail to get account number and sequence number from Omni")
		return noTxID, fmt.Errorf("fail to get account number and sequence number from Omni : %w", err)
	}

	flags := flag.NewFlagSet("Omni", 0)

	ctx := b.GetContext()
	factory, err := clienttx.NewFactoryCLI(ctx, flags)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to initiate tx factory")
		return noTxID, err
	}

	factory = factory.WithAccountNumber(accountNumber)
	factory = factory.WithSequence(seqNumber)
	factory = factory.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
	factory = factory.WithGas(uint64(2000000))

	builder, err := factory.BuildUnsignedTx(msgs...)
	if err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return noTxID, err
	}

	builder.SetGasLimit(4000000)
	err = clienttx.Sign(factory, ctx.GetFromName(), builder, true)
	if err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return noTxID, err
	}

	txBytes, err := ctx.TxConfig.TxEncoder()(builder.GetTx())
	if err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return noTxID, err
	}

	// broadcast to a Tendermint node
	commit, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		b.logger.Error().Err(err).Msgf("fail to broadcast tx: %s", err.Error())
		return noTxID, fmt.Errorf("fail to broadcast tx: %w", err)
	}

	txHash, err := NewTxID(commit.TxHash)
	if err != nil {
		b.logger.Error().Err(err).Msgf("fail to convert txhash: %s", err.Error())
		return BlankTxID, fmt.Errorf("fail to convert txhash: %w", err)
	}

	// Code will be the tendermint ABICode , it start at 1 , so if it is an error , code will not be zero
	if commit.Code > 0 {
		b.logger.Error().Err(err).Msgf("fail to broadcast to OmniChain,code:%d, log:%s", commit.Code, commit.RawLog)
		return txHash, fmt.Errorf("fail to broadcast to OmniChain,code:%d, log:%s", commit.Code, commit.RawLog)
	}

	b.logger.Info().Msgf("Broadcasted a tx: %s", txHash)

	return txHash, nil
}
