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

// TestSendHeaders tests the MsgSendHeaders API against the latest protocol
// version.
func TestSendHeaders(t *testing.T) {
	var SendHeadersVersion uint32 = 70012
	pver := common.ProtocolVersion
	enc := BaseEncoding

	// Ensure the command is expected value.
	wantCmd := "sendheaders"
	msg := NewMsgSendHeaders()
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgSendHeaders: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value.
	wantPayload := uint32(0)
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
		t.Errorf("encode of MsgSendHeaders failed %v err <%v>", msg,
			err)
	}

	// Older protocol versions should fail encode since message didn't
	// exist yet.
	//we do not care version
	oldPver := SendHeadersVersion - 1
	err = msg.VVSEncode(&buf, oldPver, enc)
	if err != nil {
		s := "encode of MsgSendHeaders passed for old protocol " +
			"version %v err <%v>"
		t.Errorf(s, msg, err)
	}

	// Test decode with latest protocol version.
	readmsg := NewMsgSendHeaders()
	err = readmsg.VVSDecode(&buf, pver, enc)
	if err != nil {
		t.Errorf("decode of MsgSendHeaders failed [%v] err <%v>", buf,
			err)
	}

	// Older protocol versions should fail decode since message didn't
	// exist yet.
	//we do not care version
	err = readmsg.VVSDecode(&buf, oldPver, enc)
	if err != nil {
		s := "decode of MsgSendHeaders passed for old protocol " +
			"version %v err <%v>"
		t.Errorf(s, msg, err)
	}
}


// TestSendHeadersCrossProtocol tests the MsgSendHeaders API when encoding with
// the latest protocol version and decoding with SendHeadersVersion.
func TestSendHeadersCrossProtocol(t *testing.T) {
	var SendHeadersVersion uint32 = 70012
	enc := BaseEncoding
	msg := NewMsgSendHeaders()

	// Encode with latest protocol version.
	var buf bytes.Buffer
	err := msg.VVSEncode(&buf, common.ProtocolVersion, enc)
	if err != nil {
		t.Errorf("encode of MsgSendHeaders failed %v err <%v>", msg,
			err)
	}

	// Decode with old protocol version.
	readmsg := NewMsgSendHeaders()
	err = readmsg.VVSDecode(&buf, SendHeadersVersion, enc)
	if err != nil {
		t.Errorf("decode of MsgSendHeaders failed [%v] err <%v>", buf,
			err)
	}
}

// TestSendHeadersWire tests the MsgSendHeaders protos encode and decode for
// various protocol versions.
func TestSendHeadersWire(t *testing.T) {
	msgSendHeaders := NewMsgSendHeaders()
	var msgSendHeadersEncoded []byte

	tests := []struct {
		in   *MsgSendHeaders // Message to encode
		out  *MsgSendHeaders // Expected decoded message
		buf  []byte          // Wire encoding
		pver uint32          // Protocol version for protos encoding
		enc  MessageEncoding // Message encoding format
	}{
		// Latest protocol version.
		{
			msgSendHeaders,
			msgSendHeaders,
			msgSendHeadersEncoded,
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
		var msg MsgSendHeaders
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
