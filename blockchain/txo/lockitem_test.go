// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txo

import (
	"encoding/hex"
	"reflect"
	"testing"
	"github.com/AsimovNetwork/asimov/common"
)

func hexToBytes(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic("invalid hex in source file: " + s)
	}
	return b
}

// TestLockItemSerialization ensures serializing and deserializing lock
// item works as expected.
func TestLockItemSerialization(t *testing.T) {
	t.Parallel()

	voteId1 := VoteId{1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,0,0,0,1}
	voteId2 := VoteId{1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,0,0,0,2}
	tests := []struct {
		name       string
		item       *LockItem
		serialized []byte
	}{
		{
			name: "One entry",
			item: &LockItem{
				map[VoteId]*LockEntry{
					voteId1: &LockEntry{
						Id: voteId1,
						Amount:     5000000000,
					},
				},
			},
			serialized: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 1, 0, 242, 5, 42, 1, 0, 0, 0},
		},
		// Adapted from block 100025 in main blockchain.
		{
			name: "Two entries",
			item: &LockItem{
				map[VoteId]*LockEntry{
					voteId1: &LockEntry{
						Id: voteId1,
						Amount:     34405000000,
					},
					voteId2: &LockEntry{
						Id: voteId2,
						Amount:     13761000000,
					},
				},
			},
			serialized: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 1, 64, 163, 178, 2, 8, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 2, 64, 50, 56, 52, 3, 0, 0, 0},
		},
	}

	for i, test := range tests {
		// Deserialize to a lock item.
		gotItem, err := DeserializeLockItem(test.serialized)
		if err != nil {
			t.Errorf("deserializeLockItem #%d (%s) "+
				"unexpected error: %v", i, test.name, err)
			continue
		}

		// Ensure that the deserialized lock item has the
		// correct properties.
		if !reflect.DeepEqual(gotItem, test.item) {
			t.Errorf("deserializeLockItem #%d (%s) "+
				"mismatched entries - got %v, want %v",
				i, test.name, gotItem, test.item)
			continue
		}

		// Ensure the item serializes to the expected value.
		gotItem, err = DeserializeLockItem(test.item.Bytes())
		if !reflect.DeepEqual(gotItem, test.item) {
			t.Errorf("serializeLockItem #%d (%s): "+
				"mismatched bytes - got %v, want %v", i,
				test.name, gotItem, test.item)
			continue
		}
	}
}

// TestLockItemErrors performs negative tests against deserializing lock
// item to ensure error paths work as expected.
func TestLockItemErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		item       *LockItem
		serialized []byte
		errType    error
	}{
		{
			name: "Force assertion due to invalid length",
			serialized: hexToBytes("1234"),
			errType:    common.AssertError(""),
		},
	}

	for _, test := range tests {
		// Ensure the expected error type is returned and the returned
		// slice is nil.
		lockitem, err := DeserializeLockItem(test.serialized)
		if reflect.TypeOf(err) != reflect.TypeOf(test.errType) {
			t.Errorf("deserializeLockItem (%s): expected "+
				"error type does not match - got %T, want %T",
				test.name, err, test.errType)
			continue
		}
		if lockitem != nil {
			t.Errorf("deserializeLockItem (%s): returned "+
				"item is not nil", test.name)
			continue
		}
	}
}

// TestLockItemMerge ensures two lock items merge work as expected.
func TestLockItemMerge(t *testing.T) {
	t.Parallel()
	voteId1 := VoteId{1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,0,0,0,1}
	voteId2 := VoteId{1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,0,0,0,2}
	tests := []struct {
		name  string
		itema *LockItem
		itemb *LockItem
		merge *LockItem
	}{
		{
			name: "New id",
			itema: &LockItem{
				map[VoteId]*LockEntry{
					voteId1:&LockEntry{
						voteId1, 1234567890,
					},
				},
			},
			itemb: &LockItem{
				map[VoteId]*LockEntry{
					voteId2:&LockEntry{
						voteId2, 1000000000,
					},
				},
			},
			merge: &LockItem{
				map[VoteId]*LockEntry{
					voteId1:&LockEntry{
						voteId1, 1234567890,
					},
					voteId2:&LockEntry{
						voteId2, 1000000000,
					},
				},
			},
		},
		{
			name: "Same Id",
			itema: &LockItem{
				map[VoteId]*LockEntry{
					voteId1:&LockEntry{
						voteId1, 1234567890,
					},
				},
			},
			itemb: &LockItem{
				map[VoteId]*LockEntry{
					voteId1:&LockEntry{
						voteId1, 1000000000,
					},
				},
			},
			merge: &LockItem{
				map[VoteId]*LockEntry{
					voteId1:&LockEntry{
						voteId1, 2234567890,
					},
				},
			},
		},
		{
			name: "Same Id while different height, use the big one",
			itema: &LockItem{
				map[VoteId]*LockEntry{
					voteId1:&LockEntry{
						voteId1, 1234567890,
					},
				},
			},
			itemb: &LockItem{
				map[VoteId]*LockEntry{
					voteId1:&LockEntry{
						voteId1, 1000000000,
					},
				},
			},
			merge: &LockItem{
				map[VoteId]*LockEntry{
					voteId1:&LockEntry{
						voteId1, 2234567890,
					},
				},
			},
		},
	}

	for _, test := range tests {
		// Ensure the merge method runs ok.
		lockitem := test.itema.Merge(test.itemb)
		if lockitem != test.itema {
			t.Errorf("lock item Merge (%s): returned item "+
				"is not match",
				test.name)
			continue
		}
		if !reflect.DeepEqual(lockitem, test.merge) {
			t.Errorf("lock item Merge (%s): returned "+
				"item is unexpected -- got (%v), want (%v)",
				test.name, lockitem, test.merge)
			continue
		}
	}
}
