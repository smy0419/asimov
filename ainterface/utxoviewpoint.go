// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package ainterface

import (
	"github.com/AsimovNetwork/asimov/blockchain/txo"
	"github.com/AsimovNetwork/asimov/database"
	"github.com/AsimovNetwork/asimov/protos"
	"github.com/AsimovNetwork/asimov/asiutil"
)

type IUtxoViewpoint interface {
	LookupEntry(outpoint protos.OutPoint) *txo.UtxoEntry
	AddTxOut(tx *asiutil.Tx, txOutIdx uint32, blockHeight int32)
	AddTxOuts(tx *asiutil.Tx, blockHeight int32) map[protos.Assets]*txo.LockItem
	RemoveEntry(outpoint protos.OutPoint)
	AddEntry(outpoint protos.OutPoint, entry *txo.UtxoEntry)
	Entries() map[protos.OutPoint]*txo.UtxoEntry
	FetchUtxosMain(db database.Transactor, outpoints map[protos.OutPoint]struct{}) error
	Clone() IUtxoViewpoint
}