// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package vrf

import (
	"bytes"
	"github.com/AsimovNetwork/asimov/common"
)

// stakeholder defines candidate of round and number of coin he owns
type stakeholder struct{
	addr common.Address
	coin float64
}

// stakeholderList defines a list of stakeholder
type stakeholderList []stakeholder

// Len returns number of stakeholders from stakeholderList
func (shl stakeholderList) Len() int { return len(shl) }

func (shl stakeholderList) Less(i, j int) bool {
	holderI := shl[i]
	holderJ := shl[j]
	return holderI.coin > holderJ.coin ||
		(holderI.coin == holderJ.coin && bytes.Compare(holderI.addr[:], holderJ.addr[:]) < 0)
}

func (shl stakeholderList) Swap(i, j int) {
	shl[i], shl[j] = shl[j], shl[i]
}

// node defines relationship between stakeholder and node
type node struct{
	left		*node
	right		*node
	stakeholder	*stakeholder
	cost        float64
}

// setNode inits struct of node
func (n *node) setNode(left *node, right *node) {
	n.left  = left
	n.right = right
	n.cost = 0
}

func (n *node) isLeaf() bool {
	return n.stakeholder != nil
}

func (n *node) hasCoin() bool {
	return n.stakeholder != nil && n.stakeholder.coin > n.cost
}

func (n *node) getCoin() float64 {
	if n.isLeaf() {
		if n.stakeholder.coin <= n.cost {
			return 0
		}
		return n.stakeholder.coin - n.cost
	}
	return n.left.getCoin() + n.right.getCoin()
}

func (n *node) reset() {
	if n.isLeaf() {
		n.cost = 0
	} else {
		n.left.reset()
		n.right.reset()
	}
}

func (n *node) takeCost(avgCost float64)  {
	n.cost += avgCost
}
