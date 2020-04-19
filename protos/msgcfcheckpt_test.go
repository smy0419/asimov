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

var testMsgCFCheckptHash = common.Hash{
	0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
	0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
	0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
	0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x05, // stopHash
}

var testFilterHeadersHash = []*common.Hash{
	{
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
		0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x06, // FilterHeaders
	},
	{
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
		0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x07, // FilterHeaders
	},
}


// TestMsgGetCFCheckpt tests the MsgGetCFCheckpt protos encode and decode for various
// numbers of headers and protocol versions.
func TestMsgCFCheckptSerialize(t *testing.T) {

	msgCFCheckpt := NewMsgCFCheckpt(1,&testMsgCFCheckptHash,2)
	msgCFCheckpt.AddCFHeader(testFilterHeadersHash[0])
	msgCFCheckpt.AddCFHeader(testFilterHeadersHash[1])
	msgCFCheckptEncoded := []byte{
		0x01,                   // filterType.
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
		0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x05, // stopHash
		0x02,
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
		0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x06, // FilterHeaders
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
		0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x07, // FilterHeaders
	}

	tests := []struct {
		in   *MsgCFCheckpt     // Message to encode
		out  *MsgCFCheckpt     // Expected decoded message
		buf  []byte          // Wire encoding
		pver uint32          // Protocol version for protos encoding
		enc  MessageEncoding // Message encoding format
	}{
		// Latest protocol version with no headers.
		{
			msgCFCheckpt,
			msgCFCheckpt,
			msgCFCheckptEncoded,
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
		var msg MsgCFCheckpt
		rbuf := bytes.NewReader(test.buf)
		err = msg.VVSDecode(rbuf, test.pver, test.enc)
		if err != nil {
			t.Errorf("VVSDecode #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(&msg, test.out) {
			t.Errorf("VVSDecode #%d\n got: %v want: %v", i,
				&msg, test.out)
			continue
		}
	}
}

// TestMsgGetCFCheckptErrors performs negative tests against protos encode and decode
// of TestMsgGetCFCheckptErrors to confirm error paths work correctly.
func TestMsgCFCheckptErrors(t *testing.T) {
	pver := common.ProtocolVersion

	msgCFCheckpt := NewMsgCFCheckpt(1,&testMsgCFCheckptHash,2)
	msgCFCheckpt.AddCFHeader(testFilterHeadersHash[0])
	msgCFCheckpt.AddCFHeader(testFilterHeadersHash[1])
	msgCFCheckptEncoded := []byte{
		0x01,                   // filterType.
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
		0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x05, // stopHash
		0x02,
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
		0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x06, // FilterHeaders
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
		0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x07, // FilterHeaders
	}

	tests := []struct {
		in       *MsgCFCheckpt     // Value to encode
		buf      []byte          // Wire encoding
		pver     uint32          // Protocol version for protos encoding
		enc      MessageEncoding // Message encoding format
		max      int             // Max size of fixed buffer to induce errors
		writeErr error           // Expected write error
		readErr  error           // Expected read error
	}{
		// Force error in filterType.
		{msgCFCheckpt, msgCFCheckptEncoded, pver, BaseEncoding, 0, io.ErrShortWrite, io.EOF},
		// Force error in stophash.
		{msgCFCheckpt, msgCFCheckptEncoded, pver, BaseEncoding, 1, io.ErrShortWrite, io.EOF},
		// Force error in FilterHeaders length.
		{msgCFCheckpt, msgCFCheckptEncoded, pver, BaseEncoding, 33, io.ErrShortWrite, io.EOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		t.Logf("=======================Running the %d test", i)
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
		var msg MsgCFCheckpt
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


func TestFilterHeadersOverflowErrors(t *testing.T) {
	pver := uint32(0)

	tests := []struct {
		buf  []byte          // Wire encoding
		pver uint32          // Protocol version for protos encoding
		enc  MessageEncoding // Message encoding format
		err  error           // Expected error
	}{
		// FilterHeaders that claims to have ~uint64(0) FilterHeaders.
		{
			[]byte{
				0x01,                   // filterType.
				0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
				0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
				0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
				0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x05, // stopHash
				0x02,
				0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
				0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
				0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
				0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x06, // FilterHeaders
			}, pver, BaseEncoding, io.EOF,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Decode from protos format.
		var msg MsgCFCheckpt
		r := bytes.NewReader(test.buf)
		err := msg.VVSDecode(r,test.pver,test.enc)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Deserialize #%d wrong error got: %v, want: %v",
				i, err, reflect.TypeOf(test.err))
			continue
		}
	}
}

