// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"bytes"
	"github.com/AsimovNetwork/asimov/common"
	"io"
	"reflect"
	"testing"

	"github.com/AsimovNetwork/asimov/common/serialization"
)

// TestPing tests the MsgPing API against the latest protocol version.
func TestPing(t *testing.T) {
	pver := common.ProtocolVersion

	// Ensure we get the same nonce back out.
	nonce, err := serialization.RandomUint64()
	if err != nil {
		t.Errorf("RandomUint64: Error generating nonce: %v", err)
	}
	msg := NewMsgPing(nonce)
	if msg.Nonce != nonce {
		t.Errorf("NewMsgPing: wrong nonce - got %v, want %v",
			msg.Nonce, nonce)
	}

	// Ensure the command is expected value.
	wantCmd := "ping"
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgPing: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value for latest protocol version.
	wantPayload := uint32(8)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload)
	}
}

// TestPingCrossProtocol tests the MsgPing API when encoding with the latest
// protocol version and decoding with ProtocolVersion.
func TestPingCrossProtocol(t *testing.T) {
	nonce, err := serialization.RandomUint64()
	if err != nil {
		t.Errorf("RandomUint64: Error generating nonce: %v", err)
	}
	msg := NewMsgPing(nonce)
	if msg.Nonce != nonce {
		t.Errorf("NewMsgPing: wrong nonce - got %v, want %v",
			msg.Nonce, nonce)
	}

	// Encode with latest protocol version.
	var buf bytes.Buffer
	err = msg.VVSEncode(&buf, common.ProtocolVersion, BaseEncoding)
	if err != nil {
		t.Errorf("encode of MsgPing failed %v err <%v>", msg, err)
	}

	// Decode with old protocol version.
	readmsg := NewMsgPing(0)
	err = readmsg.VVSDecode(&buf, common.ProtocolVersion, BaseEncoding)
	if err != nil {
		t.Errorf("decode of MsgPing failed [%v] err <%v>", buf, err)
	}
}

// TestPingWire tests the MsgPing protos encode and decode for various protocol
// versions.
func TestPingWire(t *testing.T) {
	tests := []struct {
		in   MsgPing         // Message to encode
		out  MsgPing         // Expected decoded message
		buf  []byte          // Wire encoding
		pver uint32          // Protocol version for protos encoding
		enc  MessageEncoding // Message encoding format
	}{
		// Latest protocol version.
		{
			MsgPing{Nonce: 123123}, // 0x1e0f3
			MsgPing{Nonce: 123123}, // 0x1e0f3
			[]byte{0xf3, 0xe0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00},
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
			t.Errorf("VVSEncode #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Decode the message from protos format.
		var msg MsgPing
		rbuf := bytes.NewReader(test.buf)
		err = msg.VVSDecode(rbuf, test.pver, test.enc)
		if err != nil {
			t.Errorf("VVSDecode #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(msg, test.out) {
			t.Errorf("VVSDecode #%d\n got: %v want: %v", i,
				msg, test.out)
			continue
		}
	}
}

// TestPingWireErrors performs negative tests against protos encode and decode
// of MsgPing to confirm error paths work correctly.
func TestPingWireErrors(t *testing.T) {
	pver := common.ProtocolVersion

	tests := []struct {
		in       *MsgPing        // Value to encode
		buf      []byte          // Wire encoding
		pver     uint32          // Protocol version for protos encoding
		enc      MessageEncoding // Message encoding format
		max      int             // Max size of fixed buffer to induce errors
		writeErr error           // Expected write error
		readErr  error           // Expected read error
	}{
		// Latest protocol version with intentional read/write errors.
		{
			&MsgPing{Nonce: 123123}, // 0x1e0f3
			[]byte{0xf3, 0xe0, 0x01, 0x00},
			pver,
			BaseEncoding,
			2,
			io.ErrShortWrite,
			io.ErrUnexpectedEOF,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := test.in.VVSEncode(w, test.pver, test.enc)
		if err != test.writeErr {
			t.Errorf("VVSEncode #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		var msg MsgPing
		r := newFixedReader(test.max, test.buf)
		err = msg.VVSDecode(r, test.pver, test.enc)
		if err != test.readErr {
			t.Errorf("VVSDecode #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}
