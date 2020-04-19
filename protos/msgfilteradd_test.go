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
)

// TestFilterAddLatest tests the MsgFilterAdd API against the latest protocol
// version.
func TestFilterAddLatest(t *testing.T) {
	enc := BaseEncoding
	pver := common.ProtocolVersion

	data := []byte{0x01, 0x02}
	msg := NewMsgFilterAdd(data)

	// Ensure the command is expected value.
	wantCmd := "filteradd"
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgFilterAdd: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value for latest protocol version.
	wantPayload := uint32(523)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload)
	}

	// Test encode with latest protocol version.
	var buf bytes.Buffer
	err := msg.VVSEncode(&buf, pver, enc)
	if err != nil {
		t.Errorf("encode of MsgFilterAdd failed %v err <%v>", msg, err)
	}

	// Test decode with latest protocol version.
	var readmsg MsgFilterAdd
	err = readmsg.VVSDecode(&buf, pver, enc)
	if err != nil {
		t.Errorf("decode of MsgFilterAdd failed [%v] err <%v>", buf, err)
	}
}


// TestFilterAddMaxDataSize tests the MsgFilterAdd API maximum data size.
func TestFilterAddMaxDataSize(t *testing.T) {
	data := bytes.Repeat([]byte{0xff}, 521)
	msg := NewMsgFilterAdd(data)

	// Encode with latest protocol version.
	var buf bytes.Buffer
	err := msg.VVSEncode(&buf, common.ProtocolVersion, LatestEncoding)
	if err == nil {
		t.Errorf("encode of MsgFilterAdd succeeded when it shouldn't "+
			"have %v", msg)
	}

	// Decode with latest protocol version.
	readbuf := bytes.NewReader(data)
	err = msg.VVSDecode(readbuf, common.ProtocolVersion, LatestEncoding)
	if err == nil {
		t.Errorf("decode of MsgFilterAdd succeeded when it shouldn't "+
			"have %v", msg)
	}
}

// TestFilterAddWireErrors performs negative tests against protos encode and decode
// of MsgFilterAdd to confirm error paths work correctly.
func TestFilterAddWireErrors(t *testing.T) {
	pver := common.ProtocolVersion
	baseData := []byte{0x01, 0x02, 0x03, 0x04}
	baseFilterAdd := NewMsgFilterAdd(baseData)
	baseFilterAddEncoded := append([]byte{0x04}, baseData...)

	tests := []struct {
		in       *MsgFilterAdd   // Value to encode
		buf      []byte          // Wire encoding
		pver     uint32          // Protocol version for protos encoding
		enc      MessageEncoding // Message encoding format
		max      int             // Max size of fixed buffer to induce errors
		writeErr error           // Expected write error
		readErr  error           // Expected read error
	}{
		// Latest protocol version with intentional read/write errors.
		// Force error in data size.
		{
			baseFilterAdd, baseFilterAddEncoded, pver, BaseEncoding, 0,
			io.ErrShortWrite, io.EOF,
		},
		// Force error in data.
		{
			baseFilterAdd, baseFilterAddEncoded, pver, BaseEncoding, 1,
			io.ErrShortWrite, io.EOF,
		},
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
		var msg MsgFilterAdd
		r := newFixedReader(test.max, test.buf)
		err = msg.VVSDecode(r, test.pver, test.enc)
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
