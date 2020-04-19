// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package vrf

import (
	"encoding/binary"
	"github.com/AsimovNetwork/asimov/common"
	"math"
	"sort"
)

// Calculate a rand number from hash, and an integer
func getRandSource(hash common.Hash, number int64) uint64 {
	numByte := [8]byte{}
	binary.BigEndian.PutUint64(numByte[:], uint64(number))
	for i := 0; i < 8; i++ {
		hash[i] = hash[i] ^ numByte[i]
	}
	sourceHash := common.HashH(hash[:])
	sourceData := binary.BigEndian.Uint64(sourceHash[:])
	return sourceData
}

// getStakeholdersFromAddress transfers given addresses to struct of stakeholder
func getStakeholdersFromAddress(candidates map[common.Address]float64) []stakeholder {
	stakeHolders := make(stakeholderList, 0, len(candidates))
	for candidate, stake := range candidates {
		stakeHolders = append(stakeHolders, stakeholder{
			candidate, stake * 10000,
		})
	}
	sort.Sort(stakeHolders)

	return stakeHolders
}

// SelectValidators returns validators from given addresses according to specify algorithm
func SelectValidators(candidates map[common.Address]float64, blockHash common.Hash, round uint32, roundSize uint16) []*common.Address {
	stakeholders := getStakeholdersFromAddress(candidates)

	validators := make([]*common.Address, roundSize)
	// Create the Merkle tree
	nodeTree := createMerkleTree(stakeholders)

	// if no coin, fill validators one by one in order
	if math.Abs(nodeTree[0].getCoin()) < 1e-8 {
		stakeholderSize := len(stakeholders)
		for i := 0; i < int(roundSize); i++ {
			validators[i] = &stakeholders[i%stakeholderSize].addr
		}
	} else {
		avgCost := nodeTree[0].getCoin() / float64(roundSize)
		startPos := int64(round) * int64(roundSize)
		for i := 0; i < int(roundSize); i++ {
			merkleRootCoin := nodeTree[0].getCoin()
			if math.Abs(merkleRootCoin) < 1e-8 {
				nodeTree[0].reset()
				merkleRootCoin = nodeTree[0].getCoin()
			}
			// generate rand number
			sourceData := getRandSource(blockHash, startPos + int64(i))
			randData := sourceData % uint64(merkleRootCoin) + 1

			// choose holder and fill into the result.
			holder := selectHolder(nodeTree, float64(randData), avgCost)
			validators[i] = &holder.addr
		}
	}
	return validators
}

// Create a merkle tree and init leaf node.
func createMerkleTree(stakeholders []stakeholder) []node {
	stakeholderSize := len(stakeholders)
	tree := make([]node, stakeholderSize*2 - 1)

	// init leaf
	for i := 0; i < stakeholderSize; i++ {
		tree[stakeholderSize - 1 + i].stakeholder = &stakeholders[i]
	}

	// init node
	for i := 0; i < stakeholderSize - 1; i++ {
		tree[i].setNode(&tree[i*2+1], &tree[i*2+2])
	}
	return tree
}

// Select a stakeholder which the rand parameter located on,
// and cost the value on the node.
func selectHolder(nodeTree []node, rand float64, avgCost float64) *stakeholder {
	i := 0
	for true {
		if nodeTree[i].isLeaf() {
			nodeTree[i].takeCost(avgCost)
			return nodeTree[i].stakeholder
		}
		leftCoin  := nodeTree[i].left.getCoin()
		if rand <= leftCoin {
			i = i*2+1
		} else {
			rand -= leftCoin
			i = i*2+2
		}
	}
	return nil
}
