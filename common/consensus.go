// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package common

const (
	//consensus type
	SOLO           = iota
	POA
	SATOSHIPLUS
	ConsensusCount
)

var (
	consensusByName = make(map[string]int32)
)

var consensusArray = [ConsensusCount]string{
	SOLO: "solo",
	POA:  "poa",
	SATOSHIPLUS:"satoshiplus",
}

// GetConsensus returns consensus type according to the given name
func GetConsensus(name string) int32 {
	if c, ok := consensusByName[name]; ok {
		return c
	}
	return -1
}

func init() {
	for k, v := range consensusArray {
		consensusByName[v] = int32(k)
	}
}
