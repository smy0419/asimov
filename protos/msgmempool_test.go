// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"bytes"
	"github.com/AsimovNetwork/asimov/common"
	"testing"
)

func TestMemPool(t *testing.T) {
	pver := common.ProtocolVersion
	enc := BaseEncoding

	// Ensure the command is expected value.
	wantCmd := "mempool"
	msg := NewMsgMemPool()
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgMemPool: wrong command - got %v want %v",
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
		t.Errorf("encode of MsgMemPool failed %v err <%v>", msg, err)
	}

	// Test decode with latest protocol version.
	readmsg := NewMsgMemPool()
	err = readmsg.VVSDecode(&buf, pver, enc)
	if err != nil {
		t.Errorf("decode of MsgMemPool failed [%v] err <%v>", buf, err)
	}
	//we do not care version, so the version parameter is no effect
	err = readmsg.VVSDecode(&buf, pver, enc)
	if err != nil {
		s := "decode of MsgMemPool passed for old protocol version %v err <%v>"
		t.Errorf(s, msg, err)
	}
}
