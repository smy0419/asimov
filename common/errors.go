// Copyright (c) 2018-2020. The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package common

import "errors"

var (
	// ErrUnknownAddressType describes an error where an address can not
	// decoded as a specific address type due to the string encoding
	// begining with an identifier byte unknown to any standard or
	// registered (via chaincfg.Register) network.
	ErrUnknownAddressType = errors.New("unknown address type")
)


// AssertError identifies an error that indicates an internal code consistency
// issue and should be treated as a critical and unrecoverable error.
type AssertError string

// Error returns the assertion error as a human-readable string and satisfies
// the error interface.
func (e AssertError) Error() string {
	return "assertion failed: " + string(e)
}


// DeserializeError signifies that a problem was encountered when deserializing
// data.
type DeserializeError string

// Error implements the error interface.
func (e DeserializeError) Error() string {
	return string(e)
}

// IsDeserializeErr returns whether or not the passed error is an DeserializeError
// error.
func IsDeserializeErr(err error) bool {
	_, ok := err.(DeserializeError)
	return ok
}
