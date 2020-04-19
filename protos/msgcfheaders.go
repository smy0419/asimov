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

const (
	// MaxCFHeaderPayload is the maximum byte size of a committed
	// filter header.
	MaxCFHeaderPayload = common.HashLength

	// MaxCFHeadersPerMsg is the maximum number of committed filter headers
	// that can be in a single bitcoin cfheaders message.
	MaxCFHeadersPerMsg = 2000
)

// MsgCFHeaders implements the Message interface and represents a bitcoin
// cfheaders message.  It is used to deliver committed filter header information
// in response to a getcfheaders message (MsgGetCFHeaders). The maximum number
// of committed filter headers per message is currently 2000. See
// MsgGetCFHeaders for details on requesting the headers.
type MsgCFHeaders struct {
	FilterType       FilterType
	StopHash         common.Hash
	PrevFilterHeader common.Hash
	FilterHashes     []*common.Hash
}

// AddCFHash adds a new filter hash to the message.
func (msg *MsgCFHeaders) AddCFHash(hash *common.Hash) error {
	if len(msg.FilterHashes)+1 > MaxCFHeadersPerMsg {
		str := fmt.Sprintf("too many block headers in message [max %v]",
			MaxBlockHeadersPerMsg)
		return messageError("MsgCFHeaders.AddCFHash", str)
	}

	msg.FilterHashes = append(msg.FilterHashes, hash)
	return nil
}

// VVSDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgCFHeaders) VVSDecode(r io.Reader, pver uint32, _ MessageEncoding) error {
	// Read filter type, stop hash, prev filter header
	if err := serialization.ReadUint8(r, (*uint8)(&msg.FilterType)); err != nil {
		return err
	}
	if err := serialization.ReadNBytes(r, msg.StopHash[:], common.HashLength); err != nil {
		return err
	}
	if err := serialization.ReadNBytes(r, msg.PrevFilterHeader[:], common.HashLength); err != nil {
		return err
	}

	// Read number of filter headers
	count, err := serialization.ReadVarInt(r, pver)
	if err != nil {
		return err
	}

	// Limit to max committed filter headers per message.
	if count > MaxCFHeadersPerMsg {
		str := fmt.Sprintf("too many committed filter headers for "+
			"message [count %v, max %v]", count,
			MaxBlockHeadersPerMsg)
		return messageError("MsgCFHeaders.VVSDecode", str)
	}

	// Create a contiguous slice of hashes to deserialize into in order to
	// reduce the number of allocations.
	msg.FilterHashes = make([]*common.Hash, 0, count)
	for i := uint64(0); i < count; i++ {
		var cfh common.Hash
		if err := serialization.ReadNBytes(r, cfh[:], common.HashLength); err != nil {
			return err
		}
		err = msg.AddCFHash(&cfh)
		if err != nil {
			return err
		}
	}

	return nil
}

// VVSEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgCFHeaders) VVSEncode(w io.Writer, pver uint32, _ MessageEncoding) error {
	// Write filter type, stop hash, prev filter header
	if err := serialization.WriteUint8(w, (uint8)(msg.FilterType)); err != nil {
		return err
	}
	if err := serialization.WriteNBytes(w, msg.StopHash[:]); err != nil {
		return err
	}
	if err := serialization.WriteNBytes(w, msg.PrevFilterHeader[:]); err != nil {
		return err
	}

	// Limit to max committed headers per message.
	count := len(msg.FilterHashes)
	if count > MaxCFHeadersPerMsg {
		str := fmt.Sprintf("too many committed filter headers for "+
			"message [count %v, max %v]", count,
			MaxBlockHeadersPerMsg)
		return messageError("MsgCFHeaders.VVSEncode", str)
	}

	err := serialization.WriteVarInt(w, pver, uint64(count))
	if err != nil {
		return err
	}

	for _, cfh := range msg.FilterHashes {
		err := serialization.WriteNBytes(w, cfh[:])
		if err != nil {
			return err
		}
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgCFHeaders) Command() string {
	return CmdCFHeaders
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver. This is part of the Message interface implementation.
func (msg *MsgCFHeaders) MaxPayloadLength(pver uint32) uint32 {
	// Hash size + filter type + num headers (varInt) +
	// (header size * max headers).
	return 1 + common.HashLength + common.HashLength + serialization.MaxVarIntPayload +
		(MaxCFHeaderPayload * MaxCFHeadersPerMsg)
}

// NewMsgCFHeaders returns a new bitcoin cfheaders message that conforms to
// the Message interface. See MsgCFHeaders for details.
func NewMsgCFHeaders() *MsgCFHeaders {
	return &MsgCFHeaders{
		FilterHashes: make([]*common.Hash, 0, MaxCFHeadersPerMsg),
	}
}
