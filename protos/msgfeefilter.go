// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"github.com/AsimovNetwork/asimov/common/serialization"
	"io"
	"unsafe"
)

// MsgFeeFilter implements the Message interface and represents a bitcoin
// feeFilter message.  It is used to request the receiving peer does not
// announce any transactions below the specified minimum fee rate.
//
// This message was not added until protocol versions starting with
// FeeFilterVersion.
type MsgFeeFilter struct {
	MinPrice int32
}

// VVSDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgFeeFilter) VVSDecode(r io.Reader, pver uint32, enc MessageEncoding) error {
	if err := serialization.ReadUint32(r, (*uint32)(unsafe.Pointer(&msg.MinPrice))); err != nil {
		return err
	}
	return nil
}

// VVSEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgFeeFilter) VVSEncode(w io.Writer, pver uint32, enc MessageEncoding) error {
	if err := serialization.WriteUint32(w, uint32(msg.MinPrice)); err != nil {
		return err
	}
	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgFeeFilter) Command() string {
	return CmdFeeFilter
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgFeeFilter) MaxPayloadLength(pver uint32) uint32 {
	return 8
}

// NewMsgFeeFilter returns a new bitcoin feeFilter message that conforms to
// the Message interface.  See MsgFeeFilter for details.
func NewMsgFeeFilter(minPrice int32) *MsgFeeFilter {
	return &MsgFeeFilter{
		MinPrice: minPrice,
	}
}
