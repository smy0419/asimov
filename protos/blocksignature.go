// Copyright (c) 2018-2020 The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"bytes"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/serialization"
	"io"
	"unsafe"
)

// MsgBlockSign describes relationship between block and signer
// signer signs for the block and ensures the block is valid
type MsgBlockSign struct {
	BlockHeight int32
	BlockHash   common.Hash
	Signer      common.Address
	Signature   [HashSignLen]byte
}

// const fixedBlockSignPayloadLen = 4 + 2 + 4 + common.HashLength + HashSignLen
const fixedBlockSignPayloadLen = 4 + common.HashLength + common.AddressLength + HashSignLen

// Serialize encodes MsgBlockSign to w using the asimov protocol encoding.
func (bs *MsgBlockSign) Serialize(w io.Writer) error {
	if err := serialization.WriteUint32(w, uint32(bs.BlockHeight)); err != nil {
		return err
	}
	if err := serialization.WriteNBytes(w, bs.BlockHash[:]); err != nil {
		return err
	}
	if err := serialization.WriteNBytes(w, bs.Signer[:]); err != nil {
		return err
	}
	if err := serialization.WriteNBytes(w, bs.Signature[:]); err != nil {
		return err
	}

	return nil
}

// Deserialize decodes MsgBlockSign from r into the receiver using a format
// that is suitable for long-term storage such as a database while respecting
// the Version field.
func (bs *MsgBlockSign) Deserialize(r io.Reader) error {
	if err := serialization.ReadUint32(r, (*uint32)(unsafe.Pointer(&bs.BlockHeight))); err != nil {
		return err
	}
	if err := serialization.ReadNBytes(r, bs.BlockHash[:], common.HashLength); err != nil {
		return err
	}
	if err := serialization.ReadNBytes(r, bs.Signer[:], common.AddressLength); err != nil {
		return err
	}
	if err := serialization.ReadNBytes(r, bs.Signature[:], HashSignLen); err != nil {
		return err
	}
	return nil
}

// SignHash generates the Hash for the block sign.
func (bs *MsgBlockSign) SignHash() common.Hash {
	buf := bytes.NewBuffer(make([]byte, 0, bs.SerializeSize()))
	_ = bs.Serialize(buf)
	return common.DoubleHashH(buf.Bytes())
}

// SerializeSize returns length of block sign payload
func (bs *MsgBlockSign) SerializeSize() int {
	return fixedBlockSignPayloadLen
}

// VVSDecode decodes MsgBlockSign from r into the receiver using a format
//// that is suitable for long-term storage such as a database while respecting
//// the Version field.
func (bs *MsgBlockSign) VVSDecode(r io.Reader, pver uint32, me MessageEncoding) error {
	return bs.Deserialize(r)
}

// VVSEncode encodes MsgBlockSign to w using the asimov protocol encoding
func (bs *MsgBlockSign) VVSEncode(w io.Writer, pver uint32, me MessageEncoding) error {
	return bs.Serialize(w)
}

func (bs *MsgBlockSign) MaxPayloadLength(uint32) uint32 {
	return fixedBlockSignPayloadLen
}

func (bs *MsgBlockSign) Command() string {
	return CmdSig
}

// MsgBlockSign records the information of signer who signs for block
type BlockSignList []*MsgBlockSign

/// Len returns the number of signers sign for block
func (bl BlockSignList) Len() int {
	return len(bl)
}

func (bl BlockSignList) Less(i, j int) bool {
	if bl[i].BlockHeight < bl[j].BlockHeight {
		return true
	}
	if bl[i].BlockHeight > bl[j].BlockHeight {
		return false
	}
	return bytes.Compare(bl[i].Signature[:], bl[j].Signature[:]) < 0
}

func (bl BlockSignList) Swap(i, j int) {
	bl[i], bl[j] = bl[j], bl[i]
}
