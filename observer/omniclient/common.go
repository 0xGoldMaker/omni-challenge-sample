package omniclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"omni/app"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/hashicorp/go-retryablehttp"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"omni/observer/logging"
	"omni/observer/logging/log"

	"omni/x/omni/types"
)

// BlankTxID represent blank
var BlankTxID = TxID("0000000000000000000000000000000000000000000000000000000000000000")

// NewTxID parse the input hash as TxID
func NewTxID(hash string) (TxID, error) {
	switch len(hash) {
	case 64:
		// do nothing
	case 66: // ETH check
		if !strings.HasPrefix(hash, "0x") {
			err := fmt.Errorf("txid error: must be 66 characters (got %d)", len(hash))
			return TxID(""), err
		}
	default:
		err := fmt.Errorf("txid error: must be 64 characters (got %d)", len(hash))
		return TxID(""), err
	}

	return TxID(strings.ToUpper(hash)), nil
}

type BridgeConfig struct {
	ChainId         string
	ChainHost       string
	ChainRPC        string
	ChainHomeFolder string
}

// OmniBridge will be used to send tx to omnichain
type OmniBridge struct {
	keys          *Keys
	cfg           BridgeConfig
	blockHeight   uint64
	httpClient    *retryablehttp.Client
	broadcastLock *sync.RWMutex
	codec         codec.Codec

	signerName           string
	lastBlockHeightCheck time.Time
	pubKey               string
	voterAddress         string

	logger logging.Logger
}

// AccountResp the response from thorclient
type AccountResp struct {
	Account struct {
		AccountNumber uint64 `json:"account_number,string"`
		Sequence      uint64 `json:"sequence,string"`
	} `json:"account"`
}

// NewOmniBridge create a new instance of OmniBridge
func NewOmniBridge(k *Keys, cfg *BridgeConfig, signer string, pubKey string, voter string) (*OmniBridge, error) {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	return &OmniBridge{
		keys:          k,
		httpClient:    httpClient,
		codec:         app.MakeEncodingConfig().Marshaler,
		signerName:    signer,
		broadcastLock: &sync.RWMutex{},
		pubKey:        pubKey,
		voterAddress:  voter,
		cfg:           *cfg,

		logger: log.With().Str("module", "OmnichainBridge").Logger(),
	}, nil
}

// GetContext return a valid context with all relevant values set
func (b *OmniBridge) GetContext() client.Context {
	ctx := client.Context{}
	ctx = ctx.WithKeyring(b.keys.GetKeybase())
	ctx = ctx.WithChainID(b.cfg.ChainId)
	ctx = ctx.WithHomeDir(b.cfg.ChainHomeFolder)
	ctx = ctx.WithFromName(b.signerName)
	addr, err := b.keys.GetSignerInfo().GetAddress()
	if err != nil {
		panic(err)
	}
	ctx = ctx.WithFromAddress(addr)
	ctx = ctx.WithBroadcastMode("sync")

	encodingConfig := app.MakeEncodingConfig()
	ctx = ctx.WithCodec(encodingConfig.Marshaler)
	ctx = ctx.WithInterfaceRegistry(encodingConfig.InterfaceRegistry)
	ctx = ctx.WithTxConfig(encodingConfig.TxConfig)
	ctx = ctx.WithLegacyAmino(encodingConfig.Amino)
	ctx = ctx.WithAccountRetriever(authtypes.AccountRetriever{})

	remote := b.cfg.ChainRPC
	if !strings.HasSuffix(b.cfg.ChainHost, "http") {
		remote = fmt.Sprintf("tcp://%s", remote)
	}
	ctx = ctx.WithNodeURI(remote)
	client, err := rpchttp.New(remote, "/websocket")
	if err != nil {
		b.logger.Panic().Err(err).Msg(err.Error())
		panic(err)
	}

	ctx = ctx.WithClient(client)
	return ctx
}

// Get Path
func (b *OmniBridge) getWithPath(path string) ([]byte, int, error) {
	return b.get(b.getOmnichainURL(path))
}

// GetomnichainURL with the given path
func (b *OmniBridge) getOmnichainURL(path string) string {
	uri := url.URL{
		Scheme: "http",
		Host:   b.cfg.ChainHost,
		Path:   path,
	}
	return uri.String()
}

// Get handle all the low level http GET calls using retryablehttp.OmniBridge
func (b *OmniBridge) get(url string) ([]byte, int, error) {
	resp, err := b.httpClient.Get(url)
	if err != nil {
		b.logger.Fatal().Err(err).Msg("ffailed to GET from omnichain")
		return nil, http.StatusNotFound, fmt.Errorf("failed to GET from omnichain: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.logger.Fatal().Err(err).Msg("failed to close response body")
		}
	}()

	buf, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return buf, resp.StatusCode, errors.New("Status code: " + resp.Status + " returned")
	}
	if err != nil {
		b.logger.Fatal().Err(err).Msg("fail_read_omnichain_resp")
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return buf, resp.StatusCode, nil
}

// GetAccountNumberAndSequenceNumber returns account and Sequence number required to post into omnichain
func (b *OmniBridge) GetAccountNumberAndSequenceNumber() (uint64, uint64, error) {
	addr, err := b.keys.GetSignerInfo().GetAddress()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get signer address: %w", err)
	}

	path := fmt.Sprintf("%s/%s", "/cosmos/auth/v1beta1/accounts", addr)

	body, _, err := b.getWithPath(path)
	if err != nil {
		b.logger.Error().Err(err).Msgf("failed to get auth accounts: %s", err.Error())
		return 0, 0, fmt.Errorf("failed to get auth accounts: %w", err)
	}

	var resp AccountResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal account resp: %w", err)
	}

	acc := resp.Account
	return acc.AccountNumber, acc.Sequence, nil
}

// GetParam returns contract address and currrent round
func (b *OmniBridge) GetParam() (string, uint64, error) {
	path := "/omni/omni/params"

	buf, _, err := b.getWithPath(path)
	if err != nil {
		b.logger.Error().Err(err).Msgf("failed to get token addresses: %s", err.Error())
		return "", 0, fmt.Errorf("failed to get token addresses: %w", err)
	}

	var resp types.QueryParamsResponse
	if err := b.codec.UnmarshalJSON(buf, &resp); err != nil {
		return "", 0, fmt.Errorf("fail to unmarshal token addresses: %w", err)
	}

	address := resp.Params.ContractAddress
	round := resp.Params.CurRound

	return address, round, nil
}

// GetVoterInfo returns public key and voter address
func (b *OmniBridge) GetVoterInfo() (string, string) {
	return b.pubKey, b.voterAddress
}

// MakeLegacyCodec creates codec
func MakeLegacyCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	banktypes.RegisterLegacyAminoCodec(cdc)
	authtypes.RegisterLegacyAminoCodec(cdc)
	sdk.RegisterLegacyAminoCodec(cdc)
	// stypes.RegisterCodec(cdc)
	return cdc
}
