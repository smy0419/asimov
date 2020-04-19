// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package netsync

import (
	"github.com/AsimovNetwork/asimov/asiutil"
	"github.com/AsimovNetwork/asimov/blockchain"
	"github.com/AsimovNetwork/asimov/chaincfg"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/crypto"
	"github.com/AsimovNetwork/asimov/mempool"
	"github.com/AsimovNetwork/asimov/mining"
	"github.com/AsimovNetwork/asimov/peer"
	"github.com/AsimovNetwork/asimov/protos"
)

// PeerNotifier exposes methods to notify peers of status changes to
// transactions, blocks, etc. Currently server (in the main package) implements
// this interface.
type PeerNotifier interface {
	AnnounceNewTransactions(newTxs []*mining.TxDesc)

	UpdatePeerHeights(latestBlkHash *common.Hash, latestHeight int32, updateSource *peer.Peer)

	RelayInventory(invVect *protos.InvVect, data interface{})

	TransactionConfirmed(tx *asiutil.Tx)

	AnnounceNewSignature(sig *asiutil.BlockSign)
}

// Config is a configuration struct used to initialize a new SyncManager.
type Config struct {
	PeerNotifier PeerNotifier
	Chain        *blockchain.BlockChain
	TxMemPool    *mempool.TxPool
	SigMemPool   *mempool.SigPool
	ChainParams  *chaincfg.Params

	DisableCheckpoints bool
	MaxPeers           int

	Account *crypto.Account
	BroadcastMessage func(msg protos.Message, exclPeers ...interface{})
}
