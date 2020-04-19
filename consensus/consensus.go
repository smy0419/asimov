// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package consensus

import (
	"errors"
	"github.com/AsimovNetwork/asimov/ainterface"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/consensus/params"
	"github.com/AsimovNetwork/asimov/consensus/poa"
	"github.com/AsimovNetwork/asimov/consensus/satoshiplus"
	"github.com/AsimovNetwork/asimov/consensus/solo"
	"github.com/AsimovNetwork/asimov/crypto"
)

// Create a new consensus service via pass config match consensus name
func NewConsensusService(consensusName string, cfg *params.Config) (ainterface.Consensus, error) {
	consensusType := common.GetConsensus(consensusName)
	switch consensusType {
	case common.SOLO:
		return solo.NewSoloService(cfg)
	case common.SATOSHIPLUS:
		return satoshiplus.NewSatoshiPlusService(cfg)
	case common.POA:
		return poa.NewService(cfg)
	}
	return nil, errors.New("unknown consensus " + consensusName)
}

// Create a new round manager under consensus name
func NewRoundManager(consensusName string, acc *crypto.Account) ainterface.IRoundManager {
	consensusType := common.GetConsensus(consensusName)
	switch consensusType {
	case common.SOLO:
		return solo.NewRoundManager([]*common.Address{acc.Address})
	case common.SATOSHIPLUS:
		return satoshiplus.NewRoundManager()
	case common.POA:
		return poa.NewRoundManager()
	}
	return nil
}
