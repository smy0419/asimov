// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dbdriver

import (
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/AsimovNetwork/asimov/database"
)

var (
	// ignoreDbTypes are types which should be ignored when running tests
	// that iterate all supported DB types.  This allows some tests to add
	// bogus drivers for testing purposes while still allowing other tests
	// to easily iterate all supported drivers.
	ignoreDbTypes = map[string]bool{"createopenfail": true}
	dbType = database.FFLDB
)

// checkDbError ensures the passed error is a database.Error with an error code
// that matches the passed  error code.
func checkDbError(t *testing.T, testName string, gotErr error, wantErrCode database.ErrorCode) bool {
	dbErr, ok := gotErr.(database.Error)
	if !ok {
		t.Errorf("%s: unexpected error type - got %T, want %T",
			testName, gotErr, database.Error{})
		return false
	}
	if dbErr.ErrorCode != wantErrCode {
		t.Errorf("%s: unexpected error code - got %s (%s), want %s",
			testName, dbErr.ErrorCode, dbErr.Description,
			wantErrCode)
		return false
	}

	return true
}

// TestAddDuplicateDriver ensures that adding a duplicate driver does not
// overwrite an existing one.
func TestAddDuplicateDriver(t *testing.T) {
	supportedDrivers := SupportedDrivers()
	if len(supportedDrivers) == 0 {
		t.Errorf("no backends to test")
		return
	}
	dbType := supportedDrivers[0]

	// bogusCreateDB is a function which acts as a bogus create and open
	// driver function and intentionally returns a failure that can be
	// detected if the interface allows a duplicate driver to overwrite an
	// existing one.
	bogusCreateDB := func(args ...interface{}) (database.Database, error) {
		return nil, fmt.Errorf("duplicate driver allowed for database "+
			"type [%v]", dbType)
	}

	// Create a driver that tries to replace an existing one.  Set its
	// create and open functions to a function that causes a test failure if
	// they are invoked.
	driver := Driver{
		DbType: dbType,
		Create: bogusCreateDB,
		Open:   bogusCreateDB,
	}
	testName := "duplicate driver registration"
	err := RegisterDriver(&driver)
	if !checkDbError(t, testName, err, database.ErrDbTypeRegistered) {
		return
	}
}

// TestCreateOpenFail ensures that errors which occur while opening or closing
// a database are handled properly.
func TestCreateOpenFail1(t *testing.T) {
	// bogusCreateDB is a function which acts as a bogus create and open
	// driver function that intentionally returns a failure which can be
	// detected.
	dbType := "createopenfail"
	openError := fmt.Errorf("failed to create or open database for "+
		"database type [%v]", dbType)
	bogusCreateDB := func(args ...interface{}) (database.Database, error) {
		return nil, openError
	}

	// Create and add driver that intentionally fails when created or opened
	// to ensure errors on database open and create are handled properly.
	driver := Driver{
		DbType: dbType,
		Create: bogusCreateDB,
		Open:   bogusCreateDB,
	}
	RegisterDriver(&driver)

	// Ensure creating a database with the new type fails with the expected
	// error.
	_, err := Create(dbType)
	if err != openError {
		t.Errorf("expected error not received - got: %v, want %v", err,
			openError)
		return
	}

	// Ensure opening a database with the new type fails with the expected
	// error.
	_, err = Open(dbType)
	if err != openError {
		t.Errorf("expected error not received - got: %v, want %v", err,
			openError)
		return
	}
}

// TestCreateOpenUnsupported ensures that attempting to create or open an
// unsupported database type is handled properly.
func TestCreateOpenUnsupported(t *testing.T) {
	// Ensure creating a database with an unsupported type fails with the
	// expected error.
	testName := "create with unsupported database type"
	dbType := "unsupported"
	_, err := Create(dbType)
	if !checkDbError(t, testName, err, database.ErrDbUnknownType) {
		return
	}

	// Ensure opening a database with the an unsupported type fails with the
	// expected error.
	testName = "open with unsupported database type"
	_, err = Open(dbType)
	if !checkDbError(t, testName, err, database.ErrDbUnknownType) {
		return
	}
}


// dbType is the database type name for this driver.
const blockDataNet = common.TestNet

// TestCreateOpenFail ensures that errors related to creating and opening a
// database are handled properly.
func TestCreateOpenFail(t *testing.T) {
	t.Parallel()

	// Ensure that attempting to open a database that doesn't exist returns
	// the expected error.
	wantErrCode := database.ErrDbDoesNotExist
	_, err := Open(dbType, "noexist", blockDataNet)
	if !checkDbError(t, "Open", err, wantErrCode) {
		return
	}

	// Ensure that attempting to open a database with the wrong number of
	// parameters returns the expected error.
	wantErr := fmt.Errorf("invalid arguments to %s.Open -- expected "+
		"database path and block network", dbType)
	_, err = Open(dbType, 1, 2, 3)
	if err.Error() != wantErr.Error() {
		t.Errorf("Open: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure that attempting to open a database with an invalid type for
	// the first parameter returns the expected error.
	wantErr = fmt.Errorf("first argument to %s.Open is invalid -- "+
		"expected database path string", dbType)
	_, err = Open(dbType, 1, blockDataNet)
	if err.Error() != wantErr.Error() {
		t.Errorf("Open: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure that attempting to open a database with an invalid type for
	// the second parameter returns the expected error.
	wantErr = fmt.Errorf("second argument to %s.Open is invalid -- "+
		"expected block network", dbType)
	_, err = Open(dbType, "noexist", "invalid")
	if err.Error() != wantErr.Error() {
		t.Errorf("Open: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure that attempting to create a database with the wrong number of
	// parameters returns the expected error.
	wantErr = fmt.Errorf("invalid arguments to %s.Create -- expected "+
		"database path and block network", dbType)
	_, err = Create(dbType, 1, 2, 3)
	if err.Error() != wantErr.Error() {
		t.Errorf("Create: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure that attempting to create a database with an invalid type for
	// the first parameter returns the expected error.
	wantErr = fmt.Errorf("first argument to %s.Create is invalid -- "+
		"expected database path string", dbType)
	_, err = Create(dbType, 1, blockDataNet)
	if err.Error() != wantErr.Error() {
		t.Errorf("Create: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure that attempting to create a database with an invalid type for
	// the second parameter returns the expected error.
	wantErr = fmt.Errorf("second argument to %s.Create is invalid -- "+
		"expected block network", dbType)
	_, err = Create(dbType, "noexist", "invalid")
	if err.Error() != wantErr.Error() {
		t.Errorf("Create: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure operations against a closed database return the expected
	// error.
	dbPath := filepath.Join(os.TempDir(), "ffldb-createfail")
	_ = os.RemoveAll(dbPath)
	db, err := Create(dbType, dbPath, blockDataNet)
	if err != nil {
		t.Errorf("Create: unexpected error: %v", err)
		return
	}
	defer os.RemoveAll(dbPath)
	db.Close()

	wantErrCode = database.ErrDbNotOpen
	err = db.View(func(tx database.Tx) error {
		return nil
	})
	if !checkDbError(t, "View", err, wantErrCode) {
		return
	}

	wantErrCode = database.ErrDbNotOpen
	err = db.Update(func(tx database.Tx) error {
		return nil
	})
	if !checkDbError(t, "Update", err, wantErrCode) {
		return
	}

	wantErrCode = database.ErrDbNotOpen
	_, err = db.Begin(false)
	if !checkDbError(t, "Begin(false)", err, wantErrCode) {
		return
	}

	wantErrCode = database.ErrDbNotOpen
	_, err = db.Begin(true)
	if !checkDbError(t, "Begin(true)", err, wantErrCode) {
		return
	}

	wantErrCode = database.ErrDbNotOpen
	err = db.Close()
	if !checkDbError(t, "Close", err, wantErrCode) {
		return
	}
}

//// TestPersistence ensures that values stored are still valid after closing and
//// reopening the database.
//func TestPersistence(t *testing.T) {
//	t.Parallel()
//
//	// Create a new database to run tests against.
//	dbPath := filepath.Join(os.TempDir(), "ffldb-persistencetest")
//	_ = os.RemoveAll(dbPath)
//	db, err := Create(dbType, dbPath, blockDataNet)
//	if err != nil {
//		t.Errorf("Failed to create test database (%s) %v", dbType, err)
//		return
//	}
//	defer os.RemoveAll(dbPath)
//	defer db.Close()
//
//	// Create a bucket, put some values into it, and store a block so they
//	// can be tested for existence on re-open.
//	bucket1Key := []byte("bucket1")
//	storeValues := map[string]string{
//		"b1key1": "foo1",
//		"b1key2": "foo2",
//		"b1key3": "foo3",
//	}
//	genesisBlock := vvsutil.NewBlock(chaincfg.MainNetParams.GenesisBlock)
//	genesisHash := chaincfg.MainNetParams.GenesisHash
//	err = db.Update(func(tx database.Tx) error {
//		metadataBucket := tx.Metadata()
//		if metadataBucket == nil {
//			return fmt.Errorf("Metadata: unexpected nil bucket")
//		}
//
//		bucket1, err := metadataBucket.CreateBucket(bucket1Key)
//		if err != nil {
//			return fmt.Errorf("CreateBucket: unexpected error: %v",
//				err)
//		}
//
//		for k, v := range storeValues {
//			err := bucket1.Put([]byte(k), []byte(v))
//			if err != nil {
//				return fmt.Errorf("Put: unexpected error: %v",
//					err)
//			}
//		}
//
//		blockBytes, _ := genesisBlock.Bytes()
//		if err := tx.StoreBlock(genesisBlock.Hash(), blockBytes); err != nil {
//			return fmt.Errorf("StoreBlock: unexpected error: %v",
//				err)
//		}
//
//		return nil
//	})
//	if err != nil {
//		t.Errorf("Update: unexpected error: %v", err)
//		return
//	}
//
//	// Close and reopen the database to ensure the values persist.
//	db.Close()
//	db, err = Open(dbType, dbPath, blockDataNet)
//	if err != nil {
//		t.Errorf("Failed to open test database (%s) %v", dbType, err)
//		return
//	}
//	defer db.Close()
//
//	// Ensure the values previously stored in the 3rd namespace still exist
//	// and are correct.
//	err = db.View(func(tx database.Tx) error {
//		metadataBucket := tx.Metadata()
//		if metadataBucket == nil {
//			return fmt.Errorf("Metadata: unexpected nil bucket")
//		}
//
//		bucket1 := metadataBucket.Bucket(bucket1Key)
//		if bucket1 == nil {
//			return fmt.Errorf("Bucket1: unexpected nil bucket")
//		}
//
//		for k, v := range storeValues {
//			gotVal := bucket1.Get([]byte(k))
//			if !reflect.DeepEqual(gotVal, []byte(v)) {
//				return fmt.Errorf("Get: key '%s' does not "+
//					"match expected value - got %s, want %s",
//					k, gotVal, v)
//			}
//		}
//
//		genesisBlockBytes, _ := genesisBlock.Bytes()
//		gotBytes, err := tx.FetchBlock(genesisHash)
//		if err != nil {
//			return fmt.Errorf("FetchBlock: unexpected error: %v",
//				err)
//		}
//		if !reflect.DeepEqual(gotBytes, genesisBlockBytes) {
//			return fmt.Errorf("FetchBlock: stored block mismatch")
//		}
//
//		return nil
//	})
//	if err != nil {
//		t.Errorf("View: unexpected error: %v", err)
//		return
//	}
//}

// TestInterface performs all interfaces tests for this database driver.
func TestInterface(t *testing.T) {
	t.Parallel()

	// Create a new database to run tests against.
	dbPath := filepath.Join(os.TempDir(), "ffldb-interfacetest")
	_ = os.RemoveAll(dbPath)
	db, err := Create(dbType, dbPath, blockDataNet)
	if err != nil {
		t.Errorf("Failed to create test database (%s) %v", dbType, err)
		return
	}
	defer os.RemoveAll(dbPath)
	defer db.Close()

	// Ensure the driver type is the expected value.
	gotDbType := db.Type()
	if gotDbType != dbType {
		t.Errorf("Type: unepxected driver type - got %v, want %v",
			gotDbType, dbType)
		return
	}

	// Run all of the interface tests against the database.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Change the maximum file size to a small value to force multiple flat
	// files with the test data set.

	// csun TODO ???
	/*ffldb.TstRunWithMaxBlockFileSize(db, 2048, func() {
		testInterface(t, db)
	})*/
}

//func TstRunWithMaxBlockFileSize(idb database.Database, size uint32, fn func()) {
//	ffldb := idb.(*db)
//	origSize := ffldb.store.maxBlockFileSize
//
//	ffldb.store.maxBlockFileSize = size
//	fn()
//	ffldb.store.maxBlockFileSize = origSize
//}


