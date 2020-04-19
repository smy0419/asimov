// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package minersync

import (
	"errors"
	"github.com/AsimovNetwork/asimov/chaincfg"
	"github.com/AsimovNetwork/asimov/rpcs/rpc"
	"math/rand"
	"net/http"
)

// BtcValidatorTx models the validatorTxs field of the listblockminerinfo command.
type BtcValidatorTx struct {
	Txid        string   `json:"txid"`
	Vin         []string `json:"vin"`
	OutAddress  string   `json:"outAddress"`
	AddressType int32    `json:"addressType"`
}

// GetBitcoinBlockMinerInfoResult models the data returned from the listblockminerinfo
// command.
type GetBitcoinBlockMinerInfoResult struct {
	Height       int32            `json:"height"`
	Pool         string           `json:"pool"`
	Hash         string           `json:"hash"`
	Time         uint32           `json:"time"`
	Address      string           `json:"address"`
	ValidatorTxs []BtcValidatorTx `json:"ValidatorTxs"`
}

// GetBlockChainInfoResult models the data returned from the getblockchaininfo
// command.
type GetBlockChainInfoResult struct {
	Chain         string `json:"chain"`
	Blocks        int32  `json:"blocks"`
	BestBlockHash string `json:"bestblockhash"`
}

// BtcClient implements the IBtcClient interface and represents a connection
// to an RPC server.
type BtcClient struct {
	client        *rpc.Client
	bitcoinParams []*chaincfg.BitcoinParams
	serverIndex   int
}

// NewBtcClient returns a new BtcClient using the active net params.
func NewBtcClient() (*BtcClient, error) {
	return NewBtcClientWithParams(chaincfg.ActiveNetParams.Params)
}

// NewBtcClientWithParams returns a new BtcClient using the net params.
func NewBtcClientWithParams(params *chaincfg.Params) (*BtcClient, error) {
	btcClient := new(BtcClient)
	//Configuration takes precedence over default parameters
	btcClient.bitcoinParams = chaincfg.Cfg.BtcParams
	//Random order of default parameters
	defaultParams := make([]*chaincfg.BitcoinParams, 0, len(params.Bitcoin))
	for _, param := range params.Bitcoin {
		defaultParams = append(defaultParams, param)
	}
	rand.Shuffle(len(defaultParams), func(i, j int) {
		defaultParams[i], defaultParams[j] = defaultParams[j], defaultParams[i]
	})

	for _, param := range defaultParams {
		btcClient.bitcoinParams = append(btcClient.bitcoinParams, param)
	}
	btcClient.serverIndex = len(btcClient.bitcoinParams) - 1
	err := btcClient.SwitchClient()
	if err != nil {
		return nil, err
	}
	return btcClient, nil
}

// SwitchClient switch target bitcoin server.
func (c *BtcClient) SwitchClient() error {
	c.serverIndex = (c.serverIndex + 1) % len(c.bitcoinParams)
	param := c.bitcoinParams[c.serverIndex]
	endpoint := "http://" + param.Host
	client, err := rpc.DialHTTPWithClientAuthorization(endpoint, new(http.Client),
		param.RpcUser, param.RpcPassword)
	if err != nil {
		return err
	}
	c.client = client
	return nil
}

// Call call an rpc request to bitcoin service, and switch target bitcoin
// server when it failed.
func (c *BtcClient) Call(result interface{}, method string, args ...interface{}) error {
	retryTimes := 0
	for retryTimes < len(c.bitcoinParams) {
		err := c.client.Call(result, method, args...)
		if err == nil {
			return nil
		}
		retryTimes++
		log.Warnf("Failed to call bitcoin server : %s", err.Error())

		for retryTimes < len(c.bitcoinParams) {
			err := c.SwitchClient()
			if err == nil {
				break
			}
			retryTimes++
			log.Warnf("SwitchClient failed : %s", err.Error())
		}
	}
	return errors.New("Failed to call all bitcoin servers")
}

// GetBitcoinMinerInfo call listblockminerinfo to bitcoin service.
func (c *BtcClient) GetBitcoinMinerInfo(result interface{}, height, count int32) error {
	err := c.Call(result, "listblockminerinfo", height, count)
	return err
}

// GetBitcoinBlockChainInfo call getblockchaininfo to bitcoin service.
func (c *BtcClient) GetBitcoinBlockChainInfo(result interface{}) error {
	err := c.Call(result, "getblockchaininfo")
	return err
}
