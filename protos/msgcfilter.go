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

// FilterType is used to represent a filter type.
type FilterType uint8

const (
	// GCSFilterRegular is the regular filter type.
	GCSFilterRegular FilterType = iota
)

const (
	// MaxCFilterDataSize is the maximum byte size of a committed filter.
	// The maximum size is currently defined as 256KiB.
	MaxCFilterDataSize = 256 * 1024
)

// MsgCFilter implements the Message interface and represents a bitcoin cfilter
// message. It is used to deliver a committed filter in response to a
// getcfilters (MsgGetCFilters) message.
type MsgCFilter struct {
	FilterType FilterType
	BlockHash  common.Hash
	Data       []byte
}

// VVSDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgCFilter) VVSDecode(r io.Reader, pver uint32, _ MessageEncoding) error {
	// Read filter type, block hash
	if err := serialization.ReadUint8(r, (*uint8)(&msg.FilterType)); err != nil {
		return err
	}
	if err := serialization.ReadNBytes(r, msg.BlockHash[:], common.HashLength); err != nil {
		return err
	}

	// Read filter data
	data, err := serialization.ReadVarBytes(r, pver, MaxCFilterDataSize,
		"cfilter data")
	if err != nil {
			return err
	}
	msg.Data = data
	return nil
}

// VVSEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgCFilter) VVSEncode(w io.Writer, pver uint32, _ MessageEncoding) error {
	size := len(msg.Data)
	if size > MaxCFilterDataSize {
		str := fmt.Sprintf("cfilter size too large for message "+
			"[size %v, max %v]", size, MaxCFilterDataSize)
		return messageError("MsgCFilter.VVSEncode", str)
	}
	if err := serialization.WriteUint8(w, uint8(msg.FilterType)); err != nil {
		return err
	}
	if err := serialization.WriteNBytes(w, msg.BlockHash[:]); err != nil {
		return err
	}

	return serialization.WriteVarBytes(w, pver, msg.Data)
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgCFilter) Command() string {
	return CmdCFilter
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgCFilter) MaxPayloadLength(pver uint32) uint32 {
	return uint32(serialization.VarIntSerializeSize(MaxCFilterDataSize)) +
		MaxCFilterDataSize + common.HashLength + 1
}

// NewMsgCFilter returns a new bitcoin cfilter message that conforms to the
// Message interface. See MsgCFilter for details.
func NewMsgCFilter(filterType FilterType, blockHash *common.Hash,
	data []byte) *MsgCFilter {
	return &MsgCFilter{
		FilterType: filterType,
		BlockHash:  *blockHash,
		Data:       data,
	}
}
