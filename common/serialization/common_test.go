// Copyright (c) 2018-2020. The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package serialization

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

const(
	MaxMessagePayload = 2 * 1024 * 1024
	ProtocolVersion = uint32(0)
)

// fakeRandReader implements the io.Reader interface and is used to force
// errors in the RandomUint64 function.
type fakeRandReader struct {
	n   int
	err error
}

// Read returns the fake reader error and the lesser of the fake reader value
// and the length of p.
func (r *fakeRandReader) Read(p []byte) (int, error) {
	n := r.n
	if n > len(p) {
		n = len(p)
	}
	return n, r.err
}

// TestBool tests protos encode and decode for bool type various elements.
func TestBool(t *testing.T) {
	tests := []struct {
		in  bool   // Value to encode
		buf []byte // serialized
	}{
		{false, []byte{0x00}},
		{true, []byte{0x01}},
	}

	t.Logf("Running bool %d tests", len(tests))
	for i, test := range tests {
		// Write to protos format.
		var buf bytes.Buffer
		err := WriteBool(&buf, test.in)
		if err != nil {
			t.Errorf("WriteBool #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteBool #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Read from protos format.
		rbuf := bytes.NewReader(test.buf)
		var val bool
		err = ReadBool(rbuf, &val)
		if err != nil {
			t.Errorf("ReadBool #%d error %v", i, err)
			continue
		}
		if val != test.in {
			t.Errorf("ReadBool #%d\n got: %v want: %v", i,
				val, test.in)
			continue
		}
	}
}

// TestBoolErrors performs negative tests against protos encode and decode
// of various element types to confirm error paths work correctly.
func TestBoolErrors(t *testing.T) {
	tests := []struct {
		in       bool // Value to encode
		max      int         // Max size of fixed buffer to induce errors
		writeErr error       // Expected write error
		readErr  error       // Expected read error
	}{
		{false, 0, io.ErrShortWrite, io.EOF},
		{true, 0, io.ErrShortWrite, io.EOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := WriteBool(w, test.in)
		if err != test.writeErr {
			t.Errorf("WriteBool #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		r := newFixedReader(test.max, nil)
		var val bool
		err = ReadBool(r, &val)
		if err != test.readErr {
			t.Errorf("ReadBool #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

// TestUint8 tests protos encode and decode for byte type various elements.
func TestUint8(t *testing.T) {
	tests := []struct {
		in   uint8   // Value to encode
		buf  []byte // serialized
	}{
		{1, []byte{0x01}},
		{128, []byte{0x80}},
		{255, []byte{0xff}},
	}

	t.Logf("Running uint8 %d tests", len(tests))
	for i, test := range tests {
		var buf bytes.Buffer
		err := WriteUint8(&buf, test.in)
		if err != nil {
			t.Errorf("WriteUint8 #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteUint8 #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Read from protos format.
		rbuf := bytes.NewReader(test.buf)
		var val uint8
		err = ReadUint8(rbuf, &val)
		if err != nil {
			t.Errorf("ReadUint8 #%d error %v", i, err)
			continue
		}
		if val != test.in {
			t.Errorf("ReadUint8 #%d\n got: %v want: %v", i,
				val, test.in)
			continue
		}
	}
}

// TestUint8Errors performs negative tests against protos encode and decode
// of various element types to confirm error paths work correctly.
func TestUint8Errors(t *testing.T) {
	tests := []struct {
		in       uint8 // Value to encode
		max      int         // Max size of fixed buffer to induce errors
		writeErr error       // Expected write error
		readErr  error       // Expected read error
	}{
		{1, 0,io.ErrShortWrite, io.EOF},
		{128, 0,io.ErrShortWrite, io.EOF},
		{255, 0,io.ErrShortWrite, io.EOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := WriteUint8(w, test.in)
		if err != test.writeErr {
			t.Errorf("WriteUint8 #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		r := newFixedReader(test.max, nil)
		var val uint8
		err = ReadUint8(r, &val)
		if err != test.readErr {
			t.Errorf("ReadUint8 #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}


// TestUint16 tests protos encode and decode for byte type various elements.
func TestUint16(t *testing.T) {
	tests := []struct {
		in   uint16   // Value to encode
		buf  []byte // serialized
	}{
		{1, []byte{0x01,0x00}},     // Min single byte
		{255, []byte{0xff,0x00}},   // Max single byte
		{256, []byte{0x00,0x01}},   // Min 2-byte
		{65535, []byte{0xff,0xff}}, // Max 2-byte
	}

	t.Logf("Running uint16 %d tests", len(tests))
	for i, test := range tests {
		var buf bytes.Buffer
		err := WriteUint16(&buf, test.in)
		if err != nil {
			t.Errorf("WriteUint16 #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteUint16 #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Read from protos format.
		rbuf := bytes.NewReader(test.buf)
		var val uint16
		err = ReadUint16(rbuf, &val)
		if err != nil {
			t.Errorf("ReadUint16 #%d error %v", i, err)
			continue
		}
		if val != test.in {
			t.Errorf("ReadUint16 #%d\n got: %v want: %v", i,
				val, test.in)
			continue
		}
	}
}

// TestUint16Errors performs negative tests against protos encode and decode
// of various element types to confirm error paths work correctly.
func TestUint16Errors(t *testing.T) {
	tests := []struct {
		in       uint16 // Value to encode
		max      int         // Max size of fixed buffer to induce errors
		writeErr error       // Expected write error
		readErr  error       // Expected read error
	}{
		{1, 0,io.ErrShortWrite, io.EOF},
		{256, 0,io.ErrShortWrite, io.EOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := WriteUint16(w, test.in)
		if err != test.writeErr {
			t.Errorf("WriteUint16 #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		r := newFixedReader(test.max, nil)
		var val uint16
		err = ReadUint16(r, &val)
		if err != test.readErr {
			t.Errorf("ReadUint16 #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}


// TestUint32 tests protos encode and decode for byte type various elements.
func TestUint32(t *testing.T) {
	tests := []struct {
		in   uint32   // Value to encode
		buf  []byte // serialized
	}{
		{1, []byte{0x01,0x00,0x00,0x00}},     // Min single byte
		{255, []byte{0xff,0x00,0x00,0x00}},   // Max single byte
		{256, []byte{0x00,0x01,0x00,0x00}},   // Min 2-byte
		{65535, []byte{0xff,0xff,0x00,0x00}}, // Max 2-byte
		{0x10000, []byte{0x00,0x00,0x01,0x00}},   // Min 4-byte
		{0xffffffff, []byte{0xff,0xff,0xff,0xff}}, // Max 4-byte
	}

	t.Logf("Running uint32 %d tests", len(tests))
	for i, test := range tests {
		var buf bytes.Buffer
		err := WriteUint32(&buf, test.in)
		if err != nil {
			t.Errorf("WriteUint32 #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteUint32 #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Read from protos format.
		rbuf := bytes.NewReader(test.buf)
		var val uint32
		err = ReadUint32(rbuf, &val)
		if err != nil {
			t.Errorf("ReadUint32 #%d error %v", i, err)
			continue
		}
		if val != test.in {
			t.Errorf("ReadUint32 #%d\n got: %v want: %v", i,
				val, test.in)
			continue
		}
	}
}

// TestUint32Errors performs negative tests against protos encode and decode
// of various element types to confirm error paths work correctly.
func TestUint32Errors(t *testing.T) {
	tests := []struct {
		in       uint32 // Value to encode
		max      int         // Max size of fixed buffer to induce errors
		writeErr error       // Expected write error
		readErr  error       // Expected read error
	}{
		{1, 0,io.ErrShortWrite, io.EOF}, // Force errors on 1-byte read/write.
		{256, 0,io.ErrShortWrite, io.EOF},  // Force errors on 2-byte read/write.
		{0x10000, 2, io.ErrShortWrite, io.ErrUnexpectedEOF}, // Force errors on 4-byte read/write.
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := WriteUint32(w, test.in)
		if err != test.writeErr {
			t.Errorf("WriteUint32 #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		r := newFixedReader(test.max, nil)
		var val uint32
		err = ReadUint32(r, &val)
		if err != test.readErr {
			t.Errorf("ReadUint32 #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

// TestUint64 tests protos encode and decode for byte type various elements.
func TestUint64(t *testing.T) {
	tests := []struct {
		in   uint64   // Value to encode
		buf  []byte // serialized
	}{
		{1, []byte{0x01,0x00,0x00,0x00,0x00,0x00,0x00,0x00}},     // Min single byte
		{255, []byte{0xff,0x00,0x00,0x00,0x00,0x00,0x00,0x00}},   // Max single byte
		{256, []byte{0x00,0x01,0x00,0x00,0x00,0x00,0x00,0x00}},   // Min 2-byte
		{65535, []byte{0xff,0xff,0x00,0x00,0x00,0x00,0x00,0x00}}, // Max 2-byte
		{0x10000, []byte{0x00,0x00,0x01,0x00,0x00,0x00,0x00,0x00}},   // Min 4-byte
		{0xffffffff, []byte{0xff,0xff,0xff,0xff,0x00,0x00,0x00,0x00}}, // Max 4-byte
		{0x100000000, []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}},   // Min 8-byte
		{0xffffffffffffffff, []byte{0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff}}, // Max 8-byte
	}

	t.Logf("Running uint64 %d tests", len(tests))
	for i, test := range tests {
		var buf bytes.Buffer
		err := WriteUint64(&buf, test.in)
		if err != nil {
			t.Errorf("WriteUint64 #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteUint64 #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Read from protos format.
		rbuf := bytes.NewReader(test.buf)
		var val uint64
		err = ReadUint64(rbuf, &val)
		if err != nil {
			t.Errorf("ReadUint64 #%d error %v", i, err)
			continue
		}
		if val != test.in {
			t.Errorf("ReadUint64 #%d\n got: %v want: %v", i,
				val, test.in)
			continue
		}
	}
}

// TestUint64Errors performs negative tests against protos encode and decode
// of various element types to confirm error paths work correctly.
func TestUint64Errors(t *testing.T) {
	tests := []struct {
		in       uint64 // Value to encode
		max      int         // Max size of fixed buffer to induce errors
		writeErr error       // Expected write error
		readErr  error       // Expected read error
	}{
		{1, 0,io.ErrShortWrite, io.EOF}, // Force errors on 1-byte read/write.
		{256, 0,io.ErrShortWrite, io.EOF},  // Force errors on 2-byte read/write.
		{0x10000, 2, io.ErrShortWrite, io.ErrUnexpectedEOF}, // Force errors on 4-byte read/write.
		{0x100000000, 4, io.ErrShortWrite, io.ErrUnexpectedEOF}, // Force errors on 8-byte read/write.
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := WriteUint64(w, test.in)
		if err != test.writeErr {
			t.Errorf("WriteUint64 #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		r := newFixedReader(test.max, nil)
		var val uint64
		err = ReadUint64(r, &val)
		if err != test.readErr {
			t.Errorf("ReadUint64 #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}


// TestUint32 tests protos encode and decode for byte type various elements.
func TestUint32B(t *testing.T) {
	tests := []struct {
		in   uint32   // Value to encode
		buf  []byte // serialized
	}{
		{1, []byte{0x00,0x00,0x00,0x01}},     // Min single byte
		{255, []byte{0x00,0x00,0x00,0xff}},   // Max single byte
		{256, []byte{0x00,0x00,0x01,0x00}},   // Min 2-byte
		{65535, []byte{0x00,0x00,0xff,0xff}}, // Max 2-byte
		{0x10000, []byte{0x00,0x01,0x00,0x00}},   // Min 4-byte
		{0xffffffff, []byte{0xff,0xff,0xff,0xff}}, // Max 4-byte
	}

	t.Logf("Running uint32B %d tests", len(tests))
	for i, test := range tests {
		var buf bytes.Buffer
		err := WriteUint32B(&buf, test.in)
		if err != nil {
			t.Errorf("WriteUint32B #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteUint32B #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Read from protos format.
		rbuf := bytes.NewReader(test.buf)
		var val uint32
		err = ReadUint32B(rbuf, &val)
		if err != nil {
			t.Errorf("ReadUint32B #%d error %v", i, err)
			continue
		}
		if val != test.in {
			t.Errorf("ReadUint32B #%d\n got: %v want: %v", i,
				val, test.in)
			continue
		}
	}
}

// TestUint32Errors performs negative tests against protos encode and decode
// of various element types to confirm error paths work correctly.
func TestUint32BErrors(t *testing.T) {
	tests := []struct {
		in       uint32 // Value to encode
		max      int         // Max size of fixed buffer to induce errors
		writeErr error       // Expected write error
		readErr  error       // Expected read error
	}{
		{1, 0,io.ErrShortWrite, io.EOF}, // Force errors on 1-byte read/write.
		{256, 0,io.ErrShortWrite, io.EOF},  // Force errors on 2-byte read/write.
		{0x10000, 2, io.ErrShortWrite, io.ErrUnexpectedEOF}, // Force errors on 4-byte read/write.
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := WriteUint32B(w, test.in)
		if err != test.writeErr {
			t.Errorf("WriteUint32B #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		r := newFixedReader(test.max, nil)
		var val uint32
		err = ReadUint32B(r, &val)
		if err != test.readErr {
			t.Errorf("ReadUint32B #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}


// TestUint64B tests protos encode and decode for byte type various elements.
func TestUint64B(t *testing.T) {
	tests := []struct {
		in   uint64   // Value to encode
		buf  []byte // serialized
	}{
		{1, []byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x01}},     // Min single byte
		{255, []byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00,0xff}},   // Max single byte
		{256, []byte{0x00,0x00,0x00,0x00,0x00,0x00,0x01,0x00}},   // Min 2-byte
		{65535, []byte{0x00,0x00,0x00,0x00,0x00,0x00,0xff,0xff}}, // Max 2-byte
		{0x10000, []byte{0x00,0x00,0x00,0x00,0x00,0x01,0x00,0x00}},   // Min 4-byte
		{0xffffffff, []byte{0x00,0x00,0x00,0x00,0xff,0xff,0xff,0xff}}, // Max 4-byte
		{0x100000000, []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}},   // Min 8-byte
		{0xffffffffffffffff, []byte{0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff}}, // Max 8-byte
	}

	t.Logf("Running uint64B %d tests", len(tests))
	for i, test := range tests {
		var buf bytes.Buffer
		err := WriteUint64B(&buf, test.in)
		if err != nil {
			t.Errorf("WriteUint64B #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteUint64B #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Read from protos format.
		rbuf := bytes.NewReader(test.buf)
		var val uint64
		err = ReadUint64B(rbuf, &val)
		if err != nil {
			t.Errorf("ReadUint64B #%d error %v", i, err)
			continue
		}
		if val != test.in {
			t.Errorf("ReadUint64B #%d\n got: %v want: %v", i,
				val, test.in)
			continue
		}
	}
}

// TestUint64BErrors performs negative tests against protos encode and decode
// of various element types to confirm error paths work correctly.
func TestUint64BErrors(t *testing.T) {
	tests := []struct {
		in       uint64 // Value to encode
		max      int         // Max size of fixed buffer to induce errors
		writeErr error       // Expected write error
		readErr  error       // Expected read error
	}{
		{1, 0,io.ErrShortWrite, io.EOF}, // Force errors on 1-byte read/write.
		{256, 0,io.ErrShortWrite, io.EOF},  // Force errors on 2-byte read/write.
		{0x10000, 2, io.ErrShortWrite, io.ErrUnexpectedEOF}, // Force errors on 4-byte read/write.
		{0x100000000, 4, io.ErrShortWrite, io.ErrUnexpectedEOF}, // Force errors on 8-byte read/write.
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := WriteUint64B(w, test.in)
		if err != test.writeErr {
			t.Errorf("WriteUint64B #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		r := newFixedReader(test.max, nil)
		var val uint64
		err = ReadUint64B(r, &val)
		if err != test.readErr {
			t.Errorf("ReadUint64B #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

func TestWriteNBytes(t *testing.T) {
	bytes256 := bytes.Repeat([]byte{0x01}, 256)
	tests := []struct {
		in   []byte // Byte Array to write
		length int
		buf  []byte // Wire encoding
	}{
		// Latest protocol version.
		// Empty byte array
		{[]byte{}, 0,[]byte{}},
		// 1-byte:
		{[]byte{0x01}, 1,[]byte{0x01}},
		// 256-byte
		{bytes256, 256,bytes256},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		var buf bytes.Buffer
		err := WriteNBytes(&buf, test.in)
		if err != nil {
			t.Errorf("WriteVarBytes #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteVarBytes #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Decode from protos format.
		rbuf := bytes.NewReader(test.buf)
		val := make([]byte,test.length)
		err = ReadNBytes(rbuf, val[:], test.length)
		if err != nil {
			t.Errorf("ReadVarBytes #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("ReadVarBytes #%d\n got: %s want: %s", i,
				val, test.buf)
			continue
		}
	}
}

func TestWriteNBytesErrors(t *testing.T) {
	bytes256 := bytes.Repeat([]byte{0x01}, 256)
	tests := []struct {
		in       []byte // byte to write
		buf      []byte // result byte
		max      int    // Max size of fixed buffer to induce errors
		writeErr error  // Expected write error
		readErr  error  // Expected read error
	}{
		{[]byte{0x01}, []byte{0x04}, 0, io.ErrShortWrite, io.EOF},
		{bytes256, bytes256, 2, io.ErrShortWrite, io.ErrUnexpectedEOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := WriteNBytes(w, test.in)
		if err != test.writeErr {
			t.Errorf("WriteNBytes #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}
		r := newFixedReader(test.max, test.buf)
		val := make([]byte,len(test.buf))
		err = ReadNBytes(r, val[:], len(test.buf))
		if err != test.readErr {
			t.Errorf("ReadVarString #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

// TestVarIntWire tests protos encode and decode for variable length integers.
func TestVarIntWire(t *testing.T) {
	pver := ProtocolVersion

	tests := []struct {
		in   uint64 // Value to encode
		out  uint64 // Expected decoded value
		buf  []byte // Wire encoding
		pver uint32 // Protocol version for protos encoding
	}{
		// Latest protocol version.
		// Single byte
		{0, 0, []byte{0x00}, pver},
		// Max single byte
		{0xfc, 0xfc, []byte{0xfc}, pver},
		// Min 2-byte
		{0xfd, 0xfd, []byte{0xfd, 0x0fd, 0x00}, pver},
		// Max 2-byte
		{0xffff, 0xffff, []byte{0xfd, 0xff, 0xff}, pver},
		// Min 4-byte
		{0x10000, 0x10000, []byte{0xfe, 0x00, 0x00, 0x01, 0x00}, pver},
		// Max 4-byte
		{0xffffffff, 0xffffffff, []byte{0xfe, 0xff, 0xff, 0xff, 0xff}, pver},
		// Min 8-byte
		{
			0x100000000, 0x100000000,
			[]byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00},
			pver,
		},
		// Max 8-byte
		{
			0xffffffffffffffff, 0xffffffffffffffff,
			[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			pver,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		var buf bytes.Buffer
		err := WriteVarInt(&buf, test.pver, test.in)
		if err != nil {
			t.Errorf("WriteVarInt #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteVarInt #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Decode from protos format.
		rbuf := bytes.NewReader(test.buf)
		val, err := ReadVarInt(rbuf, test.pver)
		if err != nil {
			t.Errorf("ReadVarInt #%d error %v", i, err)
			continue
		}
		if val != test.out {
			t.Errorf("ReadVarInt #%d\n got: %d want: %d", i,
				val, test.out)
			continue
		}
	}
}

// TestVarIntWireErrors performs negative tests against protos encode and decode
// of variable length integers to confirm error paths work correctly.
func TestVarIntWireErrors(t *testing.T) {
	pver := ProtocolVersion

	tests := []struct {
		in       uint64 // Value to encode
		buf      []byte // Wire encoding
		pver     uint32 // Protocol version for protos encoding
		max      int    // Max size of fixed buffer to induce errors
		writeErr error  // Expected write error
		readErr  error  // Expected read error
	}{
		// Force errors on discriminant.
		{0, []byte{0x00}, pver, 0, io.ErrShortWrite, io.EOF},
		// Force errors on 2-byte read/write.
		{0xfd, []byte{0xfd}, pver, 2, io.ErrShortWrite, io.ErrUnexpectedEOF},
		// Force errors on 4-byte read/write.
		{0x10000, []byte{0xfe}, pver, 2, io.ErrShortWrite, io.ErrUnexpectedEOF},
		// Force errors on 8-byte read/write.
		{0x100000000, []byte{0xff}, pver, 2, io.ErrShortWrite, io.ErrUnexpectedEOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := WriteVarInt(w, test.pver, test.in)
		if err != test.writeErr {
			t.Errorf("WriteVarInt #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		r := newFixedReader(test.max, test.buf)
		_, err = ReadVarInt(r, test.pver)
		if err != test.readErr {
			t.Errorf("ReadVarInt #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

// TestVarIntNonCanonical ensures variable length integers that are not encoded
// canonically return the expected error.
func TestVarIntNonCanonical(t *testing.T) {
	pver := ProtocolVersion

	tests := []struct {
		name string // Test name for easier identification
		in   []byte // Value to decode
		pver uint32 // Protocol version for protos encoding
	}{
		{
			"0 encoded with 3 bytes", []byte{0xfd, 0x00, 0x00},
			pver,
		},
		{
			"max single-byte value encoded with 3 bytes",
			[]byte{0xfd, 0xfc, 0x00}, pver,
		},
		{
			"0 encoded with 5 bytes",
			[]byte{0xfe, 0x00, 0x00, 0x00, 0x00}, pver,
		},
		{
			"max three-byte value encoded with 5 bytes",
			[]byte{0xfe, 0xff, 0xff, 0x00, 0x00}, pver,
		},
		{
			"0 encoded with 9 bytes",
			[]byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			pver,
		},
		{
			"max five-byte value encoded with 9 bytes",
			[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00},
			pver,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Decode from protos format.
		rbuf := bytes.NewReader(test.in)
		val, err := ReadVarInt(rbuf, test.pver)
		if _, ok := err.(*MessageError); !ok {
			t.Errorf("ReadVarInt #%d (%s) unexpected error %v", i,
				test.name, err)
			continue
		}
		if val != 0 {
			t.Errorf("ReadVarInt #%d (%s)\n got: %d want: 0", i,
				test.name, val)
			continue
		}
	}
}

// TestVarIntWire tests the serialize size for variable length integers.
func TestVarIntSerializeSize(t *testing.T) {
	tests := []struct {
		val  uint64 // Value to get the serialized size for
		size int    // Expected serialized size
	}{
		// Single byte
		{0, 1},
		// Max single byte
		{0xfc, 1},
		// Min 2-byte
		{0xfd, 3},
		// Max 2-byte
		{0xffff, 3},
		// Min 4-byte
		{0x10000, 5},
		// Max 4-byte
		{0xffffffff, 5},
		// Min 8-byte
		{0x100000000, 9},
		// Max 8-byte
		{0xffffffffffffffff, 9},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		serializedSize := VarIntSerializeSize(test.val)
		if serializedSize != test.size {
			t.Errorf("VarIntSerializeSize #%d got: %d, want: %d", i,
				serializedSize, test.size)
			continue
		}
	}
}

// TestVarStringWire tests protos encode and decode for variable length strings.
func TestVarStringWire(t *testing.T) {
	pver := ProtocolVersion

	// str256 is a string that takes a 2-byte varint to encode.
	str256 := strings.Repeat("test", 64)

	tests := []struct {
		in   string // String to encode
		out  string // String to decoded value
		buf  []byte // Wire encoding
		pver uint32 // Protocol version for protos encoding
	}{
		// Latest protocol version.
		// Empty string
		{"", "", []byte{0x00}, pver},
		// Single byte varint + string
		{"Test", "Test", append([]byte{0x04}, []byte("Test")...), pver},
		// 2-byte varint + string
		{str256, str256, append([]byte{0xfd, 0x00, 0x01}, []byte(str256)...), pver},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		var buf bytes.Buffer
		err := WriteVarString(&buf, test.pver, test.in)
		if err != nil {
			t.Errorf("WriteVarString #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteVarString #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Decode from protos format.
		rbuf := bytes.NewReader(test.buf)
		val, err := ReadVarString(rbuf, test.pver)
		if err != nil {
			t.Errorf("ReadVarString #%d error %v", i, err)
			continue
		}
		if val != test.out {
			t.Errorf("ReadVarString #%d\n got: %s want: %s", i,
				val, test.out)
			continue
		}
	}
}

// TestVarStringWireErrors performs negative tests against protos encode and
// decode of variable length strings to confirm error paths work correctly.
func TestVarStringWireErrors(t *testing.T) {
	pver := ProtocolVersion

	// str256 is a string that takes a 2-byte varint to encode.
	str256 := strings.Repeat("test", 64)

	tests := []struct {
		in       string // Value to encode
		buf      []byte // Wire encoding
		pver     uint32 // Protocol version for protos encoding
		max      int    // Max size of fixed buffer to induce errors
		writeErr error  // Expected write error
		readErr  error  // Expected read error
	}{
		// Latest protocol version with intentional read/write errors.
		// Force errors on empty string.
		{"", []byte{0x00}, pver, 0, io.ErrShortWrite, io.EOF},
		// Force error on single byte varint + string.
		{"Test", []byte{0x04}, pver, 2, io.ErrShortWrite, io.ErrUnexpectedEOF},
		// Force errors on 2-byte varint + string.
		{str256, []byte{0xfd}, pver, 2, io.ErrShortWrite, io.ErrUnexpectedEOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := WriteVarString(w, test.pver, test.in)
		if err != test.writeErr {
			t.Errorf("WriteVarString #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		r := newFixedReader(test.max, test.buf)
		_, err = ReadVarString(r, test.pver)
		if err != test.readErr {
			t.Errorf("ReadVarString #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

// TestVarStringOverflowErrors performs tests to ensure deserializing variable
// length strings intentionally crafted to use large values for the string
// length are handled properly.  This could otherwise potentially be used as an
// attack vector.
func TestVarStringOverflowErrors(t *testing.T) {
	pver := ProtocolVersion

	tests := []struct {
		buf  []byte // Wire encoding
		pver uint32 // Protocol version for protos encoding
		err  error  // Expected error
	}{
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			pver, &MessageError{}},
		{[]byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			pver, &MessageError{}},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Decode from protos format.
		rbuf := bytes.NewReader(test.buf)
		_, err := ReadVarString(rbuf, test.pver)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("ReadVarString #%d wrong error got: %v, "+
				"want: %v", i, err, reflect.TypeOf(test.err))
			continue
		}
	}

}

// TestVarBytesWire tests protos encode and decode for variable length byte array.
func TestVarBytesWire(t *testing.T) {
	pver := ProtocolVersion

	// bytes256 is a byte array that takes a 2-byte varint to encode.
	bytes256 := bytes.Repeat([]byte{0x01}, 256)

	tests := []struct {
		in   []byte // Byte Array to write
		buf  []byte // Wire encoding
		pver uint32 // Protocol version for protos encoding
	}{
		// Latest protocol version.
		// Empty byte array
		{[]byte{}, []byte{0x00}, pver},
		// Single byte varint + byte array
		{[]byte{0x01}, []byte{0x01, 0x01}, pver},
		// 2-byte varint + byte array
		{bytes256, append([]byte{0xfd, 0x00, 0x01}, bytes256...), pver},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		var buf bytes.Buffer
		err := WriteVarBytes(&buf, test.pver, test.in)
		if err != nil {
			t.Errorf("WriteVarBytes #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteVarBytes #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Decode from protos format.
		rbuf := bytes.NewReader(test.buf)
		val, err := ReadVarBytes(rbuf, test.pver, MaxMessagePayload,
			"test payload")
		if err != nil {
			t.Errorf("ReadVarBytes #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("ReadVarBytes #%d\n got: %s want: %s", i,
				val, test.buf)
			continue
		}
	}
}

// TestVarBytesWireErrors performs negative tests against protos encode and
// decode of variable length byte arrays to confirm error paths work correctly.
func TestVarBytesWireErrors(t *testing.T) {
	pver := uint32(0)

	// bytes256 is a byte array that takes a 2-byte varint to encode.
	bytes256 := bytes.Repeat([]byte{0x01}, 256)

	tests := []struct {
		in       []byte // Byte Array to write
		buf      []byte // Wire encoding
		pver     uint32 // Protocol version for protos encoding
		max      int    // Max size of fixed buffer to induce errors
		writeErr error  // Expected write error
		readErr  error  // Expected read error
	}{
		// Latest protocol version with intentional read/write errors.
		// Force errors on empty byte array.
		{[]byte{}, []byte{0x00}, pver, 0, io.ErrShortWrite, io.EOF},
		// Force error on single byte varint + byte array.
		{[]byte{0x01, 0x02, 0x03}, []byte{0x04}, pver, 2, io.ErrShortWrite, io.ErrUnexpectedEOF},
		// Force errors on 2-byte varint + byte array.
		{bytes256, []byte{0xfd}, pver, 2, io.ErrShortWrite, io.ErrUnexpectedEOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := WriteVarBytes(w, test.pver, test.in)
		if err != test.writeErr {
			t.Errorf("WriteVarBytes #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from protos format.
		r := newFixedReader(test.max, test.buf)
		_, err = ReadVarBytes(r, test.pver, MaxVarStringPayload,
			"test payload")
		if err != test.readErr {
			t.Errorf("ReadVarBytes #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

// TestVarBytesOverflowErrors performs tests to ensure deserializing variable
// length byte arrays intentionally crafted to use large values for the array
// length are handled properly.  This could otherwise potentially be used as an
// attack vector.
func TestVarBytesOverflowErrors(t *testing.T) {
	pver := uint32(0)

	tests := []struct {
		buf  []byte // Wire encoding
		pver uint32 // Protocol version for protos encoding
		err  error  // Expected error
	}{
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			pver, &MessageError{}},
		{[]byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			pver, &MessageError{}},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Decode from protos format.
		rbuf := bytes.NewReader(test.buf)
		_, err := ReadVarBytes(rbuf, test.pver, MaxVarStringPayload,
			"test payload")
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("ReadVarBytes #%d wrong error got: %v, "+
				"want: %v", i, err, reflect.TypeOf(test.err))
			continue
		}
	}

}

// TestRandomUint64 exercises the randomness of the random number generator on
// the system by ensuring the probability of the generated numbers.  If the RNG
// is evenly distributed as a proper cryptographic RNG should be, there really
// should only be 1 number < 2^56 in 2^8 tries for a 64-bit number.  However,
// use a higher number of 5 to really ensure the test doesn't fail unless the
// RNG is just horrendous.
func TestRandomUint64(t *testing.T) {
	tries := 1 << 8              // 2^8
	watermark := uint64(1 << 56) // 2^56
	maxHits := 5
	badRNG := "The random number generator on this system is clearly " +
		"terrible since we got %d values less than %d in %d runs " +
		"when only %d was expected"

	numHits := 0
	for i := 0; i < tries; i++ {
		nonce, err := RandomUint64()
		if err != nil {
			t.Errorf("RandomUint64 iteration %d failed - err %v",
				i, err)
			return
		}
		if nonce < watermark {
			numHits++
		}
		if numHits > maxHits {
			str := fmt.Sprintf(badRNG, numHits, watermark, tries, maxHits)
			t.Errorf("Random Uint64 iteration %d failed - %v %v", i,
				str, numHits)
			return
		}
	}
}

// TestRandomUint64Errors uses a fake reader to force error paths to be executed
// and checks the results accordingly.
func TestRandomUint64Errors(t *testing.T) {
	// Test short reads.
	fr := &fakeRandReader{n: 2, err: io.EOF}
	nonce, err := randomUint64(fr)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("Error not expected value of %v [%v]",
			io.ErrUnexpectedEOF, err)
	}
	if nonce != 0 {
		t.Errorf("Nonce is not 0 [%v]", nonce)
	}
}
