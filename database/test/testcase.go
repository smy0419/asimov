// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package test

import (
	"errors"
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/database"
)

type TestCase interface {
	execute( ctx DbContext ) error
	isWritable() bool
}

// StoreBlock
type TxStoreBlockCase struct {
	blockHash *common.Hash
	blockBytes []byte

	autoRollback bool
}

func NewStoreBlockCase(
	blockHash *common.Hash,
	blockBytes []byte,
	autoRollback bool ) *TxStoreBlockCase {
		return &TxStoreBlockCase{
			blockHash,
			blockBytes,
			autoRollback,
		}
}

func (t *TxStoreBlockCase) execute( ctx DbContext ) error {
	err := ctx.db.Update(func(tx database.Tx) error {
		err := tx.StoreBlock( t.blockHash, t.blockBytes )

		if err == nil && t.autoRollback {
			return errors.New( "auto rollback" )
		}

		return err

	} )

	if err != nil {
		fmt.Printf("TxStoreBlockCase error %s\n", err.Error() )
	} else {
		fmt.Printf("TxStoreBlockCase success\n" )
	}

	return err
}

func (t *TxStoreBlockCase) isWritable() bool {
	return true
}


// HasBlock
type TxHasBlockCase struct {
	blockHash *common.Hash

	result bool
}

func NewHasBlockCase(
	blockHash *common.Hash) *TxHasBlockCase {
	return &TxHasBlockCase{
		blockHash,
		false,
	}
}

func (t *TxHasBlockCase) execute( ctx DbContext ) error {
	err := ctx.db.View(func(tx database.Tx) error {
		exist, err := tx.HasBlock( t.blockHash )
		t.result = exist
		return err
	} )

	if err != nil {
		fmt.Printf("TxHasBlockCase error %s\n", err.Error() )
	} else {
		fmt.Printf("TxHasBlockCase success\n" )
	}

	return err
}

func (t *TxHasBlockCase) isWritable() bool {
	return false
}

// FetchBlockHeader
type TxFetchBlockHeaderCase struct {
	blockHash *common.Hash

	result []byte
}

func NewFetchBlockHeaderCase(
	blockHash *common.Hash ) *TxFetchBlockHeaderCase {
	return &TxFetchBlockHeaderCase{
		blockHash,
		nil,
	}
}

func (t *TxFetchBlockHeaderCase) execute( ctx DbContext ) error {
	err := ctx.db.View(func(tx database.Tx) error {
		headerData, err := tx.FetchBlockHeader( t.blockHash )
		t.result = headerData
		return err
	} )

	if err != nil {
		fmt.Printf("TxFetchBlockHeaderCase error %s\n", err.Error() )
	} else {
		fmt.Printf("TxFetchBlockHeaderCase success\n" )
	}

	return err
}

func (t *TxFetchBlockHeaderCase) isWritable() bool {
	return false
}

// FetchBlock
type TxFetchBlockCase struct {
	blockHash *common.Hash

	result []byte
}

func NewFetchBlockCase(
	blockHash *common.Hash ) *TxFetchBlockCase {
	return &TxFetchBlockCase{
		blockHash:blockHash,
	}
}

func (t *TxFetchBlockCase) execute( ctx DbContext ) error {
	err := ctx.db.View(func(tx database.Tx) error {
		blockData, err := tx.FetchBlock( t.blockHash )
		t.result = blockData
		return err
	} )

	if err != nil {
		fmt.Printf("TxFetchBlockCase error %s\n", err.Error() )
	} else {
		fmt.Printf("TxFetchBlockCase success %v\n", t.result[:10] )
	}

	return err
}

func (t *TxFetchBlockCase) isWritable() bool {
	return false
}

// FetchBlockRegion
type TxFetchBlockRegionCase struct {
	hash *common.Hash
	offset uint32
	len uint32

	result []byte
}

func NewFetchBlockRegionCase(
	hash *common.Hash,
	offset, len uint32 ) *TxFetchBlockRegionCase {
		return &TxFetchBlockRegionCase{
			hash:hash,
			offset:offset,
			len:len,
		}
}

func (t *TxFetchBlockRegionCase) execute( ctx DbContext ) error {
	region := database.BlockRegion{
		t.hash,
		t.offset,
		t.len,
	}

	err := ctx.db.View(func(tx database.Tx) error {
		regionData, err := tx.FetchBlockRegion( &region )
		t.result = regionData
		return err
	} )

	if err != nil {
		fmt.Printf("TxFetchBlockRegionCase error %s\n", err.Error() )
	} else {
		fmt.Printf("TxFetchBlockRegionCase success\n" )
	}

	return err
}

func (t *TxFetchBlockRegionCase) isWritable() bool {
	return false
}


// TxRollback not used
type TxRollbackCase struct {

}

func NewTxRollbackCase() *TxRollbackCase {
	return &TxRollbackCase{}
}

func (t *TxRollbackCase) execute( ctx DbContext ) error {
	err := ctx.tx.Rollback()

	if err != nil {
		fmt.Printf("TxRollbackCase error %s\n", err.Error() )
	} else {
		fmt.Printf("TxRollbackCase success\n" )
	}

	return err
}

func (t *TxRollbackCase) isWritable() bool {
	return true
}

// TxCommit not used
type TxCommitCase struct {

}

func NewTxCommitCase() *TxCommitCase {
	return &TxCommitCase{}
}

func (t *TxCommitCase) execute( ctx DbContext ) error {
	err := ctx.tx.Commit()

	if err != nil {
		fmt.Printf("TxCommitCase error %s\n", err.Error() )
	} else {
		fmt.Printf("TxCommitCase success\n" )
	}

	return err
}

func (t *TxCommitCase) isWritable() bool {
	return true
}


// db put
type DbPutCase struct {
	key []byte
	val []byte
}

func NewDbPutCase(
	key, val []byte ) *DbPutCase {
	return &DbPutCase{
		key:key,
		val:val,
	}
}

func (t *DbPutCase) execute( ctx DbContext ) error {
	err := ctx.db.Put( t.key, t.val )
	if err != nil {
		fmt.Printf("DbPutCase error %s\n", err.Error() )
	} else {
		fmt.Printf("DbPutCase success\n" )
	}

	return err
}

func (t *DbPutCase) isWritable() bool {
	return true
}


// db query
type DbQueryCase struct {
	key []byte

	result []byte
}

func NewDbQueryCase(
	key []byte ) *DbQueryCase {
	return &DbQueryCase{
		key:key,
	}
}

func (t *DbQueryCase) execute( ctx DbContext ) error {
	_, err := ctx.db.Has( t.key )
	if err != nil {
		fmt.Printf("DbQueryCase has error %s\n", err.Error() )
	} else {
		fmt.Printf("DbQueryCase has success\n" )
	}

	data, err := ctx.db.Get( t.key )
	t.result = data

	if err != nil {
		fmt.Printf("DbQueryCase get error %s\n", err.Error() )
	} else {
		fmt.Printf("DbQueryCase get success\n" )
	}

	return err
}

func (t *DbQueryCase) isWritable() bool {
	return false
}

// db delete
type DbDeleteCase struct {
	key []byte
}

func NewDbDeleteCase(
	key []byte ) *DbDeleteCase {
	return &DbDeleteCase{
		key,
	}
}

func (t *DbDeleteCase) execute( ctx DbContext ) error {
	err := ctx.db.Delete( t.key )

	if err != nil {
		fmt.Printf( "DbDeleteCase error %s\n", err.Error() )
	} else {
		fmt.Printf( "DbDeleteCase success\n" )
	}

	return err
}

func (t *DbDeleteCase) isWritable() bool {
	return true
}

// db batch put
type DbBatchPutCase struct {
	key []byte
	val []byte
}

func NewBatchPutCase( key, val []byte ) *DbBatchPutCase {

	return &DbBatchPutCase{
		key,
		val,
	}
}

func (t *DbBatchPutCase) execute( ctx DbContext ) error {
	err := ctx.batch.Put( t.key, t.val )

	if err != nil {
		fmt.Printf("DbBatchPutCase error %s\n", err.Error() )
	} else {
		fmt.Printf("DbBatchPutCase success\n" )
	}

	return err
}

func (t *DbBatchPutCase) isWritable() bool {
	return true
}

// db batch delete
type DbBatchDeleteCase struct {
	key []byte
}

func NewDbBatchDeleteCase( key []byte ) *DbBatchDeleteCase {

	return &DbBatchDeleteCase{
		key:key,
	}
}

func (t *DbBatchDeleteCase) execute( ctx DbContext ) error {
	err := ctx.batch.Delete( t.key )

	if err != nil {
		fmt.Printf("DbBatchDeleteCase error %s\n", err.Error() )
	} else {
		fmt.Printf("DbBatchDeleteCase success\n" )
	}

	return err
}

func (t *DbBatchDeleteCase) isWritable() bool {
	return true
}

// db Batch write
type DbBatchWriteCase struct {
}

func NewDbBatchWriteCase() *DbBatchWriteCase {
	return &DbBatchWriteCase{}
}

func (t *DbBatchWriteCase) execute( ctx DbContext ) error {
	sz := ctx.batch.ValueSize()
	fmt.Printf( "DbBatchWriteCase execute value size is %d\n", sz )
	return ctx.batch.Write()
}

func (t *DbBatchWriteCase) isWritable() bool {
	return true
}

// db Batch reset
type DbBatchResetCase struct {

}

func NewDbBatchResetCase() *DbBatchResetCase {
	return &DbBatchResetCase{}
}

func (t *DbBatchResetCase) execute( ctx DbContext ) error {
	ctx.batch.Reset()
	fmt.Printf( "DbBatchResetCase execute\n" )
	return nil
}

func (t *DbBatchResetCase) isWritable() bool {
	return true
}

// bucket put
type DbBucketPutCase struct {
	key []byte
	val []byte
}

func NewDbBucketPutCase( key, val []byte ) *DbBucketPutCase {
	return &DbBucketPutCase{
		key:key,
		val:val,
	}
}

func (t *DbBucketPutCase) execute( ctx DbContext ) error {
	err := ctx.bucket.Put( t.key, t.val )

	if err != nil {
		fmt.Printf("DbBucketPutCase error %s\n", err.Error() )
	} else {
		fmt.Printf("DbBucketPutCase success\n" )
	}

	return err
}

func (t *DbBucketPutCase) isWritable() bool {
	return true
}


// bucket delete
type DbBucketDeleteCase struct {
	key []byte
}

func NewDbBucketDeleteCase( key []byte ) *DbBucketDeleteCase {
	return &DbBucketDeleteCase{
		key: key,
	}
}

func (t *DbBucketDeleteCase) execute( ctx DbContext ) error {
	err := ctx.bucket.Delete( t.key )

	if err != nil {
		fmt.Printf("DbBucketDeleteCase error %s\n", err.Error() )
	} else {
		fmt.Printf("DbBucketDeleteCase success\n" )
	}

	return err
}

func (t *DbBucketDeleteCase) isWritable() bool {
	return true
}

// bucket get
type DbBucketGetCase struct {
	key []byte
	result []byte
}

func NewDbBucketGetCase( key []byte ) *DbBucketGetCase {
	return &DbBucketGetCase{
		key:key,
	}
}

func (t *DbBucketGetCase) execute( ctx DbContext ) error {
	t.result = ctx.bucket.Get( t.key )

	fmt.Printf( "DbBucketGetCase execute\n" )
	return nil
}

func (t *DbBucketGetCase) isWritable() bool {
	return false
}

// create bucket
type DbBucketCreateCase struct {
	key []byte
	result database.Bucket
}

func NewDbBucketCreateCase( key []byte ) *DbBucketCreateCase {
	return &DbBucketCreateCase{
		key:key,
	}
}

func (t *DbBucketCreateCase) execute( ctx DbContext ) error {
	bucket, err := ctx.bucket.CreateBucket( t.key )
	t.result = bucket

	if err != nil {
		fmt.Printf("DbBucketCreateCase error %s\n", err.Error() )
	} else {
		fmt.Printf("DbBucketCreateCase success\n" )
	}

	return err
}

func (t *DbBucketCreateCase) isWritable() bool {
	return true
}

// create bucket if not exists
type DbBucketCreateIfNotExistsCase struct {
	key []byte
	result database.Bucket
}

func NewDbBucketCreateIfNotExistsCase( key []byte ) *DbBucketCreateIfNotExistsCase {
	return &DbBucketCreateIfNotExistsCase{
		key:key,
	}
}

func (t *DbBucketCreateIfNotExistsCase) execute( ctx DbContext ) error {
	bucket, err := ctx.bucket.CreateBucketIfNotExists( t.key )
	t.result = bucket

	if err != nil {
		fmt.Printf("DbBucketCreateIfNotExistsCase error %s\n", err.Error() )
	} else {
		fmt.Printf("DbBucketCreateIfNotExistsCase success\n" )
	}
	return err
}

func (t *DbBucketCreateIfNotExistsCase) isWritable() bool {
	return true
}

// delete buckdet
type DbBucketDeleteBucketCase struct {
	key []byte
}

func NewDbBucketDeleteBucketCase( key []byte ) *DbBucketDeleteBucketCase {
	return &DbBucketDeleteBucketCase{
		key:key,
	}
}

func (t *DbBucketDeleteBucketCase) execute( ctx DbContext ) error {
	err := ctx.bucket.DeleteBucket( t.key )

	if err != nil {
		fmt.Printf("DbBucketDeleteBucketCase error %s\n", err.Error() )
	} else {
		fmt.Printf("DbBucketDeleteBucketCase success\n" )
	}

	return err
}

func (t *DbBucketDeleteBucketCase) isWritable() bool {
	return true
}

// bucket Foreach
type DbBucketForEachCase struct {
}

func NewDbBucketForEachCase() *DbBucketForEachCase {
	return &DbBucketForEachCase{}
}

func (t *DbBucketForEachCase) execute( ctx DbContext ) error {
	fmt.Println("DbBucketForEachCase list data:" )
	err := ctx.bucket.ForEach(func(k, v []byte) error {
		fmt.Printf( "	key=%v, val=%v\n", k, v )
		return nil
	} )

	if err != nil {
		fmt.Printf("DbBucketForEachCase error %s\n", err.Error() )
	} else {
		fmt.Printf("DbBucketForEachCase success\n" )
	}

	return err
}

func (t *DbBucketForEachCase) isWritable() bool {
	return false
}

// bucket ForEachBucket
type DbBucketForEachBucketCase struct {
}

func NewDbbucketForEachbucketCase() *DbBucketForEachBucketCase {
	return &DbBucketForEachBucketCase{}
}

func (t *DbBucketForEachBucketCase) execute( ctx DbContext ) error {
	fmt.Println("DbBucketForEachBucketCase list data:" )
	err := ctx.bucket.ForEachBucket(func(k []byte) error {
		fmt.Printf( "	key is %v\n", k )
		return nil
	} )

	if err != nil {
		fmt.Printf("DbBucketForEachBucketCase error %s\n", err.Error() )
	} else {
		fmt.Printf("DbBucketForEachBucketCase success\n" )
	}

	return err
}

func  (t *DbBucketForEachBucketCase) isWritable() bool {
	return false
}


// cursor delete
type DbCursorDeleteCase struct {

}

func NewDbCursorDeleteCase() *DbCursorDeleteCase {
	return &DbCursorDeleteCase{}
}

func (t *DbCursorDeleteCase) execute(ctx DbContext) error {
	err := ctx.cursor.Delete()

	if err != nil {
		fmt.Printf("DbCursorDeleteCase error %s\n", err.Error() )
	} else {
		fmt.Printf("DbCursorDeleteCase success\n" )
	}

	return err
}

func (t *DbCursorDeleteCase) isWritable() bool {
	return true
}

// cursor first
type DbCursorFirstCase struct {
	result bool
}

func NewDbCursorFirstCase() *DbCursorFirstCase {
	return &DbCursorFirstCase{}
}

func (t *DbCursorFirstCase) execute( ctx DbContext ) error {
	t.result =  ctx.cursor.First()

	fmt.Printf("DbCursorFirstCase execute %v\n", t.result )

	return nil
}

func (t *DbCursorFirstCase) isWritable() bool {
	return false
}

// cursor last
type DbCursorLastCase struct {
	result bool
}

func NewDbCursorLastCase() *DbCursorLastCase {
	return &DbCursorLastCase{}
}

func (t *DbCursorLastCase) execute( ctx DbContext ) error {
	t.result = ctx.cursor.Last()

	fmt.Printf("DbCursorLastCase execute %v\n", t.result )
	return nil
}

func (t *DbCursorLastCase) isWritable() bool {
	return false
}

// cursor next
type DbCursorNextCase struct {
	result bool
}

func NewDbCursorNextCase() *DbCursorNextCase{
	return &DbCursorNextCase{}
}

func (t *DbCursorNextCase) execute( ctx DbContext ) error {
	t.result = ctx.cursor.Next()
	fmt.Printf("DbCursorNextCase execute %v\n", t.result )
	return nil
}

func (t *DbCursorNextCase) isWritable() bool {
	return false
}

// cursor prev
type DbCursrPrevCase struct {
	result bool
}

func NewDbCursorPrevCase() *DbCursrPrevCase{
	return &DbCursrPrevCase{}
}

func (t *DbCursrPrevCase) execute( ctx DbContext ) error {
	t.result = ctx.cursor.Prev()

	fmt.Printf("DbCursrPrevCase execute %v\n", t.result )
	return nil
}

func (t *DbCursrPrevCase) isWritable() bool {
	return false
}

// cursor seek
type DbCursorSeekCase struct {
	seek []byte

	result bool
}

func NewDbCursorSeekCase( seek []byte ) *DbCursorSeekCase{
	return &DbCursorSeekCase{
		seek:seek,
	}
}

func (t *DbCursorSeekCase) execute( ctx DbContext ) error {
	t.result = ctx.cursor.Seek( t.seek )

	fmt.Printf("DbCursorSeekCase execute %v\n", t.result )

	return nil
}

func (t *DbCursorSeekCase) isWritable() bool {
	return false
}

// cursor key
type DbCursorKeyCase struct {
	result []byte
}

func NewDbCursorKeyCase() *DbCursorKeyCase {
	return &DbCursorKeyCase{}
}

func (t *DbCursorKeyCase) execute( ctx DbContext ) error {
	t.result = ctx.cursor.Key()

	fmt.Printf("DbCursorKeyCase execute %v\n", t.result )
	return nil
}

func (t *DbCursorKeyCase) isWritable() bool {
	return false
}

// cursor value
type DbCursorValueCase struct {
	result []byte
}

func NewDbCursorValueCase() *DbCursorValueCase {
	return &DbCursorValueCase{}
}

func (t *DbCursorValueCase) execute( ctx DbContext ) error {
	t.result = ctx.cursor.Value()

	fmt.Printf("DbCursorValueCase execute %v\n", t.result )

	return nil
}

func (t *DbCursorValueCase) isWritable() bool {
	return false
}
