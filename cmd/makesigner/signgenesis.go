// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"crypto/ecdsa"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/hexutil"

	"fmt"
	"github.com/AsimovNetwork/asimov/chaincfg"
	"github.com/AsimovNetwork/asimov/crypto"
	//"github.com/AsimovNetwork/asimov/crypto"
)

func phex(src []byte) string {
	dst := make([]byte, 6 * len(src))
	const hextable = "0123456789abcdef"
	for i, v := range src {
		dst[i*6] = '0'
		dst[i*6+1] = 'x'
		dst[i*6+2] = hextable[v>>4]
		dst[i*6+3] = hextable[v&0x0f]
		dst[i*6+4] = ','
		if i % 8 == 7 {
			dst[i*6+5] = '\n'
		} else {
			dst[i*6+5] = ' '
		}
	}
	return string(dst)
}

func main()  {
	privs := []string{
		"0x7eecba084c2ef8a22ef3d2aa4f7ecf2ce0e9d36b717f2f66f35717ae2806e56c",
		"0x98ca5264f6919fc12536a77c122dfaeb491ab01ed657c6db32e14a252a8125e3",
		"0x476c9b84ae64d263465152e91acb5f3076112f927a9c8dc8852ea67a9990f729",
		"0x10437dafd473a7fc2c01fd205f3a3ab0b02558e85f24a01d129af26ca49822c0",
		"0xf09c331d8b4ee7c1df83d7d5fd7b9c27220157d14d5ba0c4c8c319e3850cd0a7",
		"8oRk3HnMnd2F2qQ7XuZWexaFfGxNP6Poj59BNCBvX8bj",
		"GXXLLWiRk9mEhbhWQNvWRTpXwTVMWx1fZdQu1LYbEyUr",
		"EpbR4wqEwjgpmTWQauh8XWhb8sPnFwCuaufScavHmpW7",
		"9ythTekaMTmuHnH7N5XLR1tu3MM6dn1UMWNrv5u43D4v",
		"ArjBqmxcbP1aLiX7xD627ZMHigYfcdnrYSFrJv9WUt7f",
		"4zjNSFsFKhgto8ziDAmm6681hk5u5RuKdg4Sdszr3hF7",
		"EmNz6woSKQ9aEWaUuo8wsffQXfB6PLKe3wu6DFc9DNuT",
		//"F",
	}
	privArr := []*crypto.PrivateKey{}
	for _, prive := range(privs) {
		bytes, _ := hexutil.Decode(prive)
		priv, pubk := crypto.PrivKeyFromBytes(crypto.S256(), bytes)
		privArr = append(privArr, priv)
		fmt.Println(hexutil.Encode(pubk.SerializeCompressed()))
		fmt.Println(hexutil.Encode(common.Hash160(pubk.SerializeCompressed())))
	}
	params := chaincfg.TestNetParams
	data := params.GenesisBlock.BlockHash()
	fmt.Printf("block hash : {%s}\n", phex(data[:]))
	for _, priv := range(privArr) {
		//signature, err := priv.Sign(data[:])
		signature, err := crypto.Sign(data[:], (*ecdsa.PrivateKey)(priv))
		if err != nil {
			fmt.Println("sign error ", err)
		} else {
			fmt.Printf("sign result : {%s}\n", phex(signature))
		}
	}
}
