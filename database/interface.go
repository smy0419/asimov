// Copyright (c) 2018-2020. The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Copyright 2014 The go-ethereum Authors
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package database

import (
	"encoding/hex"
	"github.com/AsimovNetwork/asimov/common"
)

// Code using batches should try to add this much data to the batch.
// The value was determined empirically.
const IdealBatchSize = 100 * 1024

type BlockType byte

const (
	FFLDB = "ffldb"
	ETHDB = "ethdb"
	// blockKeyLength is the size of block key
	BlockKeyLength = 33

	BlockNormal BlockType = 0
	BlockVirtual BlockType = 1
)

type BlockKey [BlockKeyLength]byte

func (blockkey* BlockKey) String() string {
	return hex.EncodeToString(blockkey[:])
}

func (blockkey* BlockKey) IsVirtual() bool {
	return blockkey[common.HashLength] == byte(BlockVirtual)
}

func NewBlockKey(hash *common.Hash, blockType BlockType) *BlockKey {
	var key BlockKey
	copy(key[:common.HashLength], hash[:])
	key[common.HashLength] = byte(blockType)
	return &key
}

func NewNormalBlockKey(hash *common.Hash) *BlockKey {
	return NewBlockKey(hash, BlockNormal)
}

func NewVirtualBlockKey(hash *common.Hash) *BlockKey {
	return NewBlockKey(hash, BlockVirtual)
}

// Putter wraps the database write operation supported by both batches and regular databases.
type Putter interface {
	Put(key []byte, value []byte) error
}

// Deleter wraps the database delete operation supported by both batches and regular databases.
type Deleter interface {
	Delete(key []byte) error
}

type Queryer interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
}

// Batch is a write-only database that commits changes to its host database
// when Write is called. Batch cannot be used concurrently.
type Batch interface {
	Putter
	Deleter
	ValueSize() int // amount of data in the batch
	Write() error
	// Reset resets the batch for reuse
	Reset()
}


type Transactor interface {
	Begin(writable bool) (Tx, error)
	View(fn func(tx Tx) error) error
	Update(fn func(tx Tx) error) error
	Close() error
}

// Database wraps all database operations. All methods are safe for concurrent use.
type Database interface {
	Type() string

	Putter
	Deleter
	Queryer

	Transactor

	NewBatch() Batch
}

// BlockRegion specifies a particular region of a block identified by the
// specified key, given an offset and length.
type BlockRegion struct {
	Key   *BlockKey
	Offset uint32
	Len    uint32
}

// Tx represents a database transaction.  It can either by read-only or
// read-write.  The transaction provides a metadata bucket against which all
// read and writes occur.
//
// As would be expected with a transaction, no changes will be saved to the
// database until it has been committed.  The transaction will only provide a
// view of the database at the time it was created.  Transactions should not be
// long running operations.
type Tx interface {
	// StoreBlock stores the provided block into the database.  There are no
	// checks to ensure the block connects to a previous block, contains
	// double spends, or any additional functionality such as transaction
	// indexing.  It simply stores the block in the database.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBlockExists when the block key already exists
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	//
	// Other errors are possible depending on the implementation.
	StoreBlock(blockKey *BlockKey, blockBytes []byte) error

	// HasBlock returns whether or not a block with the given key exists
	// in the database.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrTxClosed if the transaction has already been closed
	//
	// Other errors are possible depending on the implementation.
	HasBlock(key *BlockKey) (bool, error)

	// HasBlocks returns whether or not the blocks with the provided keys
	// exist in the database.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrTxClosed if the transaction has already been closed
	//
	// Other errors are possible depending on the implementation.
	HasBlocks(keys []BlockKey) ([]bool, error)

	// FetchBlockHeader returns the raw serialized bytes for the block
	// header identified by the given key.  The raw bytes are in the format
	// returned by Serialize on a protos.BlockHeader.
	//
	// It is highly recommended to use this function (or FetchBlockHeaders)
	// to obtain block headers over the FetchBlockRegion(s) functions since
	// it provides the backend drivers the freedom to perform very specific
	// optimizations which can result in significant speed advantages when
	// working with headers.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBlockNotFound if the requested block key does not exist
	//   - ErrTxClosed if the transaction has already been closed
	//   - ErrCorruption if the database has somehow become corrupted
	//
	// NOTE: The data returned by this function is only valid during a
	// database transaction.  Attempting to access it after a transaction
	// has ended results in undefined behavior.  This constraint prevents
	// additional data copies and allows support for memory-mapped database
	// implementations.
	FetchBlockHeader(key *BlockKey) ([]byte, error)

	// FetchBlockHeaders returns the raw serialized bytes for the block
	// headers identified by the given keys.  The raw bytes are in the
	// format returned by Serialize on a protos.BlockHeader.
	//
	// It is highly recommended to use this function (or FetchBlockHeader)
	// to obtain block headers over the FetchBlockRegion(s) functions since
	// it provides the backend drivers the freedom to perform very specific
	// optimizations which can result in significant speed advantages when
	// working with headers.
	//
	// Furthermore, depending on the specific implementation, this function
	// can be more efficient for bulk loading multiple block headers than
	// loading them one-by-one with FetchBlockHeader.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBlockNotFound if any of the request block keys do not exist
	//   - ErrTxClosed if the transaction has already been closed
	//   - ErrCorruption if the database has somehow become corrupted
	//
	// NOTE: The data returned by this function is only valid during a
	// database transaction.  Attempting to access it after a transaction
	// has ended results in undefined behavior.  This constraint prevents
	// additional data copies and allows support for memory-mapped database
	// implementations.
	FetchBlockHeaders(keys []BlockKey) ([][]byte, error)

	// FetchBlock returns the raw serialized bytes for the block identified
	// by the given key.  The raw bytes are in the format returned by
	// Serialize on a protos.MsgBlock.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBlockNotFound if the requested block key does not exist
	//   - ErrTxClosed if the transaction has already been closed
	//   - ErrCorruption if the database has somehow become corrupted
	//
	// NOTE: The data returned by this function is only valid during a
	// database transaction.  Attempting to access it after a transaction
	// has ended results in undefined behavior.  This constraint prevents
	// additional data copies and allows support for memory-mapped database
	// implementations.
	FetchBlock(key *BlockKey) ([]byte, error)

	// FetchBlocks returns the raw serialized bytes for the blocks
	// identified by the given keys.  The raw bytes are in the format
	// returned by Serialize on a protos.MsgBlock.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBlockNotFound if the any of the requested block keys do not
	//     exist
	//   - ErrTxClosed if the transaction has already been closed
	//   - ErrCorruption if the database has somehow become corrupted
	//
	// NOTE: The data returned by this function is only valid during a
	// database transaction.  Attempting to access it after a transaction
	// has ended results in undefined behavior.  This constraint prevents
	// additional data copies and allows support for memory-mapped database
	// implementations.
	FetchBlocks(keys []BlockKey) ([][]byte, error)

	// FetchBlockRegion returns the raw serialized bytes for the given
	// block region.
	//
	// For example, it is possible to directly extract Bitcoin transactions
	// and/or scripts from a block with this function.  Depending on the
	// backend implementation, this can provide significant savings by
	// avoiding the need to load entire blocks.
	//
	// The raw bytes are in the format returned by Serialize on a
	// protos.MsgBlock and the Offset field in the provided BlockRegion is
	// zero-based and relative to the start of the block (byte 0).
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBlockNotFound if the requested block key does not exist
	//   - ErrBlockRegionInvalid if the region exceeds the bounds of the
	//     associated block
	//   - ErrTxClosed if the transaction has already been closed
	//   - ErrCorruption if the database has somehow become corrupted
	//
	// NOTE: The data returned by this function is only valid during a
	// database transaction.  Attempting to access it after a transaction
	// has ended results in undefined behavior.  This constraint prevents
	// additional data copies and allows support for memory-mapped database
	// implementations.
	FetchBlockRegion(region *BlockRegion) ([]byte, error)

	// FetchBlockRegions returns the raw serialized bytes for the given
	// block regions.
	//
	// For example, it is possible to directly extract Bitcoin transactions
	// and/or scripts from various blocks with this function.  Depending on
	// the backend implementation, this can provide significant savings by
	// avoiding the need to load entire blocks.
	//
	// The raw bytes are in the format returned by Serialize on a
	// protos.MsgBlock and the Offset fields in the provided BlockRegions are
	// zero-based and relative to the start of the block (byte 0).
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBlockNotFound if any of the requested block do not
	//     exist
	//   - ErrBlockRegionInvalid if one or more region exceed the bounds of
	//     the associated block
	//   - ErrTxClosed if the transaction has already been closed
	//   - ErrCorruption if the database has somehow become corrupted
	//
	// NOTE: The data returned by this function is only valid during a
	// database transaction.  Attempting to access it after a transaction
	// has ended results in undefined behavior.  This constraint prevents
	// additional data copies and allows support for memory-mapped database
	// implementations.
	FetchBlockRegions(regions []BlockRegion) ([][]byte, error)

	// Metadata returns the top-most bucket for all metadata storage.
	Metadata() Bucket

	//Start() error

	// ******************************************************************
	// Methods related to both atomic metadata storage and block storage.
	// ******************************************************************

	// Commit commits all changes that have been made to the metadata or
	// block storage.  Depending on the backend implementation this could be
	// to a cache that is periodically synced to persistent storage or
	// directly to persistent storage.  In any case, all transactions which
	// are started after the commit finishes will include all changes made
	// by this transaction.  Calling this function on a managed transaction
	// will result in a panic.
	Commit() error

	// Rollback undoes all changes that have been made to the metadata or
	// block storage.  Calling this function on a managed transaction will
	// result in a panic.
	Rollback() error
}


// Bucket represents a collection of key/value pairs.
type Bucket interface {
	// Bucket retrieves a nested bucket with the given key.  Returns nil if
	// the bucket does not exist.
	Bucket(key []byte) Bucket

	// CreateBucket creates and returns a new nested bucket with the given
	// key.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBucketExists if the bucket already exists
	//   - ErrBucketNameRequired if the key is empty
	//   - ErrIncompatibleValue if the key is otherwise invalid for the
	//     particular implementation
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	CreateBucket(key []byte) (Bucket, error)

	// CreateBucketIfNotExists creates and returns a new nested bucket with
	// the given key if it does not already exist.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBucketNameRequired if the key is empty
	//   - ErrIncompatibleValue if the key is otherwise invalid for the
	//     particular implementation
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	CreateBucketIfNotExists(key []byte) (Bucket, error)

	// DeleteBucket removes a nested bucket with the given key.  This also
	// includes removing all nested buckets and keys under the bucket being
	// deleted.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBucketNotFound if the specified bucket does not exist
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	DeleteBucket(key []byte) error

	// ForEach invokes the passed function with every key/value pair in the
	// bucket.  This does not include nested buckets or the key/value pairs
	// within those nested buckets.
	//
	// WARNING: It is not safe to mutate data while iterating with this
	// method.  Doing so may cause the underlying cursor to be invalidated
	// and return unexpected keys and/or values.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrTxClosed if the transaction has already been closed
	//
	// NOTE: The slices returned by this function are only valid during a
	// transaction.  Attempting to access them after a transaction has ended
	// results in undefined behavior.  Additionally, the slices must NOT
	// be modified by the caller.  These constraints prevent additional data
	// copies and allows support for memory-mapped database implementations.
	ForEach(func(k, v []byte) error) error

	// ForEachBucket invokes the passed function with the key of every
	// nested bucket in the current bucket.  This does not include any
	// nested buckets within those nested buckets.
	//
	// WARNING: It is not safe to mutate data while iterating with this
	// method.  Doing so may cause the underlying cursor to be invalidated
	// and return unexpected keys and/or values.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrTxClosed if the transaction has already been closed
	//
	// NOTE: The keys returned by this function are only valid during a
	// transaction.  Attempting to access them after a transaction has ended
	// results in undefined behavior.  This constraint prevents additional
	// data copies and allows support for memory-mapped database
	// implementations.
	ForEachBucket(func(k []byte) error) error

	// Cursor returns a new cursor, allowing for iteration over the bucket's
	// key/value pairs and nested buckets in forward or backward order.
	//
	// You must seek to a position using the First, Last, or Seek functions
	// before calling the Next, Prev, Key, or Value functions.  Failure to
	// do so will result in the same return values as an exhausted cursor,
	// which is false for the Prev and Next functions and nil for Key and
	// Value functions.
	Cursor() Cursor

	// Writable returns whether or not the bucket is writable.
	Writable() bool

	// Put saves the specified key/value pair to the bucket.  Keys that do
	// not already exist are added and keys that already exist are
	// overwritten.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrKeyRequired if the key is empty
	//   - ErrIncompatibleValue if the key is the same as an existing bucket
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	//
	// NOTE: The slices passed to this function must NOT be modified by the
	// caller.  This constraint prevents the requirement for additional data
	// copies and allows support for memory-mapped database implementations.
	Put(key, value []byte) error

	// Get returns the value for the given key.  Returns nil if the key does
	// not exist in this bucket.  An empty slice is returned for keys that
	// exist but have no value assigned.
	//
	// NOTE: The value returned by this function is only valid during a
	// transaction.  Attempting to access it after a transaction has ended
	// results in undefined behavior.  Additionally, the value must NOT
	// be modified by the caller.  These constraints prevent additional data
	// copies and allows support for memory-mapped database implementations.
	Get(key []byte) []byte

	// Delete removes the specified key from the bucket.  Deleting a key
	// that does not exist does not return an error.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrKeyRequired if the key is empty
	//   - ErrIncompatibleValue if the key is the same as an existing bucket
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	Delete(key []byte) error
}


// Cursor represents a cursor over key/value pairs and nested buckets of a
// bucket.
//
// Note that open cursors are not tracked on bucket changes and any
// modifications to the bucket, with the exception of Cursor.Delete, invalidates
// the cursor.  After invalidation, the cursor must be repositioned, or the keys
// and values returned may be unpredictable.
type Cursor interface {
	// Bucket returns the bucket the cursor was created for.
	Bucket() Bucket

	// Delete removes the current key/value pair the cursor is at without
	// invalidating the cursor.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrIncompatibleValue if attempted when the cursor points to a
	//     nested bucket
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	Delete() error

	// First positions the cursor at the first key/value pair and returns
	// whether or not the pair exists.
	First() bool

	// Last positions the cursor at the last key/value pair and returns
	// whether or not the pair exists.
	Last() bool

	// Next moves the cursor one key/value pair forward and returns whether
	// or not the pair exists.
	Next() bool

	// Prev moves the cursor one key/value pair backward and returns whether
	// or not the pair exists.
	Prev() bool

	// Seek positions the cursor at the first key/value pair that is greater
	// than or equal to the passed seek key.  Returns whether or not the
	// pair exists.
	Seek(seek []byte) bool

	// Key returns the current key the cursor is pointing to.
	Key() []byte

	// Value returns the current value the cursor is pointing to.  This will
	// be nil for nested buckets.
	Value() []byte
}

/*
// NewLDBDatabase returns a LevelDB wrapped object.
func NewLDBDatabase(file string, cache int, handles int) (*LDBDatabase, error) {
	logger := log.New("database", file)

	// Ensure we have some minimal caching and file guarantees
	if cache < 16 {
		cache = 16
	}
	if handles < 16 {
		handles = 16
	}
	logger.Info("Allocated cache and file handles", "cache", cache, "handles", handles)

	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(file, &opt.Options{
		OpenFilesCacheCapacity: handles,
		BlockCacheCapacity:     cache / 2 * opt.MiB,
		WriteBuffer:            cache / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(file, nil)
	}
	// (Re)check for errors and abort if opening of the db failed
	if err != nil {
		return nil, err
	}
	return &LDBDatabase{
		fn:  file,
		db:  db,
		log: logger,
	}, nil
}*/