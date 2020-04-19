// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package asiutil

import (
	"bytes"
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/database"
	"io"

	"github.com/AsimovNetwork/asimov/protos"
)

// MissVBlockError describes an error due to accessing an element that is out
// of range.
type MissVBlockError string

// Error satisfies the error interface and prints human-readable errors.
func (e MissVBlockError) Error() string {
	return string(e)
}

// VBlock defines an asimov virtual block that provides easier and more
// efficient manipulation of raw virtual blocks.  It also memorizes hashes for
// the virtual block and its virtual transactions on their first access so
// subsequent accesses don't have to repeat the relatively expensive hashing operations.
type VBlock struct {
	msgvBlock        *protos.MsgVBlock // Underlying MsgVBlock
	serializedBytes []byte             // Serialized bytes for the virtual block
	blockHash       common.Hash        // Cached corresponding block hash
	transactions    []*Tx              // Transactions
	txnsGenerated   bool               // ALL wrapped transactions generated
}

// MsgBlock returns the underlying protos.MsgBlock for the Block.
func (b *VBlock) MsgVBlock() *protos.MsgVBlock {
	// Return the cached block.
	return b.msgvBlock
}

// Bytes returns the serialized bytes for the Block.  This is equivalent to
// calling Serialize on the underlying protos.MsgBlock, however it caches the
// result so subsequent calls are more efficient.
func (b *VBlock) Bytes() ([]byte, error) {
	// Return the cached serialized bytes if it has already been generated.
	if len(b.serializedBytes) != 0 {
		return b.serializedBytes, nil
	}

	// Serialize the MsgBlock.
	w := bytes.NewBuffer(make([]byte, 0, b.msgvBlock.SerializeSize()))
	err := b.msgvBlock.Serialize(w)
	if err != nil {
		return nil, err
	}
	serializedBytes := w.Bytes()

	// Cache the serialized bytes and return them.
	b.serializedBytes = serializedBytes
	return serializedBytes, nil
}

// Hash returns the block identifier hash for the Block.  This is equivalent to
// calling BlockHash on the underlying protos.MsgBlock, however it caches the
// result so subsequent calls are more efficient.
func (b *VBlock) Hash() *common.Hash {
	// Return the cached block hash if it has already been generated.
	return &b.blockHash
}

// Tx returns a wrapped transaction (vvsutil.Tx) for the transaction at the
// specified index in the Block.  The supplied index is 0 based.  That is to
// say, the first transaction in the block is txNum 0.  This is nearly
// equivalent to accessing the raw transaction (protos.MsgTx) from the
// underlying protos.MsgBlock, however the wrapped transaction has some helpful
// properties such as caching the hash so subsequent calls are more efficient.
func (b *VBlock) Tx(txNum int) (*Tx, error) {
	// Ensure the requested transaction is in range.
	numTx := len(b.msgvBlock.VTransactions)
	if txNum < 0 || txNum > numTx {
		str := fmt.Sprintf("transaction index %d is out of range [0, %d)",
			txNum, numTx)
		return nil, OutOfRangeError(str)
	}

	// Generate slice to hold all of the wrapped transactions if needed.
	if b.transactions == nil {
		b.transactions = make([]*Tx, numTx)
	}

	// Return the wrapped transaction if it has already been generated.
	if b.transactions[txNum] != nil {
		return b.transactions[txNum], nil
	}

	// Generate and cache the wrapped transaction and return it.
	newTx := NewTx(b.msgvBlock.VTransactions[txNum])
	newTx.SetIndex(txNum)
	b.transactions[txNum] = newTx
	return newTx, nil
}

// Transactions returns a slice of wrapped transactions (vvsutil.Tx) for all
// transactions in the Block.  This is nearly equivalent to accessing the raw
// transactions (protos.MsgTx) in the underlying protos.MsgBlock, however it
// instead provides easy access to wrapped versions (vvsutil.Tx) of them.
func (b *VBlock) Transactions() []*Tx {
	// Return transactions if they have ALL already been generated.  This
	// flag is necessary because the wrapped transactions are lazily
	// generated in a sparse fashion.
	if b.txnsGenerated {
		return b.transactions
	}

	// Generate slice to hold all of the wrapped transactions if needed.
	if len(b.transactions) == 0 {
		b.transactions = make([]*Tx, len(b.msgvBlock.VTransactions))
	}

	// Generate and cache the wrapped transactions for all that haven't
	// already been done.
	for i, tx := range b.transactions {
		if tx == nil {
			newTx := NewTx(b.msgvBlock.VTransactions[i])
			newTx.SetIndex(i)
			b.transactions[i] = newTx
		}
	}

	b.txnsGenerated = true
	return b.transactions
}

// TxHash returns the hash for the requested transaction number in the Block.
// The supplied index is 0 based.  That is to say, the first transaction in the
// block is txNum 0.  This is equivalent to calling TxHash on the underlying
// protos.MsgTx, however it caches the result so subsequent calls are more
// efficient.
func (b *VBlock) TxHash(txNum int) (*common.Hash, error) {
	// Attempt to get a wrapped transaction for the specified index.  It
	// will be created lazily if needed or simply return the cached version
	// if it has already been generated.
	tx, err := b.Tx(txNum)
	if err != nil {
		return nil, err
	}

	// Defer to the wrapped transaction which will return the cached hash if
	// it has already been generated.
	return tx.Hash(), nil
}

// TxLoc returns the offsets and lengths of each transaction in a raw block.
// It is used to allow fast indexing into transactions within the raw byte
// stream.
func (b *VBlock) TxLoc() ([]protos.TxLoc, error) {
	rawMsg, err := b.Bytes()
	if err != nil {
		return nil, err
	}
	rbuf := bytes.NewBuffer(rawMsg)

	var vblock protos.MsgVBlock
	txLocs, err := vblock.DeserializeTxLoc(rbuf)
	if err != nil {
		return nil, err
	}
	return txLocs, err
}

// NewVBlock returns a new instance of an asimov block given an underlying
// protos.MsgVBlock.  See VBlock.
func NewVBlock(msgvBlock *protos.MsgVBlock, blockHash *common.Hash) *VBlock {
	return &VBlock{
		msgvBlock: msgvBlock,
		blockHash: *blockHash,
	}
}

// NewVBlockFromBytes returns a new instance of an asimov block given the
// serialized bytes.  See VBlock.
func NewVBlockFromBytes(serializedBytes []byte, blockHash *common.Hash) (*VBlock, error) {
	br := bytes.NewReader(serializedBytes)
	b, err := newVBlockFromReader(br)
	if err != nil {
		return nil, err
	}
	b.serializedBytes = serializedBytes
	b.blockHash = *blockHash
	return b, nil
}

// NewBlockFromReader returns a new instance of a flow block given a
// Reader to deserialize the block.  See Block.
func newVBlockFromReader(r io.Reader) (*VBlock, error) {
	// Deserialize the bytes into a MsgBlock.
	var msgvBlock protos.MsgVBlock
	err := msgvBlock.Deserialize(r)
	if err != nil {
		return nil, err
	}

	b := VBlock{
		msgvBlock: &msgvBlock,
	}
	return &b, nil
}

// GetBlockPair returns a new instance of an asimov block
func GetBlockPair(db database.Transactor, hash *common.Hash) (block *Block, vblock *VBlock, err error) {
	err = db.View(func(dbTx database.Tx) error {
		blockKey := database.NewNormalBlockKey(hash)
		blockBytes, err := dbTx.FetchBlock(blockKey)
		if err != nil {
			return err
		}
		block, err = NewBlockFromBytes(blockBytes)
		if err != nil {
			return err
		}
		blockKey[common.HashLength] = byte(database.BlockVirtual)

		blockBytes, err = dbTx.FetchBlock(blockKey)
		if err != nil {
			return MissVBlockError(err.Error())
		}
		vblock, err = NewVBlockFromBytes(blockBytes, hash)
		if err != nil {
			return err
		}
		return err
	})
	return
}
