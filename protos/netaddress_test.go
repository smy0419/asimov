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
	"testing"
	"time"

	"io"
)

// TestNetAddress tests the NetAddress API.
func TestNetAddress(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	port := 8333

	// Test NewNetAddress.
	na := NewNetAddress(&net.TCPAddr{IP: ip, Port: port}, 0)

	// Ensure we get the same ip, port, and services back out.
	if !na.IP.Equal(ip) {
		t.Errorf("NetNetAddress: wrong ip - got %v, want %v", na.IP, ip)
	}
	if na.Port != uint16(port) {
		t.Errorf("NetNetAddress: wrong port - got %v, want %v", na.Port,
			port)
	}
	if na.Services != 0 {
		t.Errorf("NetNetAddress: wrong services - got %v, want %v",
			na.Services, 0)
	}
	if na.HasService(common.SFNodeNetwork) {
		t.Errorf("HasService: SFNodeNetwork service is set")
	}

	// Ensure adding the full service node flag works.
	na.AddService(common.SFNodeNetwork)
	if na.Services != common.SFNodeNetwork {
		t.Errorf("AddService: wrong services - got %v, want %v",
			na.Services, common.SFNodeNetwork)
	}
	if !na.HasService(common.SFNodeNetwork) {
		t.Errorf("HasService: SFNodeNetwork service not set")
	}
}

// TestNetAddressWire tests the NetAddress protos encode and decode for various
// protocol versions and timestamp flag combinations.
func TestNetAddressWire(t *testing.T) {
	// baseNetAddr is used in the various tests as a baseline NetAddress.
	baseNetAddr := NetAddress{
		Timestamp: time.Unix(0x495fab29, 0), // 2009-01-03 12:15:05 -0600 CST
		Services:  common.SFNodeNetwork,
		IP:        net.ParseIP("127.0.0.1"),
		Port:      8333,
	}

	// baseNetAddrEncoded is the protos encoded bytes of baseNetAddr.
	baseNetAddrEncoded := []byte{
		0x29, 0xab, 0x5f, 0x49, 0x00, 0x00, 0x00, 0x00,// Timestamp
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // SFNodeNetwork
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xff, 0x7f, 0x00, 0x00, 0x01, // IP 127.0.0.1
		0x8d, 0x20, // Port 8333 in small-endian
	}

	tests := []struct {
		in   NetAddress // NetAddress to encode
		out  NetAddress // Expected decoded NetAddress
		buf  []byte     // Wire encoding
		pver uint32     // Protocol version for protos encoding
	}{
		{
			baseNetAddr,
			baseNetAddr,
			baseNetAddrEncoded,
			common.ProtocolVersion,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		var buf bytes.Buffer
		err := writeNetAddress(&buf, test.pver, &test.in)
		if err != nil {
			t.Errorf("writeNetAddress #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("writeNetAddress #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Decode the message from protos format.
		var na NetAddress
		rbuf := bytes.NewReader(test.buf)
		err = readNetAddress(rbuf, test.pver, &na)
		if err != nil {
			t.Errorf("readNetAddress #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(na, test.out) {
			t.Errorf("readNetAddress #%d\n got: %v want: %v", i,
				na, test.out)
			continue
		}
	}
}

// TestNetAddressWireErrors performs negative tests against protos encode and
// decode NetAddress to confirm error paths work correctly.
func TestNetAddressWireErrors(t *testing.T) {
	pver := common.ProtocolVersion

	// baseNetAddr is used in the various tests as a baseline NetAddress.
	baseNetAddr := NetAddress{
		Timestamp: time.Unix(0x495fab29, 0), // 2009-01-03 12:15:05 -0600 CST
		Services:  common.SFNodeNetwork,
		IP:        net.ParseIP("127.0.0.1"),
		Port:      8333,
	}

	tests := []struct {
		in       *NetAddress // Value to encode
		buf      []byte      // Wire encoding
		pver     uint32      // Protocol version for protos encoding
		ts       bool        // Include timestamp flag
		max      int         // Max size of fixed buffer to induce errors
		writeErr error       // Expected write error
		readErr  error       // Expected read error
	}{
		// Latest protocol version with timestamp and intentional
		// read/write errors.
		// Force errors on timestamp.
		{&baseNetAddr, []byte{}, pver, true, 0, io.ErrShortWrite, io.EOF},
		// Force errors on services.
		{&baseNetAddr, []byte{}, pver, true, 4, io.ErrShortWrite, io.ErrUnexpectedEOF},
		// Force errors on ip.
		{&baseNetAddr, []byte{}, pver, true, 12, io.ErrShortWrite, io.ErrUnexpectedEOF},
		// Force errors on port.
		{&baseNetAddr, []byte{}, pver, true, 28, io.ErrShortWrite, io.ErrUnexpectedEOF},

		// Latest protocol version with no timestamp and intentional
		// read/write errors.
		// Force errors on services.
		{&baseNetAddr, []byte{}, pver, false, 0, io.ErrShortWrite, io.EOF},
		// Force errors on ip.
		{&baseNetAddr, []byte{}, pver, false, 8, io.ErrShortWrite, io.EOF},
		// Force errors on port.
		{&baseNetAddr, []byte{}, pver, false, 24, io.ErrShortWrite, io.ErrUnexpectedEOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := writeNetAddress(w, test.pver, test.in)
		if err != test.writeErr {
			t.Errorf("writeNetAddress #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		var na NetAddress
		r := newFixedReader(test.max, test.buf)
		err = readNetAddress(r, test.pver, &na)
		if err != test.readErr {
			t.Errorf("readNetAddress #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}
