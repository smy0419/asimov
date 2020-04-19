// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package ffldb

import (
	"encoding/binary"
	"github.com/AsimovNetwork/asimov/protos"
)

const (
	// metadataDbName is the name used for the metadata database.
	metadataDbName = "metadata"

	// blockHdrSize is the size of a block header.  This is simply the
	// constant from protos and is only provided here for convenience since
	// protos.BlockHeaderPayload is quite long.
	blockHdrSize = protos.BlockHeaderPayload

	// errDbNotOpenStr is the text to use for the database.ErrDbNotOpen
	// error code.
	errDbNotOpenStr = "database is not open"

	// errTxClosedStr is the text to use for the database.ErrTxClosed error
	// code.
	errTxClosedStr = "database tx is closed"
)

var (
	// byteOrder is the preferred byte order used through the database and
	// block files.  Sometimes big endian will be used to allow ordered byte
	// sortable integer values.
	byteOrder = binary.LittleEndian

	// bucketIndexPrefix is the prefix used for all entries in the bucket
	// index.
	bucketIndexPrefix = []byte("bidx")

	// curBucketIDKeyName is the name of the key used to keep track of the
	// current bucket ID counter.
	curBucketIDKeyName = []byte("bidx-cbid")

	// metadataBucketID is the ID of the top-level metadata bucket.
	// It is the value 0 encoded as an unsigned big-endian uint32.
	metadataBucketID = [4]byte{}

	// blockIdxBucketID is the ID of the internal block metadata bucket.
	// It is the value 1 encoded as an unsigned big-endian uint32.
	blockIdxBucketID = [4]byte{0x00, 0x00, 0x00, 0x01}

	// blockIdxBucketName is the bucket used internally to track block
	// metadata.
	blockIdxBucketName = []byte("ffldb-blockidx")

	// writeLocKeyName is the key used to store the current write file
	// location.
	writeLocKeyName = []byte("ffldb-writeloc")
)

