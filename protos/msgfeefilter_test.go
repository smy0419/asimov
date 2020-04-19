// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"bytes"
	"github.com/AsimovNetwork/asimov/common"
	"io"
	"math/rand"
	"reflect"
	"testing"
)

// TestFeeFilterLatest tests the MsgFeeFilter API against the latest protocol version.
func TestFeeFilterLatest(t *testing.T) {
	pver := common.ProtocolVersion

	minPrice := rand.Int31()
	msg := NewMsgFeeFilter(minPrice)
	if msg.MinPrice != minPrice {
		t.Errorf("NewMsgFeeFilter: wrong minfee - got %v, want %v",
			msg.MinPrice, minPrice)
	}

	// Ensure the command is expected value.
	wantCmd := "feefilter"
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgFeeFilter: wrong command - got %v want %v",
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

	// Test encode with latest protocol version.
	var buf bytes.Buffer
	err := msg.VVSEncode(&buf, pver, BaseEncoding)
	if err != nil {
		t.Errorf("encode of MsgFeeFilter failed %v err <%v>", msg, err)
	}

	// Test decode with latest protocol version.
	readmsg := NewMsgFeeFilter(0)
	err = readmsg.VVSDecode(&buf, pver, BaseEncoding)
	if err != nil {
		t.Errorf("decode of MsgFeeFilter failed [%v] err <%v>", buf, err)
	}

	// Ensure minfee is the same.
	if msg.MinPrice != readmsg.MinPrice {
		t.Errorf("Should get same minfee for protocol version %d", pver)
	}
}

// TestFeeFilterWire tests the MsgFeeFilter protos encode and decode for various protocol
// versions.
func TestFeeFilterWire(t *testing.T) {
	tests := []struct {
		in   MsgFeeFilter // Message to encode
		out  MsgFeeFilter // Expected decoded message
		buf  []byte       // Wire encoding
		pver uint32       // Protocol version for protos encoding
	}{
		// Latest protocol version.
		{
			MsgFeeFilter{MinPrice: 123123}, // 0x1e0f3
			MsgFeeFilter{MinPrice: 123123}, // 0x1e0f3
			[]byte{0xf3, 0xe0, 0x01, 0x00},
			common.ProtocolVersion,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode the message to protos format.
		var buf bytes.Buffer
		err := test.in.VVSEncode(&buf, test.pver, BaseEncoding)
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
		var msg MsgFeeFilter
		rbuf := bytes.NewReader(test.buf)
		err = msg.VVSDecode(rbuf, test.pver, BaseEncoding)
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

// TestFeeFilterWireErrors performs negative tests against protos encode and decode
// of MsgFeeFilter to confirm error paths work correctly.
func TestFeeFilterWireErrors(t *testing.T) {
	pver := common.ProtocolVersion
	baseFeeFilter := NewMsgFeeFilter(123123) // 0x1e0f3
	baseFeeFilterEncoded := []byte{
		0xf3, 0xe0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	tests := []struct {
		in       *MsgFeeFilter // Value to encode
		buf      []byte        // Wire encoding
		pver     uint32        // Protocol version for protos encoding
		max      int           // Max size of fixed buffer to induce errors
		writeErr error         // Expected write error
		readErr  error         // Expected read error
	}{
		// Latest protocol version with intentional read/write errors.
		// Force error in minfee.
		{baseFeeFilter, baseFeeFilterEncoded, pver, 0, io.ErrShortWrite, io.EOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := test.in.VVSEncode(w, test.pver, BaseEncoding)
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
		var msg MsgFeeFilter
		r := newFixedReader(test.max, test.buf)
		err = msg.VVSDecode(r, test.pver, BaseEncoding)
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
