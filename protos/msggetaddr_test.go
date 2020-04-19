// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"bytes"
	"github.com/AsimovNetwork/asimov/common"
	"reflect"
	"testing"
)

// TestGetAddr tests the MsgGetAddr API.
func TestGetAddr(t *testing.T) {
	pver := common.ProtocolVersion

	// Ensure the command is expected value.
	wantCmd := "getaddr"
	msg := NewMsgGetAddr()
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgGetAddr: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value for latest protocol version.
	// Num addresses (varInt) + max allowed addresses.
	wantPayload := uint32(0)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload)
	}
}

// TestGetAddrWire tests the MsgGetAddr protos encode and decode for various
// protocol versions.
func TestGetAddrWire(t *testing.T) {
	msgGetAddr := NewMsgGetAddr()
	var msgGetAddrEncoded []byte

	tests := []struct {
		in   *MsgGetAddr     // Message to encode
		out  *MsgGetAddr     // Expected decoded message
		buf  []byte          // Wire encoding
		pver uint32          // Protocol version for protos encoding
		enc  MessageEncoding // Message encoding variant.
	}{
		// Latest protocol version.
		{
			msgGetAddr,
			msgGetAddr,
			msgGetAddrEncoded,
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
				buf.Bytes(), test.buf)
			continue
		}

		// Decode the message from protos format.
		var msg MsgGetAddr
		rbuf := bytes.NewReader(test.buf)
		err = msg.VVSDecode(rbuf, test.pver, test.enc)
		if err != nil {
			t.Errorf("VVSDecode #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(&msg, test.out) {
			t.Errorf("VVSDecode #%d\n got: %v want: %v", i,
				msg, test.out)
			continue
		}
	}
}
