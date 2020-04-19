// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"bytes"
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/serialization"
	"io"
	"strings"
	"time"
	"unsafe"
)

// MaxUserAgentLen is the maximum allowed length for the user agent field in a
// version message (MsgVersion).
const MaxUserAgentLen = 256

// DefaultUserAgent for protos in the stack
const DefaultUserAgent = "/asimovproto:0.0.1/"

// MsgVersion implements the Message interface and represents a bitcoin version
// message.  It is used for a peer to advertise itself as soon as an outbound
// connection is made.  The remote peer then uses this information along with
// its own to negotiate.  The remote peer must then respond with a version
// message of its own containing the negotiated values followed by a verack
// message (MsgVerAck).  This exchange must take place before any further
// communication is allowed to proceed.
type MsgVersion struct {
	// Version of the protocol the node is using.
	ProtocolVersion uint32

	// Bitfield which identifies the enabled services.
	Services common.ServiceFlag

	// Time the message was generated.  This is encoded as an int64 on the protos.
	Timestamp int64

	// 4 bytes
	Magic common.AsimovNet

	// ChainId is used to distinguish different chain, the main chain occupy zero,
	// each subchain take a positive integer
	ChainId uint64

	// Address of the remote peer.
	AddrYou NetAddress

	// Address of the local peer.
	AddrMe NetAddress

	// Unique value associated with message that is used to detect self
	// connections.
	Nonce uint64

	// The user agent that generated messsage.  This is an encoded as a varString
	// on the protos.  This has a max length of MaxUserAgentLen.
	UserAgent string

	// Last block seen by the generator of the version message.
	LastBlock int32

	// Don't announce transactions to peer.
	DisableRelayTx bool
}

// HasService returns whether the specified service is supported by the peer
// that generated the message.
func (msg *MsgVersion) HasService(service common.ServiceFlag) bool {
	return msg.Services&service == service
}

// AddService adds service as a supported service by the peer generating the
// message.
func (msg *MsgVersion) AddService(service common.ServiceFlag) {
	msg.Services |= service
}

// VVSDecode decodes r using the bitcoin protocol encoding into the receiver.
// The version message is special in that the protocol version hasn't been
// negotiated yet.  As a result, the pver field is ignored and any fields which
// are added in new versions are optional.  This also mean that r must be a
// *bytes.Buffer so the number of remaining bytes can be ascertained.
//
// This is part of the Message interface implementation.
func (msg *MsgVersion) VVSDecode(r io.Reader, pver uint32, enc MessageEncoding) error {
	buf, ok := r.(*bytes.Buffer)
	if !ok {
		return fmt.Errorf("MsgVersion.VVSDecode reader is not a " +
			"*bytes.Buffer")
	}

	if err := serialization.ReadUint32(buf, &msg.ProtocolVersion); err != nil {
		return err
	}
	if err := serialization.ReadUint64(buf, (*uint64)(&msg.Services)); err != nil {
		return err
	}
	if err := serialization.ReadUint64(buf, (*uint64)(unsafe.Pointer(&msg.Timestamp))); err != nil {
		return err
	}
	if err := serialization.ReadUint32(buf, (*uint32)(&msg.Magic)); err != nil {
		return err
	}
	if err := serialization.ReadUint64(buf, &msg.ChainId); err != nil {
		return err
	}

	if err := readNetAddress(buf, pver, &msg.AddrYou); err != nil {
		return err
	}

	if err := readNetAddress(buf, pver, &msg.AddrMe); err != nil {
		return err
	}
	if err := serialization.ReadUint64(buf, &msg.Nonce); err != nil {
		return err
	}
	userAgent, err := serialization.ReadVarString(buf, pver)
	if err != nil {
		return err
	}
	err = validateUserAgent(userAgent)
	if err != nil {
		return err
	}
	msg.UserAgent = userAgent

	if err := serialization.ReadUint32(buf, (*uint32)(unsafe.Pointer(&msg.LastBlock))); err != nil {
		return err
	}

	if err := serialization.ReadBool(r, &msg.DisableRelayTx); err != nil {
		return err
	}
	return nil
}

// VVSEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgVersion) VVSEncode(w io.Writer, pver uint32, enc MessageEncoding) error {
	err := validateUserAgent(msg.UserAgent)
	if err != nil {
		return err
	}

	if err := serialization.WriteUint32(w, msg.ProtocolVersion); err != nil {
		return err
	}
	if err := serialization.WriteUint64(w, uint64(msg.Services)); err != nil {
		return err
	}
	if err := serialization.WriteUint64(w, uint64(msg.Timestamp)); err != nil {
		return err
	}
	if err := serialization.WriteUint32(w, uint32(msg.Magic)); err != nil {
		return err
	}
	if err := serialization.WriteUint64(w, msg.ChainId); err != nil {
		return err
	}

	err = writeNetAddress(w, pver, &msg.AddrYou)
	if err != nil {
		return err
	}

	err = writeNetAddress(w, pver, &msg.AddrMe)
	if err != nil {
		return err
	}

	err = serialization.WriteUint64(w, msg.Nonce)
	if err != nil {
		return err
	}

	err = serialization.WriteVarString(w, pver, msg.UserAgent)
	if err != nil {
		return err
	}

	err = serialization.WriteUint32(w, uint32(msg.LastBlock))
	if err != nil {
		return err
	}

	err = serialization.WriteBool(w, msg.DisableRelayTx)
	if err != nil {
		return err
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgVersion) Command() string {
	return CmdVersion
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgVersion) MaxPayloadLength(pver uint32) uint32 {
	// Protocol version 4 bytes + services 8 bytes + timestamp 8 bytes +
	// magic 4 bytes + chainId 8 bytes +
	// remote and local net addresses + nonce 8 bytes + length of user
	// agent (varInt) + max allowed useragent length + last block 4 bytes +
	// relay transactions flag 1 byte.
	return 45 + (maxNetAddressPayload * 2) + serialization.MaxVarIntPayload +
		MaxUserAgentLen
}

// NewMsgVersion returns a new bitcoin version message that conforms to the
// Message interface using the passed parameters and defaults for the remaining
// fields.
func NewMsgVersion(me *NetAddress, you *NetAddress, nonce uint64,
	lastBlock int32, net common.AsimovNet) *MsgVersion {

	// Limit the timestamp to one second precision since the protocol
	// doesn't support better.
	return &MsgVersion{
		ProtocolVersion: common.ProtocolVersion,
		Services:        0,
		Timestamp:       time.Now().Unix(),
		AddrYou:         *you,
		AddrMe:          *me,
		Nonce:           nonce,
		UserAgent:       DefaultUserAgent,
		LastBlock:       lastBlock,
		DisableRelayTx:  false,
		Magic:           net,
	}
}

// validateUserAgent checks userAgent length against MaxUserAgentLen
func validateUserAgent(userAgent string) error {
	if len(userAgent) > MaxUserAgentLen {
		str := fmt.Sprintf("user agent too long [len %v, max %v]",
			len(userAgent), MaxUserAgentLen)
		return messageError("MsgVersion", str)
	}
	return nil
}

// AddUserAgent adds a user agent to the user agent string for the version
// message.  The version string is not defined to any strict format, although
// it is recommended to use the form "major.minor.revision" e.g. "2.6.41".
func (msg *MsgVersion) AddUserAgent(name string, version string,
	comments ...string) error {

	newUserAgent := fmt.Sprintf("%s:%s", name, version)
	if len(comments) != 0 {
		newUserAgent = fmt.Sprintf("%s(%s)", newUserAgent,
			strings.Join(comments, "; "))
	}
	newUserAgent = fmt.Sprintf("%s%s/", msg.UserAgent, newUserAgent)
	err := validateUserAgent(newUserAgent)
	if err != nil {
		return err
	}
	msg.UserAgent = newUserAgent
	return nil
}
