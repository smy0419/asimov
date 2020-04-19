// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"io"
)

// MsgMemPool implements the Message interface and represents a bitcoin mempool
// message.  It is used to request a list of transactions still in the active
// memory pool of a relay.
//
// This message has no payload and was not added until protocol versions
// starting with BIP0035Version.
type MsgMemPool struct{}

// VVSDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgMemPool) VVSDecode(r io.Reader, pver uint32, enc MessageEncoding) error {
	return nil
}

// VVSEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgMemPool) VVSEncode(w io.Writer, pver uint32, enc MessageEncoding) error {
	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgMemPool) Command() string {
	return CmdMemPool
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgMemPool) MaxPayloadLength(pver uint32) uint32 {
	return 0
}

// NewMsgMemPool returns a new bitcoin pong message that conforms to the Message
// interface.  See MsgPong for details.
func NewMsgMemPool() *MsgMemPool {
	return &MsgMemPool{}
}
