// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/serialization"
	"io"
)

// MsgGetCFHeaders is a message similar to MsgGetHeaders, but for committed
// filter headers. It allows to set the FilterType field to get headers in the
// chain of basic (0x00) or extended (0x01) headers.
type MsgGetCFHeaders struct {
	FilterType  FilterType
	StartHeight uint32
	StopHash    common.Hash
}

// VVSDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgGetCFHeaders) VVSDecode(r io.Reader, pver uint32, _ MessageEncoding) error {
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
func (msg *MsgGetCFHeaders) VVSEncode(w io.Writer, pver uint32, _ MessageEncoding) error {
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
func (msg *MsgGetCFHeaders) Command() string {
	return CmdGetCFHeaders
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgGetCFHeaders) MaxPayloadLength(pver uint32) uint32 {
	// Filter type + uint32 + block hash
	return 1 + 4 + common.HashLength
}

// NewMsgGetCFHeaders returns a new bitcoin getcfheader message that conforms to
// the Message interface using the passed parameters and defaults for the
// remaining fields.
func NewMsgGetCFHeaders(filterType FilterType, startHeight uint32,
	stopHash *common.Hash) *MsgGetCFHeaders {
	return &MsgGetCFHeaders{
		FilterType:  filterType,
		StartHeight: startHeight,
		StopHash:    *stopHash,
	}
}
