package observer

import (
	"context"
	"fmt"
	"math/big"
	"omni/observer/logging"
	"omni/observer/logging/log"
	"omni/observer/omniclient"
	"sync"
	"time"

	"omni/x/omni/types"

	"github.com/cenkalti/backoff"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Observer observer service
type Observer struct {
	logger          logging.Logger
	lock            *sync.Mutex
	stopChan        chan struct{}
	OmniChainBridge *omniclient.OmniBridge

	// Tx Msg List
	msgList []sdktypes.Msg

	// Msgs Array for each Tx
	msgChan chan sdktypes.Msg

	// RPC endpoint of the node
	nodeRPC string
}

// NewObserver create a new instance of Observer for chain
func NewObserver(chainBridge *omniclient.OmniBridge, nodeRPC string) (*Observer, error) {
	return &Observer{
		logger:          log.With().Str("module", "Observer").Logger(),
		lock:            &sync.Mutex{},
		stopChan:        make(chan struct{}),
		OmniChainBridge: chainBridge,
		msgList:         make([]sdktypes.Msg, 0),
		msgChan:         make(chan sdktypes.Msg),
		nodeRPC:         nodeRPC,
	}, nil
}

// Start Entire Process
func (o *Observer) Start() {
	// Starts TSS Process
	go o.StartObserveProcess()

	// Process Broadcasting Msgs to Omni
	go o.ProcessIncomingTx()

	// Starts Omnichain MSG broadcasting process
	go o.StartMSGBroadcastingProcess()
}

func (o *Observer) Stop() error {
	close(o.stopChan)

	return nil
}

// Start Observe Process
func (o *Observer) StartObserveProcess() {
	for {
		select {
		case <-o.stopChan:
			return
		case <-time.After(10 * time.Second):
			o.ObserveStorageData()
		}
	}
}

// Observe storage balance
func (o *Observer) ObserveStorageData() {
	// Fetch the smart contract address from Omni chain
	contractAddress, currentRound, err := o.OmniChainBridge.GetParam()
	if err != nil {
		log.Fatal().Err(err).Msg("fail to get param")
		return
	}

	// Get voter
	_, voter := o.OmniChainBridge.GetVoterInfo()

	// Initiates the client
	client, err := ethclient.Dial(o.nodeRPC)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to connect ETH RPC node")
		return
	}

	defer client.Close()

	// Context
	ctx := context.Background()

	// Contract hex address in address format
	contractAddressHex := common.HexToAddress(contractAddress)

	// Pick up the first store
	key := common.HexToHash("0x0")

	// Fetch balance
	data, err := client.StorageAt(ctx, contractAddressHex, key, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to get storage value")
		return
	}

	// Balance in sdk.Int type
	balance := sdktypes.NewIntFromBigInt(big.NewInt(0).SetBytes(data))

	// Broadcast observe vote to omni chain
	o.msgChan <- types.NewMsgObserveVote(voter, currentRound, balance, time.Now())

	log.Info().Msg("Balance observed")
}

// Store Incoming Tx
func (o *Observer) ProcessIncomingTx() {
	for {
		select {
		case <-o.stopChan:
			return
		case msg := <-o.msgChan:
			o.lock.Lock()
			o.msgList = append(o.msgList, msg)
			o.lock.Unlock()
		}
	}
}

// Broadcast Tx to Omnichain
func (o *Observer) StartMSGBroadcastingProcess() {
	bf := backoff.NewExponentialBackOff()
	bf.MaxElapsedTime = time.Second * 25

	for {
		select {
		case <-o.stopChan:
			return
		default:
			o.lock.Lock()
			if len(o.msgList) < 1 {
				o.lock.Unlock()

				// Sleep
				time.Sleep(time.Second * 5)
				break
			}

			// Copy msg array
			msgs := make([]sdktypes.Msg, 0)
			for _, msg := range o.msgList {
				m := msg
				msgs = append(msgs, m)
			}
			o.msgList = o.msgList[:0]
			o.lock.Unlock()

			// Broadcasting
			err := backoff.Retry(func() error {
				_, err := o.OmniChainBridge.Broadcast(msgs...)
				if err != nil {
					return fmt.Errorf("fail to send the tx to Omnichain: %w", err)
				}
				return nil
			}, bf)

			// If it wasn't successful to broadcast
			if err != nil {
				o.lock.Lock()
				o.msgList = append(o.msgList, msgs...)
				o.lock.Unlock()
			}

			// Sleep
			time.Sleep(5 * time.Second)
		}
	}
}
