// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"io"

	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/serialization"
)

// MaxGetCFiltersReqRange the maximum number of filters that may be requested in
// a getcfheaders message.
const MaxGetCFiltersReqRange = 1000

// MsgGetCFilters implements the Message interface and represents a bitcoin
// getcfilters message. It is used to request committed filters for a range of
// blocks.
type MsgGetCFilters struct {
	FilterType  FilterType
	StartHeight uint32
	StopHash    common.Hash
}

// VVSDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgGetCFilters) VVSDecode(r io.Reader, pver uint32, _ MessageEncoding) error {
	if err := serialization.ReadUint8(r, (*uint8)(&msg.FilterType)); err != nil {
		return err
	}
	if err := serialization.ReadUint32(r, &msg.StartHeight); err != nil {
		return err
	}
	if err := serialization.ReadNBytes(r, msg.StopHash[:], common.HashLength); err != nil {
		return err
	}
	return nil
}

// VVSEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgGetCFilters) VVSEncode(w io.Writer, pver uint32, _ MessageEncoding) error {
	if err := serialization.WriteUint8(w, (uint8)(msg.FilterType)); err != nil {
		return err
	}
	if err := serialization.WriteUint32(w, msg.StartHeight); err != nil {
		return err
	}
	if err := serialization.WriteNBytes(w, msg.StopHash[:]); err != nil {
		return err
	}
	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgGetCFilters) Command() string {
	return CmdGetCFilters
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgGetCFilters) MaxPayloadLength(pver uint32) uint32 {
	// Filter type + uint32 + block hash
	return 1 + 4 + common.HashLength
}

// NewMsgGetCFilters returns a new bitcoin getcfilters message that conforms to
// the Message interface using the passed parameters and defaults for the
// remaining fields.
func NewMsgGetCFilters(filterType FilterType, startHeight uint32,
	stopHash *common.Hash) *MsgGetCFilters {
	return &MsgGetCFilters{
		FilterType:  filterType,
		StartHeight: startHeight,
		StopHash:    *stopHash,
	}
}
