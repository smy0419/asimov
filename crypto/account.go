// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package crypto

import (
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/hexutil"
)

/* crypto object */
type Account struct {
	PrivateKey PrivateKey
	PublicKey  PublicKey
	Address    *common.Address
}

func NewAccount(privateKey string) (*Account, error) {
	privKeyBytes, err := hexutil.Decode(privateKey)
	if err != nil {
		return nil, err
	}

	privKey, publicKey := PrivKeyFromBytes(S256(), privKeyBytes)
	acc := Account {
		PrivateKey: *privKey,
		PublicKey:  *publicKey,
	}
	acc.Address, err = common.NewAddressWithId(common.PubKeyHashAddrID,
		common.Hash160(acc.PublicKey.SerializeCompressed()))
	return &acc, err
}