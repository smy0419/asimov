// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"bytes"
	"github.com/AsimovNetwork/asimov/common"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/AsimovNetwork/asimov/common/serialization"
	"io"
)

// TestVersion tests the MsgVersion API.
func TestVersion(t *testing.T) {
	pver := common.ProtocolVersion

	// Create version message data.
	lastBlock := int32(234234)
	tcpAddrMe := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8333}
	me := NewNetAddress(tcpAddrMe, common.SFNodeNetwork)
	tcpAddrYou := &net.TCPAddr{IP: net.ParseIP("192.168.0.1"), Port: 8333}
	you := NewNetAddress(tcpAddrYou, common.SFNodeNetwork)
	nonce, err := serialization.RandomUint64()
	if err != nil {
		t.Errorf("RandomUint64: error generating nonce: %v", err)
	}

	// Ensure we get the correct data back out.
	msg := NewMsgVersion(me, you, nonce, lastBlock, common.MainNet)
	if msg.ProtocolVersion != pver {
		t.Errorf("NewMsgVersion: wrong protocol version - got %v, want %v",
			msg.ProtocolVersion, pver)
	}
	if !reflect.DeepEqual(&msg.AddrMe, me) {
		t.Errorf("NewMsgVersion: wrong me address - got %v, want %v",
			&msg.AddrMe, me)
	}
	if !reflect.DeepEqual(&msg.AddrYou, you) {
		t.Errorf("NewMsgVersion: wrong you address - got %v, want %v",
			&msg.AddrYou, you)
	}
	if msg.Nonce != nonce {
		t.Errorf("NewMsgVersion: wrong nonce - got %v, want %v",
			msg.Nonce, nonce)
	}
	if msg.UserAgent != DefaultUserAgent {
		t.Errorf("NewMsgVersion: wrong user agent - got %v, want %v",
			msg.UserAgent, DefaultUserAgent)
	}
	if msg.LastBlock != lastBlock {
		t.Errorf("NewMsgVersion: wrong last block - got %v, want %v",
			msg.LastBlock, lastBlock)
	}
	if msg.DisableRelayTx {
		t.Errorf("NewMsgVersion: disable relay tx is not false by "+
			"default - got %v, want %v", msg.DisableRelayTx, false)
	}

	msg.AddUserAgent("myclient", "1.2.3", "optional", "comments")
	customUserAgent := DefaultUserAgent + "myclient:1.2.3(optional; comments)/"
	if msg.UserAgent != customUserAgent {
		t.Errorf("AddUserAgent: wrong user agent - got %s, want %s",
			msg.UserAgent, customUserAgent)
	}

	msg.AddUserAgent("mygui", "3.4.5")
	customUserAgent += "mygui:3.4.5/"
	if msg.UserAgent != customUserAgent {
		t.Errorf("AddUserAgent: wrong user agent - got %s, want %s",
			msg.UserAgent, customUserAgent)
	}

	// accounting for ":", "/"
	err = msg.AddUserAgent(strings.Repeat("t",
		MaxUserAgentLen-len(customUserAgent)-2+1), "")
	if _, ok := err.(*MessageError); !ok {
		t.Errorf("AddUserAgent: expected error not received "+
			"- got %v, want %T", err, MessageError{})

	}

	// Version message should not have any services set by default.
	if msg.Services != 0 {
		t.Errorf("NewMsgVersion: wrong default services - got %v, want %v",
			msg.Services, 0)

	}
	if msg.HasService(common.SFNodeNetwork) {
		t.Errorf("HasService: SFNodeNetwork service is set")
	}

	// Ensure the command is expected value.
	wantCmd := "version"
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgVersion: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value.
	// Protocol version 4 bytes + services 8 bytes + timestamp 8 bytes +
	// magic 4 bytes + chainId 8 bytes +
	// remote and local net addresses + nonce 8 bytes + length of user agent
	// (varInt) + max allowed user agent length + last block 4 bytes +
	// relay transactions flag 1 byte.
	wantPayload := uint32(370)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload)
	}

	// Ensure adding the full service node flag works.
	msg.AddService(common.SFNodeNetwork)
	if msg.Services != common.SFNodeNetwork {
		t.Errorf("AddService: wrong services - got %v, want %v",
			msg.Services, common.SFNodeNetwork)
	}
	if !msg.HasService(common.SFNodeNetwork) {
		t.Errorf("HasService: SFNodeNetwork service not set")
	}
}

// TestVersionWire tests the MsgVersion protos encode and decode for various
// protocol versions.
func TestVersionWire(t *testing.T) {
	baseProtocolVersionCopy := *baseProtocolVersion
	verRelayTxFalse := &baseProtocolVersionCopy
	verRelayTxFalse.DisableRelayTx = true
	verRelayTxFalseEncoded := make([]byte, len(baseProtocolVersionEncoded))
	copy(verRelayTxFalseEncoded, baseProtocolVersionEncoded)
	verRelayTxFalseEncoded[len(verRelayTxFalseEncoded)-1] = 0

	tests := []struct {
		in   *MsgVersion     // Message to encode
		out  *MsgVersion     // Expected decoded message
		buf  []byte          // Wire encoding
		pver uint32          // Protocol version for protos encoding
		enc  MessageEncoding // Message encoding format
	}{
		// Latest protocol version.
		{
			baseProtocolVersion,
			baseProtocolVersion,
			baseProtocolVersionEncoded,
			common.ProtocolVersion,
			BaseEncoding,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode the message to protos format.
		var buf bytes.Buffer
		err := test.in.VVSEncode(&buf, test.pver, test.enc)
		if err != nil {
			t.Errorf("VVSEncode #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("VVSEncode #%d\n got: %s want: %s", i,
				spew.Sdump(buf.Bytes()), spew.Sdump(test.buf))
			continue
		}

		// Decode the message from protos format.
		var msg MsgVersion
		rbuf := bytes.NewBuffer(test.buf)
		err = msg.VVSDecode(rbuf, test.pver, test.enc)
		if err != nil {
			t.Errorf("VVSDecode #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(&msg, test.out) {
			t.Errorf("VVSDecode #%d\n got: %s want: %s", i,
				spew.Sdump(msg), spew.Sdump(test.out))
			continue
		}
	}
}

// TestVersionWireErrors performs negative tests against protos encode and
// decode of MsgGetHeaders to confirm error paths work correctly.
func TestVersionWireErrors(t *testing.T) {
	// Use protocol version 60002 specifically here instead of the latest
	// because the test data is using bytes encoded with that protocol
	// version.
	pver := uint32(1)
	enc := BaseEncoding
	wireErr := &MessageError{}

	// Ensure calling MsgVersion.VVSDecode with a non *bytes.Buffer returns
	// error.
	fr := newFixedReader(0, []byte{})
	if err := baseProtocolVersion.VVSDecode(fr, pver, enc); err == nil {
		t.Errorf("Did not received error when calling " +
			"MsgVersion.VVSDecode with non *bytes.Buffer")
	}

	// Copy the base version and change the user agent to exceed max limits.
	bvc := *baseProtocolVersion
	exceedUAVer := &bvc
	newUA := "/" + strings.Repeat("t", MaxUserAgentLen-8+1) + ":0.0.1/"
	exceedUAVer.UserAgent = newUA

	// Encode the new UA length as a varint.
	var newUAVarIntBuf bytes.Buffer
	err := serialization.WriteVarInt(&newUAVarIntBuf, pver, uint64(len(newUA)))
	if err != nil {
		t.Errorf("serialization.WriteVarInt: error %v", err)
	}

	// Make a new buffer big enough to hold the base version plus the new
	// bytes for the bigger varint to hold the new size of the user agent
	// and the new user agent string.  Then stich it all together.
	newLen := len(baseProtocolVersionEncoded) - len(baseProtocolVersion.UserAgent)
	newLen = newLen + len(newUAVarIntBuf.Bytes()) - 1 + len(newUA)
	exceedUAVerEncoded := make([]byte, newLen)
	copy(exceedUAVerEncoded, baseProtocolVersionEncoded[0:108])
	copy(exceedUAVerEncoded[108:], newUAVarIntBuf.Bytes())
	copy(exceedUAVerEncoded[111:], []byte(newUA))
	copy(exceedUAVerEncoded[111+len(newUA):], baseProtocolVersionEncoded[124:])

	tests := []struct {
		in       *MsgVersion     // Value to encode
		buf      []byte          // Wire encoding
		pver     uint32          // Protocol version for protos encoding
		enc      MessageEncoding // Message encoding format
		max      int             // Max size of fixed buffer to induce errors
		writeErr error           // Expected write error
		readErr  error           // Expected read error
	}{
		//Force error in protocol version.
		{baseProtocolVersion, baseProtocolVersionEncoded, pver, BaseEncoding, 0, io.ErrShortWrite, io.EOF}, //4byte
		// Force error in services.
		{baseProtocolVersion, baseProtocolVersionEncoded, pver, BaseEncoding, 4, io.ErrShortWrite, io.EOF}, //8byte
		// Force error in timestamp.
		{baseProtocolVersion, baseProtocolVersionEncoded, pver, BaseEncoding, 12, io.ErrShortWrite, io.EOF},  //8byte
		// Force error in Magic.
		{baseProtocolVersion, baseProtocolVersionEncoded, pver, BaseEncoding, 20, io.ErrShortWrite, io.EOF},  //4byte
		// Force error in chainid.
		{baseProtocolVersion, baseProtocolVersionEncoded, pver, BaseEncoding, 24, io.ErrShortWrite, io.EOF},  //8byte
		// Force error in remote address.
		{baseProtocolVersion, baseProtocolVersionEncoded, pver, BaseEncoding, 32, io.ErrShortWrite, io.EOF},  //34byte
		// Force error in local address.
		{baseProtocolVersion, baseProtocolVersionEncoded, pver, BaseEncoding, 66, io.ErrShortWrite, io.EOF},  //34byte
		// Force error in nonce.
		{baseProtocolVersion, baseProtocolVersionEncoded, pver, BaseEncoding, 100, io.ErrShortWrite, io.EOF}, //8byte
		// Force error in user agent.
		{baseProtocolVersion, baseProtocolVersionEncoded, pver, BaseEncoding, 108, io.ErrShortWrite, io.EOF}, //16byte + 1byte len
		// Force error in last block.
		{baseProtocolVersion, baseProtocolVersionEncoded, pver, BaseEncoding, 125, io.ErrShortWrite, io.EOF}, //4byte
		// Force error in DisableRelayTx.
		{baseProtocolVersion, baseProtocolVersionEncoded, pver, BaseEncoding, 129, io.ErrShortWrite, io.EOF},  //1byte
		//Force error due to user agent too big
		{exceedUAVer, exceedUAVerEncoded, pver, BaseEncoding, newLen, wireErr, wireErr},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := test.in.VVSEncode(w, test.pver, test.enc)
		if reflect.TypeOf(err) != reflect.TypeOf(test.writeErr) {
			t.Errorf("VVSEncode #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// For errors which are not of type MessageError, check them for
		// equality.
		if _, ok := err.(*MessageError); !ok {
			if err != test.writeErr {
				t.Errorf("VVSEncode #%d wrong error got: %v, "+
					"want: %v", i, err, test.writeErr)
				continue
			}
		}

		// Decode from protos format.
		var msg MsgVersion
		buf := bytes.NewBuffer(test.buf[0:test.max])
		err = msg.VVSDecode(buf, test.pver, test.enc)
		if reflect.TypeOf(err) != reflect.TypeOf(test.readErr) {
			t.Errorf("VVSDecode #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}

		// For errors which are not of type MessageError, check them for
		// equality.
		if _, ok := err.(*MessageError); !ok {
			if err != test.readErr {
				t.Errorf("VVSDecode #%d wrong error got: %v, "+
					"want: %v", i, err, test.readErr)
				continue
			}
		}
	}
}

// baseProtocolVersion is used in the various tests as a baseline MsgVersion
var baseProtocolVersion = &MsgVersion{
	ProtocolVersion: 1,
	Services:        common.SFNodeNetwork,
	Timestamp:       0x495fab29, // 2009-01-03 12:15:05 -0600 CST)
	Magic:           common.AsimovNet(0),
	ChainId:         1,
	AddrYou: NetAddress{
		Timestamp: time.Unix(0x495fab29, 0),
		Services:  common.SFNodeNetwork,
		IP:        net.ParseIP("192.168.0.1"),
		Port:      8333,
	},
	AddrMe: NetAddress{
		Timestamp: time.Unix(0x495fab29, 0),
		Services:  common.SFNodeNetwork,
		IP:        net.ParseIP("127.0.0.1"),
		Port:      8333,
	},
	Nonce:     123123, // 0x1e0f3
	UserAgent: "/btcdtest:0.0.1/",
	LastBlock: 234234, // 0x392fa
	DisableRelayTx:false,
}

// baseProtocolVersionEncoded is the protos encoded bytes for baseProtocolVersion
// used in the various tests.
var baseProtocolVersionEncoded = []byte{
	0x01, 0x00, 0x00, 0x00, // Protocol version
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // SFNodeNetwork
	0x29, 0xab, 0x5f, 0x49, 0x00, 0x00, 0x00, 0x00, // 64-bit Timestamp
	0x00, 0x00, 0x00, 0x00, // Magic
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // ChainId
	// AddrYou -- No timestamp for NetAddress in version message
	//0x00, 0x09, 0x6e, 0x88, 0xf1, 0xff, 0xff, 0xff, //timestamp
	0x29, 0xab, 0x5f, 0x49, 0x00, 0x00, 0x00, 0x00, // 64-bit Timestamp
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // SFNodeNetwork
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xc0, 0xa8, 0x00, 0x01, // IP 192.168.0.1
	0x8d, 0x20,  // Port 8333
	// AddrMe -- No timestamp for NetAddress in version message
	//0x00, 0x09, 0x6e, 0x88, 0xf1, 0xff, 0xff, 0xff, //timestamp
	0x29, 0xab, 0x5f, 0x49, 0x00, 0x00, 0x00, 0x00, // 64-bit Timestamp
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // SFNodeNetwork
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x7f, 0x00, 0x00, 0x01, // IP 127.0.0.1
	0x8d, 0x20, // Port 8333
	0xf3, 0xe0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, // Nonce
	0x10, // Varint for user agent length
	0x2f, 0x62, 0x74, 0x63, 0x64, 0x74, 0x65, 0x73,
	0x74, 0x3a, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x2f, // User agent
	0xfa, 0x92, 0x03, 0x00, // Last block
	0x00, // Relay tx
}
