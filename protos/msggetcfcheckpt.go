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

// MsgGetCFCheckpt is a request for filter headers at evenly spaced intervals
// throughout the blockchain history. It allows to set the FilterType field to
// get headers in the chain of basic (0x00) or extended (0x01) headers.
type MsgGetCFCheckpt struct {
	FilterType FilterType
	StopHash   common.Hash
}

// VVSDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgGetCFCheckpt) VVSDecode(r io.Reader, pver uint32, _ MessageEncoding) error {
	err := serialization.ReadUint8(r, (*uint8)(&msg.FilterType))
	if err != nil {
		return err
	}

	return serialization.ReadNBytes(r, msg.StopHash[:], common.HashLength)
}

// VVSEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgGetCFCheckpt) VVSEncode(w io.Writer, pver uint32, _ MessageEncoding) error {
	if err := serialization.WriteUint8(w, (uint8)(msg.FilterType)); err != nil {
		return err
	}
	if err := serialization.WriteNBytes(w, msg.StopHash[:]); err != nil {
		return err
	}
	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgGetCFCheckpt) Command() string {
	return CmdGetCFCheckpt
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgGetCFCheckpt) MaxPayloadLength(pver uint32) uint32 {
	// Filter type + uint32 + block hash
	return 1 + common.HashLength
}

// NewMsgGetCFCheckpt returns a new bitcoin getcfcheckpt message that conforms
// to the Message interface using the passed parameters and defaults for the
// remaining fields.
func NewMsgGetCFCheckpt(filterType FilterType, stopHash *common.Hash) *MsgGetCFCheckpt {
	return &MsgGetCFCheckpt{
		FilterType: filterType,
		StopHash:   *stopHash,
	}
}
