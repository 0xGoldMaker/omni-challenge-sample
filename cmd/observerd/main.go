package main

import (
	tssconfig "omni/observer/config"
	"omni/observer/logging/log"
	"omni/observer/observer"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"omni/observer/omniclient"

	"omni/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	flag "github.com/spf13/pflag"
)

func initSDKConfig() {
	// Set prefixes
	accountPubKeyPrefix := app.AccountAddressPrefix + "pub"
	validatorAddressPrefix := app.AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := app.AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := app.AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := app.AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}

func main() {
	homePath := flag.StringP("home", "h", "$HOME/.omni", "node home path")
	flag.Parse()

	// Init prefix
	initSDKConfig()

	// load configuration files
	tssCfg, err := tssconfig.LoadConfig(*homePath, "config")

	if err != nil {
		log.Fatal().Err(err).Msg("fail to load config ")
		return
	}

	if len(tssCfg.SignerName) == 0 {
		log.Fatal().Msg("signer name is empty")
		return
	}

	if len(tssCfg.SignerPasswd) == 0 {
		log.Fatal().Msg("signer password is empty")
		return
	}

	SignerName := tssCfg.SignerName //args[2]
	SignerPasswd := tssCfg.SignerPasswd

	kb, info, err := omniclient.GetKeyringKeybase("", SignerName, SignerPasswd)
	if err != nil {
		log.Fatal().Msg("fail to get keyring keybase")
		return
	}

	// get Keyring base
	k := omniclient.NewKeysWithKeybase(kb, SignerName, SignerPasswd)

	// Get chain Id
	chainId := strings.ToLower(tssCfg.ChainID)
	//
	cfg := &omniclient.BridgeConfig{
		ChainId:         chainId,
		ChainHost:       tssCfg.ChainHost,
		ChainRPC:        tssCfg.ChainRPC,
		ChainHomeFolder: *homePath,
	}

	accAddr, err := info.GetAddress()
	// Signer pubkey
	pubKey, _ := info.GetPubKey()
	if err != nil {
		log.Fatal().Msg("fail to get keyring keybase")
		return
	}

	addr := accAddr.String()
	// Omni Bridge Service
	OmnichainBridge, err := omniclient.NewOmniBridge(k, cfg, SignerName, pubKey.String(), addr)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create omni tx bridge")
		return
	}

	// Creates Observer instance
	obs, err := observer.NewObserver(OmnichainBridge, tssCfg.NodeRPC)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create observer")
		return
	}

	// Observation start
	obs.Start()

	// Wait....
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("stop signal received")
	// Stop observer
	if err := obs.Stop(); err != nil {
		log.Fatal().Err(err).Msg("fail to stop observer")
		return
	}

	log.Info().Msg("Stopped whole process!!")
}
