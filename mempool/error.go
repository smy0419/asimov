// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"github.com/AsimovNetwork/asimov/blockchain"
	"github.com/AsimovNetwork/asimov/protos"
)

// RuleError identifies a rule violation.  It is used to indicate that
// processing of a transaction failed due to one of the many validation
// rules.  The caller can use type assertions to determine if a failure was
// specifically due to a rule violation and use the Err field to access the
// underlying error, which will be either a TxRuleError or a
// blockchain.RuleError.
type RuleError struct {
	Err error
}

// Error satisfies the error interface and prints human-readable errors.
func (e RuleError) Error() string {
	if e.Err == nil {
		return "<nil>"
	}
	return e.Err.Error()
}

// TxRuleError identifies a rule violation.  It is used to indicate that
// processing of a transaction failed due to one of the many validation
// rules.  The caller can use type assertions to determine if a failure was
// specifically due to a rule violation and access the ErrorCode field to
// ascertain the specific reason for the rule violation.
type TxRuleError struct {
	RejectCode  protos.RejectCode // The code to send with reject messages
	Description string            // Human readable description of the issue
}

// Error satisfies the error interface and prints human-readable errors.
func (e TxRuleError) Error() string {
	return e.Description
}

// txRuleError creates an underlying TxRuleError with the given a set of
// arguments and returns a RuleError that encapsulates it.
func txRuleError(c protos.RejectCode, desc string) RuleError {
	return RuleError{
		Err: TxRuleError{RejectCode: c, Description: desc},
	}
}

// chainRuleError returns a RuleError that encapsulates the given
// blockchain.RuleError.
func chainRuleError(chainErr blockchain.RuleError) RuleError {
	return RuleError{
		Err: chainErr,
	}
}

// extractRejectCode attempts to return a relevant reject code for a given error
// by examining the error for known types.  It will return true if a code
// was successfully extracted.
func extractRejectCode(err error) (protos.RejectCode, bool) {
	// Pull the underlying error out of a RuleError.
	if rerr, ok := err.(RuleError); ok {
		err = rerr.Err
	}

	switch err := err.(type) {
	case blockchain.RuleError:
		// Convert the chain error to a reject code.
		var code protos.RejectCode
		switch err.ErrorCode {
		// Rejected due to duplicate.
		case blockchain.ErrDuplicateBlock:
			code = protos.RejectDuplicate

		// Rejected due to obsolete version.
		case blockchain.ErrBlockVersionTooOld:
			code = protos.RejectObsolete

		// Rejected due to checkpoint.
		case blockchain.ErrCheckpointTimeTooOld:
			fallthrough
		case blockchain.ErrBadCheckpoint:
			fallthrough
		case blockchain.ErrForkTooOld:
			code = protos.RejectCheckpoint

		// Everything else is due to the block or transaction being invalid.
		default:
			code = protos.RejectInvalid
		}

		return code, true

	case TxRuleError:
		return err.RejectCode, true

	case nil:
		return protos.RejectInvalid, false
	}

	return protos.RejectInvalid, false
}

// ErrToRejectErr examines the underlying type of the error and returns a reject
// code and string appropriate to be sent in a protos.MsgReject message.
func ErrToRejectErr(err error) (protos.RejectCode, string) {
	// Return the reject code along with the error text if it can be
	// extracted from the error.
	rejectCode, found := extractRejectCode(err)
	if found {
		return rejectCode, err.Error()
	}

	// Return a generic rejected string if there is no error.  This really
	// should not happen unless the code elsewhere is not setting an error
	// as it should be, but it's best to be safe and simply return a generic
	// string rather than allowing the following code that dereferences the
	// err to panic.
	if err == nil {
		return protos.RejectInvalid, "rejected"
	}

	// When the underlying error is not one of the above cases, just return
	// protos.RejectInvalid with a generic rejected string plus the error
	// text.
	return protos.RejectInvalid, "rejected: " + err.Error()
}
