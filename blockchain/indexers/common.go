// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package indexers

import (
	"encoding/binary"
	"errors"
)

var (
	// byteOrder is the preferred byte order used for serializing numeric
	// fields for storage in the database.
	byteOrder = binary.LittleEndian

	// errInterruptRequested indicates that an operation was cancelled due
	// to a user-requested interrupt.
	errInterruptRequested = errors.New("interrupt requested")
)

// NeedsInputser provides a generic interface for an indexer to specify the it
// requires the ability to look up inputs for a transaction.
type NeedsInputser interface {
	NeedsInputs() bool
}

// internalBucket is an abstraction over a database bucket.  It is used to make
// the code easier to test since it allows mock objects in the tests to only
// implement these functions instead of everything a database.Bucket supports.
type internalBucket interface {
	Get(key []byte) []byte
	Put(key []byte, value []byte) error
	Delete(key []byte) error
}

// interruptRequested returns true when the provided channel has been closed.
// This simplifies early shutdown slightly since the caller can just use an if
// statement instead of a select.
func interruptRequested(interrupted <-chan struct{}) bool {
	select {
	case <-interrupted:
		return true
	default:
	}

	return false
}
