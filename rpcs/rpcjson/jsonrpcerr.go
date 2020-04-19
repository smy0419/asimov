// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.


package rpcjson

// Standard JSON-RPC 2.0 errors.
var (
	ErrRPCInternal = &RPCError{
		Code:    -32603,
		Message: "Internal error",
	}
)

// General application defined JSON errors.
const (
	ErrRPCMisc                RPCErrorCode = -11
	ErrRPCInvalidOutPutAmount RPCErrorCode = -12
	ErrRPCInvalidAddressOrKey RPCErrorCode = -13
	ErrRPCInvalidParameter    RPCErrorCode = -14
	ErrRPCDeserialization     RPCErrorCode = -15
)

// Specific Errors related to commands.  These are the ones a user of the RPC
// server are most likely to see.  Generally, the codes should match one of the
// more general errors above.
const (
	ErrRPCBlockNotFound       RPCErrorCode = -200
	ErrRPCBlockHeaderNotFound RPCErrorCode = -201
	ErrRPCBlockHeightNotFound RPCErrorCode = -202
	ErrRPCBlockHashNotFound   RPCErrorCode = -203
	ErrRPCNoTxInfo            RPCErrorCode = -204
	ErrRPCInvalidTxVout       RPCErrorCode = -205
	ErrRPCDecodeHexString     RPCErrorCode = -206
)

