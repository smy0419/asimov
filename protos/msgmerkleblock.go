// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"fmt"
	"io"

	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/serialization"
)

// maxFlagsPerMerkleBlock is the maximum number of flag bytes that could
// possibly fit into a merkle block.  Since each transaction is represented by
// a single bit, this is the max number of transactions per block divided by
// 8 bits per byte.  Then an extra one to cover partials.
const maxFlagsPerMerkleBlock = maxTxPerBlock / 8

// MsgMerkleBlock implements the Message interface and represents a bitcoin
// merkleblock message which is used to reset a Bloom filter.
//
// This message was not added until protocol version BIP0037Version.
type MsgMerkleBlock struct {
	Header       BlockHeader
	Transactions uint32
	Hashes       []*common.Hash
	Flags        []byte
}

// AddTxHash adds a new transaction hash to the message.
func (msg *MsgMerkleBlock) AddTxHash(hash *common.Hash) error {
	if len(msg.Hashes)+1 > maxTxPerBlock {
		str := fmt.Sprintf("too many tx hashes for message [max %v]",
			maxTxPerBlock)
		return messageError("MsgMerkleBlock.AddTxHash", str)
	}

	msg.Hashes = append(msg.Hashes, hash)
	return nil
}

// VVSDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgMerkleBlock) VVSDecode(r io.Reader, pver uint32, enc MessageEncoding) error {
	err := readBlockHeader(r, &msg.Header)
	if err != nil {
		return err
	}

	err = serialization.ReadUint32(r, &msg.Transactions)
	if err != nil {
		return err
	}

	// Read num block locator hashes and limit to max.
	count, err := serialization.ReadVarInt(r, pver)
	if err != nil {
		return err
	}
	if count > maxTxPerBlock {
		str := fmt.Sprintf("too many transaction hashes for message "+
			"[count %v, max %v]", count, maxTxPerBlock)
		return messageError("MsgMerkleBlock.VVSDecode", str)
	}

	// Create a contiguous slice of hashes to deserialize into in order to
	// reduce the number of allocations.
	hashes := make([]common.Hash, count)
	msg.Hashes = make([]*common.Hash, 0, count)
	for i := uint64(0); i < count; i++ {
		hash := &hashes[i]
		err := serialization.ReadNBytes(r, hash[:], common.HashLength)
		if err != nil {
			return err
		}
		err = msg.AddTxHash(hash)
		if err != nil {
			return err
		}
	}

	msg.Flags, err = serialization.ReadVarBytes(r, pver, maxFlagsPerMerkleBlock,
		"merkle block flags size")
	return err
}

// VVSEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgMerkleBlock) VVSEncode(w io.Writer, pver uint32, enc MessageEncoding) error {
	// Read num transaction hashes and limit to max.
	numHashes := len(msg.Hashes)
	if numHashes > maxTxPerBlock {
		str := fmt.Sprintf("too many transaction hashes for message "+
			"[count %v, max %v]", numHashes, maxTxPerBlock)
		return messageError("MsgMerkleBlock.VVSDecode", str)
	}
	numFlagBytes := len(msg.Flags)
	if numFlagBytes > maxFlagsPerMerkleBlock {
		str := fmt.Sprintf("too many flag bytes for message [count %v, "+
			"max %v]", numFlagBytes, maxFlagsPerMerkleBlock)
		return messageError("MsgMerkleBlock.VVSDecode", str)
	}

	err := writeBlockHeader(w, &msg.Header)
	if err != nil {
		return err
	}

	err = serialization.WriteUint32(w, msg.Transactions)
	if err != nil {
		return err
	}

	err = serialization.WriteVarInt(w, pver, uint64(numHashes))
	if err != nil {
		return err
	}
	for _, hash := range msg.Hashes {
		err = serialization.WriteNBytes(w, hash[:])
		if err != nil {
			return err
		}
	}

	return serialization.WriteVarBytes(w, pver, msg.Flags)
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgMerkleBlock) Command() string {
	return CmdMerkleBlock
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgMerkleBlock) MaxPayloadLength(pver uint32) uint32 {
	return MaxBlockPayload
}

// NewMsgMerkleBlock returns a new bitcoin merkleblock message that conforms to
// the Message interface.  See MsgMerkleBlock for details.
func NewMsgMerkleBlock(bh *BlockHeader) *MsgMerkleBlock {
	return &MsgMerkleBlock{
		Header:       *bh,
		Transactions: 0,
		Hashes:       make([]*common.Hash, 0),
		Flags:        make([]byte, 0),
	}
}
