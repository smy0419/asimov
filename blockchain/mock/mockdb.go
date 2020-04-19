// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mock

import (
	"github.com/AsimovNetwork/asimov/database"
)

type MockTx struct {
	metadata database.Bucket
}

func (m *MockTx) StoreBlock(blockKey *database.BlockKey, blockBytes []byte) error {
	return nil
}

func (m *MockTx) HasBlock(key *database.BlockKey) (bool, error) {
	return true, nil
}

func (m *MockTx) HasBlocks(keys []database.BlockKey) ([]bool, error) {
	return nil, nil
}

func (m *MockTx) FetchBlockHeader(key *database.BlockKey) ([]byte, error) {
	return nil, nil
}

func (m *MockTx) FetchBlockHeaders(keys []database.BlockKey) ([][]byte, error) {
	return nil, nil
}

func (m *MockTx) FetchBlock(key *database.BlockKey) ([]byte, error) {
	return nil, nil
}

func (m *MockTx) FetchBlocks(keys []database.BlockKey) ([][]byte, error) {
	return nil, nil
}

func (m *MockTx) FetchBlockRegion(region *database.BlockRegion) ([]byte, error) {
	return nil, nil
}

func (m *MockTx) FetchBlockRegions(regions []database.BlockRegion) ([][]byte, error) {
	return nil, nil
}

func (m *MockTx) Metadata() database.Bucket {
	return m.metadata
}

func (m *MockTx) Commit() error {
	return nil
}

func (m *MockTx) Rollback() error {
	return nil
}

var _ database.Tx = (*MockTx)(nil)

func NewMockTx() database.Tx {
	mockTx := &MockTx{}
	mockTx.metadata = NewMockBucket()
	return mockTx
}

type MockBucket struct {
	buckets map[string]database.Bucket
	cache map[string][]byte

}

func (mb *MockBucket) Bucket(key []byte) database.Bucket {
	if bucket, ok := mb.buckets[string(key)]; ok {
		return bucket
	}
	return nil
}

func (mb *MockBucket) CreateBucket(key []byte) (database.Bucket, error) {
	if _, ok := mb.buckets[string(key)]; ok {
		str := "bucket already exists"
		return nil, database.MakeError(database.ErrBucketExists, str, nil)
	}
	newBucket := NewMockBucket()
	mb.buckets[string(key)] = newBucket
	return newBucket, nil
}

func (mb *MockBucket) CreateBucketIfNotExists(key []byte) (database.Bucket, error) {
	return nil, nil
}

func (mb *MockBucket) DeleteBucket(key []byte) error {
	return nil
}

func (mb *MockBucket) ForEach(func(k, v []byte) error) error {
	return nil
}

func (mb *MockBucket) ForEachBucket(func(k []byte) error) error {
	return nil
}

func (mb *MockBucket) Cursor() database.Cursor {
	return nil
}

func (mb *MockBucket) Writable() bool {
	return true
}

func (mb *MockBucket) Put(key, value []byte) error {
	mb.cache[string(key)] = value
	return nil
}

func (mb *MockBucket) Get(key []byte) []byte {
	if value, ok := mb.cache[string(key)]; ok {
		return value
	}
	return nil
}

func (mb *MockBucket) Delete(key []byte) error {
	return nil
}

func NewMockBucket() database.Bucket {
	mb := &MockBucket{}
	mb.buckets = make(map[string]database.Bucket)
	mb.cache = make(map[string][]byte)
	return mb
}

var _ database.Bucket = (*MockBucket)(nil)

type MockCursor struct {
}

func (mc *MockCursor) Bucket() database.Bucket {
	return nil
}

func (mc *MockCursor) Delete() error {
	return nil
}

func (mc *MockCursor) First() bool {
	return true
}

func (mc *MockCursor) Last() bool {
	return true
}

func (mc *MockCursor) Next() bool {
	return true
}

func (mc *MockCursor) Prev() bool {
	return true
}

func (mc *MockCursor) Seek(seek []byte) bool {
	return true
}

func (mc *MockCursor) Key() []byte {
	return nil
}

func (mc *MockCursor) Value() []byte {
	return nil
}

var _ database.Cursor = (*MockCursor)(nil)
