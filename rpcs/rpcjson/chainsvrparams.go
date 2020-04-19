// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.


// NOTE: This file is intended to house the RPC commands that are supported by
// a chain server.
package rpcjson

import "math/big"


// AddNodeSubCmd defines the type used in the addnode JSON-RPC command for the
// sub command field.
type AddNodeSubCmd string

const (
	// ANAdd indicates the specified host should be added as a persistent
	// peer.
	ANAdd AddNodeSubCmd = "add"

	// ANRemove indicates the specified peer should be removed.
	ANRemove AddNodeSubCmd = "remove"

	// ANOneTry indicates the specified host should try to connect once,
	// but it should not be made persistent.
	ANOneTry AddNodeSubCmd = "onetry"
)

// TransactionInput represents the inputs to a transaction.  Specifically a
// transaction hash and output number pair.
type TransactionInput struct {
	Txid         string `json:"txid"`
	Vout         uint32 `json:"vout"`
	Assets       string `json:"assets"`
	ScriptPubKey string `json:"scriptpubkey"`
}

type TransactionOutput struct {
	Amount       string `json:"amount"`
	Assets       string `json:"assets"`
	Address      string `json:"address"`
	Data         string `json:"data"`
	ContractType string `json:"contracttype,omitempty"` //create call
}

type TestCallData struct {
	Name      string  `json:"name"`
	Data      string  `json:"data"`
	Caller    string  `json:"caller"`
	Amount    uint32  `json:"amount"`
	Asset     string  `json:"asset"`
	AmountB   big.Int `json:"amountb"`
	VoteValue string  `json:"voteValue"`
}
