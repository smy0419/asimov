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

// TestFilterCLearLatest tests the MsgFilterClear API against the latest
// protocol version.
func TestFilterClearLatest(t *testing.T) {
	pver := common.ProtocolVersion

	msg := NewMsgFilterClear()

	// Ensure the command is expected value.
	wantCmd := "filterclear"
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgFilterClear: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value for latest protocol version.
	wantPayload := uint32(0)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload)
	}
}

// TestFilterClearWire tests the MsgFilterClear protos encode and decode for
// various protocol versions.
func TestFilterClearWire(t *testing.T) {
	msgFilterClear := NewMsgFilterClear()
	var msgFilterClearEncoded []byte

	tests := []struct {
		in   *MsgFilterClear // Message to encode
		out  *MsgFilterClear // Expected decoded message
		buf  []byte          // Wire encoding
		pver uint32          // Protocol version for protos encoding
		enc  MessageEncoding // Message encoding format
	}{
		// Latest protocol version.
		{
			msgFilterClear,
			msgFilterClear,
			msgFilterClearEncoded,
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
		var msg MsgFilterClear
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
