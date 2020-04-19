// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"bytes"
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/serialization"
	"io"
)

// MsgVBlock implements the Message interface and represents an asimov
// block message. It is used to deliver block and transaction information in
// response to a getdata message (MsgGetData) for a given block hash.
type MsgVBlock struct {
	Version       uint32
	MerkleRoot    common.Hash
	VTransactions []*MsgTx
}

// AddTransaction adds a transaction to the message.
func (msg *MsgVBlock) AddTransaction(tx *MsgTx) {
	msg.VTransactions = append(msg.VTransactions, tx)
}

// Deserialize decodes a virtual block from r into the receiver using a format
// that is suitable for long-term storage such as a database while respecting
// the Version field in the virtual block. The protos encoding can technically
// differ depending on the protocol version and doesn't even really need to
// match the format of a stored virtual block at all.  As of the time this
// comment was written, the encoded virtual block is the same in both instances,
// but there is a distinct difference and separating the two allows the API to
// be flexible enough to deal with changes.
func (msg *MsgVBlock) Deserialize(r io.Reader) error {
	err := serialization.ReadUint32(r, &msg.Version)
	if err != nil {
		return err
	}

	txCount, err := serialization.ReadVarInt(r, 0)
	if err != nil {
		return err
	}

	// Prevent more transactions than could possibly fit into a block.
	// It would be possible to cause memory exhaustion and panics without
	// a sane upper bound on this count.
	if txCount > maxTxPerBlock {
		str := fmt.Sprintf("too many transactions to fit into a vblock "+
			"[count %d, max %d]", txCount, maxTxPerBlock)
		return messageError("MsgVBlock.Deserialize", str)
	}

	msg.VTransactions = make([]*MsgTx, 0, txCount)
	for i := uint64(0); i < txCount; i++ {
		tx := MsgTx{}
		err := tx.Deserialize(r)
		if err != nil {
			return err
		}
		msg.VTransactions = append(msg.VTransactions, &tx)
	}

	err = serialization.ReadNBytes(r, msg.MerkleRoot[:], common.HashLength)
	if err != nil {
		return err
	}

	return nil
}

// DeserializeTxLoc decodes r in the same manner Deserialize does, but it takes
// a byte buffer instead of a generic reader and returns a slice containing the
// start and length of each transaction within the raw data that is being
// deserialized.
func (msg *MsgVBlock) DeserializeTxLoc(r *bytes.Buffer) ([]TxLoc, error) {
	fullLen := r.Len()
	err := serialization.ReadUint32(r, &msg.Version)
	if err != nil {
		return nil, err
	}

	txCount, err := serialization.ReadVarInt(r, 0)
	if err != nil {
		return nil, err
	}

	// Prevent more transactions than could possibly fit into a block.
	// It would be possible to cause memory exhaustion and panics without
	// a sane upper bound on this count.
	if txCount > maxTxPerBlock {
		str := fmt.Sprintf("too many transactions to fit into a vblock "+
			"[count %d, max %d]", txCount, maxTxPerBlock)
		return nil, messageError("MsgVBlock.DeserializeTxLoc", str)
	}

	// Deserialize each transaction while keeping track of its location
	// within the byte stream.
	msg.VTransactions = make([]*MsgTx, 0, txCount)
	txLocs := make([]TxLoc, txCount)
	for i := uint64(0); i < txCount; i++ {
		txLocs[i].TxStart = fullLen - r.Len()
		tx := MsgTx{}
		err := tx.Deserialize(r)
		if err != nil {
			return nil, err
		}
		msg.VTransactions = append(msg.VTransactions, &tx)
		txLocs[i].TxLen = (fullLen - r.Len()) - txLocs[i].TxStart
	}

	return txLocs, nil
}

// Serialize encodes the virtual block to w using a format that suitable for
// long-term storage such as a database while respecting the Version field in
// the block. The protos encoding can technically differ depending on the
// protocol version and doesn't even really need to match the format of a
// stored virtual block at all.  As of the time this comment was written, the
// encoded virtual block is the same in both instances, but there is a
// distinct difference and separating the two allows the API to be flexible
// enough to deal with changes.
func (msg *MsgVBlock) Serialize(w io.Writer) error {
	err := serialization.WriteUint32(w, msg.Version)
	if err != nil {
		return err
	}

	err = serialization.WriteVarInt(w, 0, uint64(len(msg.VTransactions)))
	if err != nil {
		return err
	}

	for _, tx := range msg.VTransactions {
		err = tx.Serialize(w)
		if err != nil {
			return err
		}
	}

	err = serialization.WriteNBytes(w, msg.MerkleRoot[:])
	if err != nil {
		return err
	}

	return nil
}

// SerializeSize returns the number of bytes it would take to serialize the
// block, factoring in any witness data within transaction.
func (msg *MsgVBlock) SerializeSize() int {
	// Version 4 bytes + BlockHash 32 bytes +
	// Serialized varint size for the number of vtransactions.
	n := 36 + serialization.VarIntSerializeSize(uint64(len(msg.VTransactions)))
	for _, tx := range msg.VTransactions {
		n += tx.SerializeSize()
	}

	return n
}
