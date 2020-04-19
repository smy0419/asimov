// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package asiutil

import (
	"errors"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/address"
	"github.com/AsimovNetwork/asimov/common/hexutil"
)

// DecodeAddress decodes the string encoding of an address and returns
// the Address if addr is a valid encoding for a known address type.
//
// The bitcoin network the address is associated with is extracted if possible.
// When the address does not encode the network, such as in the case of a raw
// public key, the address will be associated with the passed defaultNet.
func DecodeAddress(addr string) (common.IAddress, error) {
	addrbytes, err := hexutil.Decode(addr)
	if err != nil {
		return nil, err
	}
	// Serialized public keys are either 65 bytes (130 hex chars) if
	// uncompressed/hybrid or 33 bytes (66 hex chars) if compressed.
	if len(addrbytes) == 65 || len(addr) == 33 {
		return address.NewAddressPubKey(addrbytes[:])
	}

	if len(addrbytes) == common.AddressLength {
		switch addrbytes[0] {
		case common.PubKeyHashAddrID:
			fallthrough
		case common.ScriptHashAddrID:
			fallthrough
		case common.ContractHashAddrID:
			return common.NewAddress(addrbytes)
		default:
			return nil, common.ErrUnknownAddressType
		}
	} else {
		return nil, errors.New("decoded address is of unknown size")
	}
}