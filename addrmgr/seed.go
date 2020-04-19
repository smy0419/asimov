// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package addrmgr

import (
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	mrand "math/rand"
	"strconv"
	"time"

	"github.com/AsimovNetwork/asimov/chaincfg"
	"github.com/AsimovNetwork/asimov/protos"
)

const (
	// These constants are used by the DNS seed code to pick a random last
	// seen time.
	secondsIn3Days int32 = 24 * 60 * 60 * 3
	secondsIn4Days int32 = 24 * 60 * 60 * 4
)

// SeedFromDNS uses DNS seeding to populate the address manager with peers.
func (am *AddrManager) SeedFromDNS(chainParams *chaincfg.Params, reqServices common.ServiceFlag) {

	for _, dnsseed := range chainParams.DNSSeeds {
		var host string
		if !dnsseed.HasFiltering || reqServices == common.SFNodeNetwork {
			host = dnsseed.Host
		} else {
			host = fmt.Sprintf("x%x.%s", uint64(reqServices), dnsseed.Host)
		}

		go func(host string) {
			randSource := mrand.New(mrand.NewSource(time.Now().UnixNano()))

			seedpeers, err := am.nap.Lookup(host)
			if err != nil {
				log.Infof("DNS discovery failed on seed %s: %v", host, err)
				return
			}
			numPeers := len(seedpeers)

			log.Infof("%d addresses found from DNS seed %s", numPeers, host)

			if numPeers == 0 {
				return
			}
			addresses := make([]*protos.NetAddress, len(seedpeers))
			// if this errors then we have *real* problems
			intPort, _ := strconv.Atoi(chainParams.DefaultPort)
			for i, peer := range seedpeers {
				addresses[i] = protos.NewNetAddressTimestamp(
					// seeds with addresses from
					// a time randomly selected between 3
					// and 7 days ago.
					time.Now().Add(-1*time.Second*time.Duration(secondsIn3Days+
						randSource.Int31n(secondsIn4Days))),
					0, peer, uint16(intPort))
			}

			// Uses a lookup of the dns seeder here. This
			// is rather strange since the values looked up by the
			// DNS seed lookups will vary quite a lot.
			// to replicate this behaviour we put all addresses as
			// having come from the first one.
			am.AddAddresses(addresses, addresses[0])
		}(host)
	}
}
