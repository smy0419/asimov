// Copyright (c) 2018-2020 The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package test

import (
	"github.com/AsimovNetwork/asimov/asiutil"
	"github.com/AsimovNetwork/asimov/chaincfg"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/database"
	"testing"
)

// test base interface, include Putter,Deleter,Queryer
func TestDbBase( t *testing.T ) {
	internalTestDbBase( t, false )
	internalTestDbBase( t, true )
}


func internalTestDbBase( t *testing.T, isMemdb bool ) {
	dbGenerator := NewDbGenerator( database.ETHDB, common.TestNet, isMemdb )
	if dbGenerator == nil {
		t.Error( "NewDbGenerator fail" )
		return
	}

	err := dbGenerator.Generate()
	if err != nil {
		t.Error( err.Error() )
		return
	}

	defer dbGenerator.db.Close()

	workFlow := NewWorkFlow( dbGenerator.db )
	dbBaseAction := NewDbBaseAction()
	workFlow.addAction( dbBaseAction )

	bytesGenerator := NewBytesGenerator()

	// create 200 db put case, key and val is generated randomly
	keys := make( [][]byte, 0, 0 )
	for i := 1; i < 200; i++ {
		keys = append( keys, bytesGenerator.CustomGenerate( uint32( i ) ) )
		testCase := NewDbPutCase( keys[i-1], bytesGenerator.CustomGenerate( uint32( 50 * i ) ) )
		dbBaseAction.addCase( testCase )
	}

	// create key1, val1 for testing put、query、delete
	key1 := bytesGenerator.CustomGenerate(16 )
	val1 := bytesGenerator.CustomGenerate(160 )
	testCase1 := NewDbPutCase( key1, val1 )
	testCase2 := NewDbQueryCase( key1 )
	testCase3 := NewDbDeleteCase( key1 )

	// create key2, val2 that length is bigger for testing put、query、delete
	key2 := bytesGenerator.CustomGenerate( 100 )
	val2 := bytesGenerator.CustomGenerate( 512 )
	testCase4 := NewDbPutCase( key2, val2 )
	testCase5 := NewDbQueryCase( key2 )
	testCase6 := NewDbDeleteCase( key2 )

	// test empty key and val
	testCase7 := NewDbPutCase( []byte{}, []byte{} )
	testCase8 := NewDbQueryCase( []byte{} )
	testCase9 := NewDbDeleteCase( []byte{} )

	// delete 200 keys that put before
	for i := 1; i < len(keys); i++ {
		testCase := NewDbDeleteCase( keys[i] )
		dbBaseAction.addCase( testCase )
	}

	// normal
	//dbBaseAction.addCase( testCase2 )
	dbBaseAction.addCase( testCase1 )
	dbBaseAction.addCase( testCase2 )
	dbBaseAction.addCase( testCase3 )
	dbBaseAction.addCase( testCase2 )
	dbBaseAction.addCase( testCase4 )
	dbBaseAction.addCase( testCase5 )
	dbBaseAction.addCase( testCase6 )
	dbBaseAction.addCase( testCase5 )

	// delete not exists
	dbBaseAction.addCase( testCase3 )
	dbBaseAction.addCase( testCase6 )

	// query not exists
	dbBaseAction.addCase( testCase2 )
	dbBaseAction.addCase( testCase5 )

	// put exists
	dbBaseAction.addCase( testCase1 )
	dbBaseAction.addCase( testCase1 )
	dbBaseAction.addCase( testCase4 )
	dbBaseAction.addCase( testCase4 )

	// param is empty
	dbBaseAction.addCase( testCase7 )
	dbBaseAction.addCase( testCase8 )
	dbBaseAction.addCase( testCase9 )

	err = workFlow.execute()
	if err != nil {
		t.Error( err.Error() )
	}
}

// test transactor interface
func TestDbTransactor( t *testing.T ) {
	internalTestDbTransactor( t, false )
	//internalTestDbTransactor( t, true )
}

func internalTestDbTransactor( t *testing.T, isMemdb bool ) {
	dbGenerator := NewDbGenerator( database.FFLDB, common.TestNet, isMemdb )
	if dbGenerator == nil {
		t.Error("NewDbGenerator fail" )
		return
	}

	err := dbGenerator.Generate()
	if err != nil {
		t.Error( err.Error() )
		return
	}

	defer dbGenerator.db.Close()
	workFlow := NewWorkFlow( dbGenerator.db )

	// add update action
	updateAction := NewDbTxUpdateAction()
	workFlow.addAction( updateAction )

	// add view action
	viewAction := NewDbTxViewAction()
	workFlow.addAction( viewAction )

	bytesGenerator := NewBytesGenerator()

	// create 200 store block cases
	blockHashes := make( []*common.Hash, 0, 0 )
	for i := 1; i < 200; i++ {
		blockHash := bytesGenerator.GenerateHash()
		blockHashes = append( blockHashes, &blockHash )

		blockBytes := bytesGenerator.CustomGenerate( uint32( i * 512 ) )
		testCase := NewStoreBlockCase( blockHashes[i-1], blockBytes, i % 2 == 0 )
		updateAction.addCase( testCase )
	}


	// create 200 has block cases
	// create 200 fetch block header cases
	// create 200 fetch block cases
	// create 200 fetch block region cases
	for i := 1; i < len(blockHashes); i++ {
		viewAction.addCase( NewHasBlockCase( blockHashes[i] ) )
		viewAction.addCase( NewFetchBlockHeaderCase( blockHashes[i] ) )
		viewAction.addCase( NewFetchBlockCase( blockHashes[i] ) )
		viewAction.addCase( NewFetchBlockRegionCase( blockHashes[i], 0, 10000 ) )
	}


	blockHash1 := bytesGenerator.GenerateHash()
	blockBytes1 := bytesGenerator.CustomGenerate( 1024 * 1024 )
	blockHash2 := bytesGenerator.GenerateHash()
	blockBytes2 := bytesGenerator.CustomGenerate( 512 *100 )
	testCase1 := NewStoreBlockCase( &blockHash1, blockBytes1, false )
	testCase2 := NewHasBlockCase( &blockHash1 )
	testCase4 := NewFetchBlockHeaderCase( &blockHash1 )
	testCase5 := NewFetchBlockRegionCase( &blockHash1, 0, 512 )
	testCase6 := NewFetchBlockCase( &blockHash1 )
	testCase7 := NewStoreBlockCase( &blockHash2, blockBytes2, false )
	testCase8 := NewStoreBlockCase( nil, []byte{}, false )

	// add 2 blocks
	updateAction.addCase( testCase1 )
	updateAction.addCase( testCase7 )
	updateAction.addCase( testCase8 )

	// query blocks
	viewAction.addCase( testCase1 )
	viewAction.addCase( testCase2 )
	viewAction.addCase( testCase4 )
	viewAction.addCase( testCase5 )
	viewAction.addCase( testCase6 )


	// test two blocks
	block1 := asiutil.NewBlock( chaincfg.MainNetParams.GenesisBlock )
	block1Bytes, _ := block1.Bytes()
	updateAction.addCase( NewStoreBlockCase( block1.Hash(), block1Bytes, false ) )
	viewAction.addCase( NewFetchBlockCase( block1.Hash() ) )
	viewAction.addCase( NewFetchBlockHeaderCase( block1.Hash() ) )

	block2 := asiutil.NewBlock( chaincfg.TestNetParams.GenesisBlock )
	block2Bytes, _ := block2.Bytes()
	updateAction.addCase( NewStoreBlockCase( block2.Hash(), block2Bytes, false ) )
	viewAction.addCase( NewFetchBlockCase( block2.Hash() ) )
	viewAction.addCase( NewFetchBlockHeaderCase( block2.Hash() ) )

	err = workFlow.execute()
	if err != nil  {
		t.Error( err.Error() )
	}
}


// test batch interface
func TestDbBatch( t *testing.T ) {
	internalTestDbBatch( t, false )
	internalTestDbBatch( t, true )
}


func internalTestDbBatch( t *testing.T, isMemdb bool ) {
	dbGenerator := NewDbGenerator( database.ETHDB, common.TestNet, isMemdb )
	if dbGenerator == nil {
		t.Error( "NewDbGenerator fail" )
		return
	}

	err := dbGenerator.Generate()
	if err != nil {
		t.Error( err.Error() )
		return
	}

	defer dbGenerator.db.Close()

	workFlow := NewWorkFlow( dbGenerator.db )

	batchAction := NewDbBatchAction( dbGenerator.db )
	workFlow.addAction( batchAction )

	dbBaseAction := NewDbBaseAction()
	workFlow.addAction( dbBaseAction )

	bytesGenerator := NewBytesGenerator()
	key1 := bytesGenerator.CustomGenerate( 100 )
	val1 := bytesGenerator.CustomGenerate( 1024 * 100 )
	key2 := bytesGenerator.CustomGenerate( 200 )
	val2 := bytesGenerator.CustomGenerate( 1024 * 512 )

	testCase1 := NewBatchPutCase( key1, val1 )
	testCase2 := NewDbBatchDeleteCase( key1 )
	testCase3 := NewDbBatchWriteCase()
	testCase4 := NewDbBatchResetCase()
	testCase5 := NewBatchPutCase( key2, val2 )

	testCase6 := NewDbQueryCase( key1 )
	testCase7 := NewDbQueryCase( key2 )

	// add op
	batchAction.addCase( testCase1 )
	batchAction.addCase( testCase2 )
	batchAction.addCase( testCase1 )
	// write
	batchAction.addCase( testCase3 )

	// reset
	batchAction.addCase( testCase4 )

	// add op
	batchAction.addCase( testCase5 )
	batchAction.addCase( testCase1 )
	// write
	batchAction.addCase( testCase3 )

	// base query
	dbBaseAction.addCase( testCase6 )
	dbBaseAction.addCase( testCase7 )

	// reset
	batchAction.addCase( testCase4 )

	//
	keys := make( [][]byte, 0, 0 )
	for i := 1; i < 100; i++ {
		key := bytesGenerator.CustomGenerate( uint32( 100+i ) )

		val := bytesGenerator.CustomGenerate( uint32( 100*i ) )
		testCase := NewBatchPutCase( key, val )
		batchAction.addCase( testCase )

		if i < 20 {
			batchAction.addCase( NewDbBatchDeleteCase( key ) )
		} else {
			keys = append( keys, key )
		}
	}

	for i := 1; i < 10; i++ {
		testCase := NewDbBatchDeleteCase( keys[i-1] )
		batchAction.addCase( testCase )
	}

	err = workFlow.execute()
	if err != nil {
		t.Error( err.Error() )
	}
}

// test bucket interface
func TestDbBucket( t *testing.T ) {
	dbGenerator := NewDbGenerator( database.FFLDB, common.TestNet, false )
	if dbGenerator == nil {
		t.Error( "NewDbGenerator fail" )
		return
	}

	err := dbGenerator.Generate()
	if err != nil {
		t.Error( err.Error() )
		return
	}

	defer dbGenerator.db.Close()

	workFlow := NewWorkFlow( dbGenerator.db )
	bucketAction := NewDbBucketAction( dbGenerator.db, true )
	workFlow.addAction( bucketAction )

	bytesGenerator := NewBytesGenerator()

	// 100 create bucket
	bucketKeys := make( [][]byte, 0, 0 )
	for i := 1; i < 100; i++ {
		bucketKeys = append( bucketKeys, bytesGenerator.CustomGenerate( uint32(i) ) )
		testCase := NewDbBucketCreateCase( bucketKeys[i-1] )
		bucketAction.addCase( testCase )
	}

	// 100 create bucket if not exist
	for i := 1; i < 100; i++ {
		testCase := NewDbBucketCreateIfNotExistsCase( bytesGenerator.CustomGenerate( uint32(i) ) )
		bucketAction.addCase( testCase )
	}

	// 100 put in bucket
	keys := make( [][]byte, 0, 0 )
	for i := 1; i < 100; i++ {
		keys = append( keys, bytesGenerator.CustomGenerate( uint32(i) ) )
		testCase := NewDbBucketPutCase( keys[i-1], bytesGenerator.CustomGenerate( uint32( i * 50 ) ) )
		bucketAction.addCase( testCase )
	}


	key1 := bytesGenerator.CustomGenerate(50 )
	key2 := bytesGenerator.CustomGenerate(512 )
	testCase1 := NewDbBucketCreateCase( key1 )
	testCase2 := NewDbBucketForEachCase()
	testCase3 := NewDbBucketCreateIfNotExistsCase( key1 )
	testCase7 := NewDbBucketCreateIfNotExistsCase( key2 )
	testCase4 := NewDbBucketDeleteBucketCase( key1 )
	testCase5 := NewDbbucketForEachbucketCase()
	testCase6 := NewDbBucketDeleteCase( key2 )

	bucketAction.addCase( testCase2 )

	// delete keys
	for i := 0; i < len(keys); i++ {
		testCase := NewDbBucketDeleteCase( keys[i] )
		bucketAction.addCase( testCase )
	}

	bucketAction.addCase( testCase2 )
	bucketAction.addCase( testCase1 )
	bucketAction.addCase( testCase2 )
	bucketAction.addCase( testCase3 )
	bucketAction.addCase( testCase7 )
	bucketAction.addCase( testCase4 )
	bucketAction.addCase( testCase5 )

	// delete buckets
	for i := 1; i < len(bucketKeys); i++ {
		testCase := NewDbBucketDeleteBucketCase( bucketKeys[i] )
		bucketAction.addCase( testCase )
	}

	bucketAction.addCase( testCase5 )
	bucketAction.addCase( testCase6 )
	bucketAction.addCase( testCase1 )

	err = workFlow.execute()

	if err != nil {
		t.Error( err.Error() )
	}
}

// test cursor interface
func TestDbCursor( t *testing.T ) {
	dbGenerator := NewDbGenerator( database.FFLDB, common.TestNet, false )
	if dbGenerator == nil {
		t.Error( "NewDbGenerator fail" )
		return
	}

	err := dbGenerator.Generate()
	if err != nil {
		t.Error( err.Error() )
		return
	}

	defer dbGenerator.db.Close()

	workFlow := NewWorkFlow( dbGenerator.db )
	cursorAction := NewDbCursorAction( dbGenerator.db, true )
	workFlow.addAction( cursorAction )

	testCase1 := NewDbCursorDeleteCase()
	testCase2 := NewDbCursorFirstCase()
	testCase3 := NewDbCursorLastCase()
	testCase4 := NewDbCursorNextCase()
	testCase5 := NewDbCursorPrevCase( )
	testCase6 := NewDbCursorSeekCase([]byte{} )
	testCase7 := NewDbCursorKeyCase()
	testCase8 := NewDbCursorValueCase()

	cursorAction.addCase( testCase1 )

	for i := 1; i < 100; i++ {
		cursorAction.addCase( testCase2 )
	}

	for i := 1; i < 100; i++ {
		cursorAction.addCase( testCase3 )
	}

	for i := 1; i < 100; i++ {
		cursorAction.addCase( testCase4 )
	}

	for i := 1; i < 100; i++ {
		cursorAction.addCase( testCase5 )
	}

	for i := 1; i < 100; i++ {
		cursorAction.addCase( testCase6 )
	}

	for i := 1; i < 100; i++ {
		cursorAction.addCase( testCase7 )
	}

	for i := 1; i < 100; i++ {
		cursorAction.addCase( testCase8 )
	}

	for i := 1; i < 100; i++ {
		cursorAction.addCase( testCase1 )
	}


	err = workFlow.execute()

	if err != nil {
		t.Error( err.Error() )
	}
}