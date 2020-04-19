// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txo

import (
	"encoding/binary"
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/hexutil"
)

const (
	LockEntryLength = common.HashLength + 8
)

type VoteId [common.HashLength]byte

func NewVoteId(addr common.Address, id uint32) *VoteId {
	var voteId VoteId
	copy(voteId[:], addr[:])
	binary.LittleEndian.PutUint32(voteId[common.HashLength-4:], id)
	return &voteId
}

// LockEntry houses details about a lock record in a lock view
// such as the duration of the block when to be unlock, lock amount,
// lock id, and which sequence lock.
type LockEntry struct {
	Id     VoteId
	Amount int64
}

func (le *LockEntry) Deserialize(serialized []byte) {
	copy(le.Id[:], serialized[:common.HashLength])
	le.Amount = int64(binary.LittleEndian.Uint64(serialized[common.HashLength:]))
}

func (le *LockEntry) Serialize(serialized []byte) {
	copy(serialized[:], le.Id[:])
	binary.LittleEndian.PutUint64(serialized[common.HashLength:], uint64(le.Amount))
}

func (le *LockEntry) Clone() *LockEntry {
	return &LockEntry{
		Id:     le.Id,
		Amount: le.Amount,
	}
}

func (le *LockEntry) toString() string {
	return fmt.Sprintf("{id: \"%s\", amount: %d}", hexutil.Encode(le.Id[:]), le.Amount)
}

type LockItem struct {
	Entries map[VoteId]*LockEntry
}

func NewLockItem() *LockItem {
	var lockItem LockItem
	lockItem.Entries = make(map[VoteId]*LockEntry)
	return &lockItem
}

func (li *LockItem) PutEntry(entry *LockEntry) {
	li.Entries[entry.Id] = entry
}

func (li *LockItem) Bytes() []byte {
	bytes := make([]byte, len(li.Entries)*LockEntryLength)
	offset := 0
	for _, entry := range li.Entries {
		entry.Serialize(bytes[offset : offset+LockEntryLength])
		offset += LockEntryLength
	}
	return bytes
}

func (li *LockItem) Clone() *LockItem {
	lockItem := NewLockItem()
	for k, v := range li.Entries {
		lockItem.Entries[k] = v.Clone()
	}
	return lockItem
}

func (li *LockItem) Merge(other *LockItem) *LockItem {
	for k, v := range other.Entries {
		if entry, ok := li.Entries[k]; ok {
			entry.Amount += v.Amount

		} else {
			li.Entries[k] = v
		}
	}
	return li
}

func (li *LockItem) String() (res string) {
	for _, v := range li.Entries {
		res = res + v.toString()
	}
	return res
}

func (li *LockItem) EntriesList() []*LockEntry {
	res := make([]*LockEntry, 0, len(li.Entries))
	for _, e := range li.Entries {
		res = append(res, e)
	}
	return res
}

// DeserializeLockItem decodes a lock item from the passed serialized byte
// slice into a new UtxoEntry using a format that is suitable for long-term
// storage.  The format is described in detail above.
func DeserializeLockItem(serialized []byte) (*LockItem, error) {
	if len(serialized)%LockEntryLength != 0 {
		return nil, common.AssertError("unexpected end of data after header")
	}
	size := len(serialized) / LockEntryLength
	entries := make([]LockEntry, size)
	lockItem := NewLockItem()
	for i := 0; i < size; i++ {
		entries[i].Deserialize(serialized[i*LockEntryLength : (i+1)*LockEntryLength])
		lockItem.PutEntry(&entries[i])
	}
	return lockItem, nil
}
