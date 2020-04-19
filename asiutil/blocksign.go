// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package asiutil

import (
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/protos"
)

// BlockSign stores the signature hashes of msgBlock
type BlockSign struct {
	MsgSign  *protos.MsgBlockSign
	signHash *common.Hash
}

// Hash returns hash in the blockSign.
// if not found, generate new one
func (t *BlockSign) Hash() *common.Hash {
	// Return the cached hash if it has already been generated.
	if t.signHash != nil {
		return t.signHash
	}

	// Cache the hash and return it.
	hash := t.MsgSign.SignHash()
	t.signHash = &hash
	return &hash
}

// NewBlockSign sets MsgBlockSign into BlockSign
func NewBlockSign(msgTx *protos.MsgBlockSign) *BlockSign {
	return &BlockSign{
		MsgSign: msgTx,
	}
}
