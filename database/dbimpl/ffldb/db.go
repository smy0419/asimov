// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package ffldb

import (
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"hash/crc32"
	"os"
	"path/filepath"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/AsimovNetwork/asimov/database"
	"github.com/AsimovNetwork/asimov/database/dbimpl/ffldb/treap"
)

// db represents a collection of namespaces which are persisted and implements
// the database.Database interface.  All database access is performed through
// transactions which are obtained through the specific Namespace.
type db struct {
	writeLock sync.Mutex   // Limit to one write transaction at a time.
	closeLock sync.RWMutex // Make database close block while txns active.
	closed    bool         // Is the database closed?
	store     *blockStore  // Handles read/writing blocks to flat files.
	cache     *dbCache     // Cache layer which wraps underlying leveldb DB.
}

// Enforce db implements the database.Database interface.
var _ database.Database = (*db)(nil)

// Type returns the database driver type the current database instance was
// created with.
//
// This function is part of the database.Database interface implementation.
func (db *db) Type() string {
	return database.FFLDB
}

func (db *db) Put(key []byte, value []byte) error {
	return database.ErrUnimplement
}

func (db *db) Delete(key []byte) error {
	return database.ErrUnimplement
}

func (db *db) Get(key []byte) ([]byte, error) {
	return nil, database.ErrUnimplement
}

func (db *db) Has(key []byte) (bool, error) {
	return false, database.ErrUnimplement
}

func (db *db) NewBatch() database.Batch {
	return nil
}

// Begin starts a transaction which is either read-only or read-write depending
// on the specified flag.  Multiple read-only transactions can be started
// simultaneously while only a single read-write transaction can be started at a
// time.  The call will block when starting a read-write transaction when one is
// already open.
//
// NOTE: The transaction must be closed by calling Rollback or Commit on it when
// it is no longer needed.  Failure to do so will result in unclaimed memory.
//
// This function is part of the database.Database interface implementation.
func (db *db) Begin(writable bool) (database.Tx, error) {
	return db.begin(writable)
}

// View invokes the passed function in the context of a managed read-only
// transaction with the root bucket for the namespace.  Any errors returned from
// the user-supplied function are returned from this function.
//
// This function is part of the database.Database interface implementation.
func (db *db) View(fn func(database.Tx) error) error {
	// Start a read-only transaction.
	tx, err := db.begin(false)
	if err != nil {
		return err
	}

	// Since the user-provided function might panic, ensure the transaction
	// releases all mutexes and resources.  There is no guarantee the caller
	// won't use recover and keep going.  Thus, the database must still be
	// in a usable state on panics due to caller issues.
	defer rollbackOnPanic(tx)

	tx.managed = true
	err = fn(tx)
	tx.managed = false
	if err != nil {
		// The error is ignored here because nothing was written yet
		// and regardless of a rollback failure, the tx is closed now
		// anyways.
		_ = tx.Rollback()
		return err
	}

	return tx.Rollback()
}

// Update invokes the passed function in the context of a managed read-write
// transaction with the root bucket for the namespace.  Any errors returned from
// the user-supplied function will cause the transaction to be rolled back and
// are returned from this function.  Otherwise, the transaction is committed
// when the user-supplied function returns a nil error.
//
// This function is part of the database.Database interface implementation.
func (db *db) Update(fn func(database.Tx) error) error {
	// Start a read-write transaction.
	tx, err := db.begin(true)
	if err != nil {
		return err
	}

	// Since the user-provided function might panic, ensure the transaction
	// releases all mutexes and resources.  There is no guarantee the caller
	// won't use recover and keep going.  Thus, the database must still be
	// in a usable state on panics due to caller issues.
	defer rollbackOnPanic(tx)

	tx.managed = true
	err = fn(tx)
	tx.managed = false
	if err != nil {
		// The error is ignored here because nothing was written yet
		// and regardless of a rollback failure, the tx is closed now
		// anyways.
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Close cleanly shuts down the database and syncs all data.  It will block
// until all database transactions have been finalized (rolled back or
// committed).
//
// This function is part of the database.Database interface implementation.
func (db *db) Close() error {
	// Since all transactions have a read lock on this mutex, this will
	// cause Close to wait for all readers to complete.
	db.closeLock.Lock()
	defer db.closeLock.Unlock()

	if db.closed {
		return database.MakeError(database.ErrDbNotOpen, errDbNotOpenStr, nil)
	}
	db.closed = true

	// NOTE: Since the above lock waits for all transactions to finish and
	// prevents any new ones from being started, it is safe to flush the
	// cache and clear all state without the individual locks.

	// Close the database cache which will flush any existing entries to
	// disk and close the underlying leveldb database.  Any error is saved
	// and returned at the end after the remaining cleanup since the
	// database will be marked closed even if this fails given there is no
	// good way for the caller to recover from a failure here anyways.
	closeErr := db.cache.Close()

	// Close any open flat files that house the blocks.
	wc := db.store.writeCursor
	if wc.curFile.file != nil {
		_ = wc.curFile.file.Close()
		wc.curFile.file = nil
	}
	for _, blockFile := range db.store.openBlockFiles {
		_ = blockFile.file.Close()
	}
	db.store.openBlockFiles = nil
	db.store.openBlocksLRU.Init()
	db.store.fileNumToLRUElem = nil

	return closeErr
}

// begin is the implementation function for the Begin database method.  See its
// documentation for more details.
//
// This function is only separate because it returns the internal transaction
// which is used by the managed transaction code while the database method
// returns the interface.
func (db *db) begin(writable bool) (*transaction, error) {
	// Whenever a new writable transaction is started, grab the write lock
	// to ensure only a single write transaction can be active at the same
	// time.  This lock will not be released until the transaction is
	// closed (via Rollback or Commit).
	if writable {
		db.writeLock.Lock()
	}

	// Whenever a new transaction is started, grab a read lock against the
	// database to ensure Close will wait for the transaction to finish.
	// This lock will not be released until the transaction is closed (via
	// Rollback or Commit).
	db.closeLock.RLock()
	if db.closed {
		db.closeLock.RUnlock()
		if writable {
			db.writeLock.Unlock()
		}
		return nil, database.MakeError(database.ErrDbNotOpen, errDbNotOpenStr,
			nil)
	}

	// Grab a snapshot of the database cache (which in turn also handles the
	// underlying database).
	snapshot, err := db.cache.Snapshot()
	if err != nil {
		db.closeLock.RUnlock()
		if writable {
			db.writeLock.Unlock()
		}

		return nil, err
	}

	// The metadata and block index buckets are internal-only buckets, so
	// they have defined IDs.
	tx := &transaction{
		writable:      writable,
		db:            db,
		snapshot:      snapshot,
		pendingKeys:   treap.NewMutable(),
		pendingRemove: treap.NewMutable(),
	}
	tx.metaBucket = &bucket{tx: tx, id: metadataBucketID}
	tx.blockIdxBucket = &bucket{tx: tx, id: blockIdxBucketID}
	return tx, nil
}

// rollbackOnPanic rolls the passed transaction back if the code in the calling
// function panics.  This is needed since the mutex on a transaction must be
// released and a panic in called code would prevent that from happening.
//
// NOTE: This can only be handled manually for managed transactions since they
// control the life-cycle of the transaction.  As the documentation on Begin
// calls out, callers opting to use manual transactions will have to ensure the
// transaction is rolled back on panic if it desires that functionality as well
// or the database will fail to close since the read-lock will never be
// released.
func rollbackOnPanic(tx *transaction) {
	if err := recover(); err != nil {
		tx.managed = false
		_ = tx.Rollback()
		panic(err)
	}
}

// filesExists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// initDB creates the initial buckets and values used by the package.  This is
// mainly in a separate function for testing purposes.
func initDB(ldb *leveldb.DB) error {
	// The starting block file write cursor location is file num 0, offset
	// 0.
	batch := new(leveldb.Batch)
	batch.Put(bucketizedKey(metadataBucketID, writeLocKeyName),
		serializeWriteRow(0, 0))

	// Create block index bucket and set the current bucket id.
	//
	// NOTE: Since buckets are virtualized through the use of prefixes,
	// there is no need to store the bucket index data for the metadata
	// bucket in the database.  However, the first bucket ID to use does
	// need to account for it to ensure there are no key collisions.
	batch.Put(bucketIndexKey(metadataBucketID, blockIdxBucketName),
		blockIdxBucketID[:])
	batch.Put(curBucketIDKeyName, blockIdxBucketID[:])

	// Write everything as a single batch.
	if err := ldb.Write(batch, nil); err != nil {
		str := fmt.Sprintf("failed to initialize metadata database: %v",
			err)
		return database.ConvertErr(str, err)
	}

	return nil
}

// openDB opens the database at the provided path.  database.ErrDbDoesNotExist
// is returned if the database doesn't exist and the create flag is not set.
func openDB(dbPath string, network common.AsimovNet, create bool) (database.Database, error) {
	// Error if the database doesn't exist and the create flag is not set.
	metadataDbPath := filepath.Join(dbPath, metadataDbName)
	dbExists := fileExists(metadataDbPath)
	if !create && !dbExists {
		str := fmt.Sprintf("database %q does not exist", metadataDbPath)
		return nil, database.MakeError(database.ErrDbDoesNotExist, str, nil)
	}

	// Ensure the full path to the database exists.
	if !dbExists {
		// The error can be ignored here since the call to
		// leveldb.OpenFile will fail if the directory couldn't be
		// created.
		_ = os.MkdirAll(dbPath, 0700)
	}

	// Open the metadata database (will create it if needed).
	opts := opt.Options{
		ErrorIfExist: create,
		Strict:       opt.DefaultStrict,
		Compression:  opt.NoCompression,
		Filter:       filter.NewBloomFilter(10),
	}
	ldb, err := leveldb.OpenFile(metadataDbPath, &opts)
	if err != nil {
		return nil, database.ConvertErr(err.Error(), err)
	}

	// Create the block store which includes scanning the existing flat
	// block files to find what the current write cursor position is
	// according to the data that is actually on disk.  Also create the
	// database cache which wraps the underlying leveldb database to provide
	// write caching.
	store := newBlockStore(dbPath, network)
	cache := newDbCache(ldb, store, defaultCacheSize, defaultFlushSecs)
	pdb := &db{store: store, cache: cache}

	// Perform any reconciliation needed between the block and metadata as
	// well as database initialization, if needed.
	return reconcileDB(pdb, create)
}

// The serialized write cursor location format is:
//
//  [0:4]  Block file (4 bytes)
//  [4:8]  File offset (4 bytes)
//  [8:12] Castagnoli CRC-32 checksum (4 bytes)

// serializeWriteRow serialize the current block file and offset where new
// will be written into a format suitable for storage into the metadata.
func serializeWriteRow(curBlockFileNum, curFileOffset uint32) []byte {
	var serializedRow [12]byte
	byteOrder.PutUint32(serializedRow[0:4], curBlockFileNum)
	byteOrder.PutUint32(serializedRow[4:8], curFileOffset)
	checksum := crc32.Checksum(serializedRow[:8], castagnoli)
	byteOrder.PutUint32(serializedRow[8:12], checksum)
	return serializedRow[:]
}

// deserializeWriteRow deserializes the write cursor location stored in the
// metadata.  Returns ErrCorruption if the checksum of the entry doesn't match.
func deserializeWriteRow(writeRow []byte) (uint32, uint32, error) {
	// Ensure the checksum matches.  The checksum is at the end.
	gotChecksum := crc32.Checksum(writeRow[:8], castagnoli)
	wantChecksumBytes := writeRow[8:12]
	wantChecksum := byteOrder.Uint32(wantChecksumBytes)
	if gotChecksum != wantChecksum {
		str := fmt.Sprintf("metadata for write cursor does not match "+
			"the expected checksum - got %d, want %d", gotChecksum,
			wantChecksum)
		return 0, 0, database.MakeError(database.ErrCorruption, str, nil)
	}

	fileNum := byteOrder.Uint32(writeRow[0:4])
	fileOffset := byteOrder.Uint32(writeRow[4:8])
	return fileNum, fileOffset, nil
}

// reconcileDB reconciles the metadata with the flat block files on disk.  It
// will also initialize the underlying database if the create flag is set.
func reconcileDB(pdb *db, create bool) (database.Database, error) {
	// Perform initial internal bucket and value creation during database
	// creation.
	if create {
		if err := initDB(pdb.cache.ldb); err != nil {
			return nil, err
		}
	}

	// Load the current write cursor position from the metadata.
	var curFileNum, curOffset uint32
	err := pdb.View(func(tx database.Tx) error {
		writeRow := tx.Metadata().Get(writeLocKeyName)
		if writeRow == nil {
			str := "write cursor does not exist"
			return database.MakeError(database.ErrCorruption, str, nil)
		}

		var err error
		curFileNum, curOffset, err = deserializeWriteRow(writeRow)
		return err
	})
	if err != nil {
		return nil, err
	}

	// When the write cursor position found by scanning the block files on
	// disk is AFTER the position the metadata believes to be true, truncate
	// the files on disk to match the metadata.  This can be a fairly common
	// occurrence in unclean shutdown scenarios while the block files are in
	// the middle of being written.  Since the metadata isn't updated until
	// after the block data is written, this is effectively just a rollback
	// to the known good point before the unclean shutdown.
	wc := pdb.store.writeCursor
	if wc.curFileNum > curFileNum || (wc.curFileNum == curFileNum &&
		wc.curOffset > curOffset) {

		log.Info("Detected unclean shutdown - Repairing...")
		log.Debugf("Metadata claims file %d, offset %d. Block data is "+
			"at file %d, offset %d", curFileNum, curOffset,
			wc.curFileNum, wc.curOffset)
		pdb.store.handleRollback(curFileNum, curOffset)
		log.Infof("Database sync complete")
	}

	// When the write cursor position found by scanning the block files on
	// disk is BEFORE the position the metadata believes to be true, return
	// a corruption error.  Since sync is called after each block is written
	// and before the metadata is updated, this should only happen in the
	// case of missing, deleted, or truncated block files, which generally
	// is not an easily recoverable scenario.  In the future, it might be
	// possible to rescan and rebuild the metadata from the block files,
	// however, that would need to happen with coordination from a higher
	// layer since it could invalidate other metadata.
	if wc.curFileNum < curFileNum || (wc.curFileNum == curFileNum &&
		wc.curOffset < curOffset) {

		str := fmt.Sprintf("metadata claims file %d, offset %d, but "+
			"block data is at file %d, offset %d", curFileNum,
			curOffset, wc.curFileNum, wc.curOffset)
		log.Warnf("***Database corruption detected***: %v", str)
		return nil, database.MakeError(database.ErrCorruption, str, nil)
	}

	return pdb, nil
}

const (
	dbType = database.FFLDB
)

// parseArgs parses the arguments from the database Open/Create methods.
func parseArgs(funcName string, args ...interface{}) (string, common.AsimovNet, error) {
	if len(args) != 2 {
		return "", 0, fmt.Errorf("invalid arguments to %s.%s -- "+
			"expected database path and block network", dbType,
			funcName)
	}

	dbPath, ok := args[0].(string)
	if !ok {
		return "", 0, fmt.Errorf("first argument to %s.%s is invalid -- "+
			"expected database path string", dbType, funcName)
	}

	network, ok := args[1].(common.AsimovNet)
	if !ok {
		return "", 0, fmt.Errorf("second argument to %s.%s is invalid -- "+
			"expected block network", dbType, funcName)
	}

	return dbPath, network, nil
}

// openDBDriver is the callback provided during driver registration that opens
// an existing database for use.
func OpenDBDriver(args ...interface{}) (database.Database, error) {
	dbPath, network, err := parseArgs("Open", args...)
	if err != nil {
		return nil, err
	}

	return openDB(dbPath, network, false)
}

// createDBDriver is the callback provided during driver registration that
// creates, initializes, and opens a database for use.
func CreateDBDriver(args ...interface{}) (database.Database, error) {
	dbPath, network, err := parseArgs("Create", args...)
	if err != nil {
		return nil, err
	}

	return openDB(dbPath, network, true)
}
