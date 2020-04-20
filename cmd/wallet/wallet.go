// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/address"
	"github.com/AsimovNetwork/asimov/common/bitcoinaddress"
	"github.com/AsimovNetwork/asimov/common/hexutil"
	"github.com/AsimovNetwork/asimov/crypto"
	"github.com/AsimovNetwork/asimov/txscript"
	"os"
	"strconv"
)

func main() {
	cfg, remainArgs, _ := loadConfig()
    if cfg == nil {
		os.Exit(1)
    }
	var decodefunc = func(s string) ([]byte, error) {
		if cfg.Format == "base64" {
			return base64.StdEncoding.DecodeString(s)
		}
		return hexutil.Decode(s)
	}
	var encodefunc = func(b []byte) string {
		if cfg.Format == "base64" {
			return base64.StdEncoding.EncodeToString(b)
		}
		return hexutil.Encode(b)
	}
    if cfg.Cmd == "genKey" || cfg.Cmd == "-gk" {
		genkeys(cfg, decodefunc, encodefunc)
		os.Exit(1)
	}
	if cfg.Cmd == "genMultiSigAddress" || cfg.Cmd == "-gmsa" {
		genMultiSigAddress(cfg, remainArgs, decodefunc, encodefunc)
		os.Exit(1)
	}
}

func genkeys(cfg *config, decodef func(string) ([]byte, error), encodef func(b []byte) string)  {
	privKey, err := crypto.NewPrivateKey(crypto.S256())
	if err != nil {
		fmt.Println("generate key pair err ", err)
		return
	}

	params := &bitcoinaddress.MainNetParams
	if cfg.Net == "dev"{
		params = &bitcoinaddress.DevNetParams
	}else if cfg.Net == "test"{
		params = &bitcoinaddress.TestNet3Params
	}else if cfg.Net == "regtest"{
		params = &bitcoinaddress.RegressionNetParams
	} else if cfg.Net != "main"{
		fmt.Printf("%v net not exist",cfg.Net)
		return
	}

	privKeyBytes := privKey.Serialize()
	pubk := privKey.PubKey().SerializeCompressed()
	pubKeyHash := common.Hash160(pubk)
	addr, _ := common.NewAddressWithId(common.PubKeyHashAddrID,pubKeyHash )

	fmt.Println("New key pair (priv, pubkey) (format:" + cfg.Format + ")")
	fmt.Println("    {", encodef(privKeyBytes), ",", encodef(pubk), "}")
	fmt.Println("    Compressed pubkey hash address:", encodef(addr.ScriptAddress()))

	pkh,_:=bitcoinaddress.NewAddressPubKeyHash(pubKeyHash,params)
	wpkh,_ := bitcoinaddress.NewAddressWitnessPubKeyHash(pubKeyHash,params)
	script,_ := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
	sh,_ := bitcoinaddress.NewAddressScriptHash(script,params)
	fmt.Println("Bitcoin base58check encode secret: ",bitcoinaddress.EncodePrivkey(privKey,params))
	fmt.Println("Bitcoin legacy type address: ",pkh.EncodeAddress())
	fmt.Println("Bitcoin p2sh-segwit type address: ",sh.EncodeAddress())
	fmt.Println("Bitcoin bech32 type address: ",wpkh.EncodeAddress())
}

func genMultiSigAddress(cfg *config, args []string, decodef func(string) ([]byte, error), encodef func(b []byte) string) {
	if len(args) < 3 {
		fmt.Println("arguments invalid count")
		return
	}
	m, err := strconv.Atoi(args[0])
	if err != nil || m < 1 {
		fmt.Println("arguments invalid m")
		return
	}
	n, err := strconv.Atoi(args[1])
	if err != nil || n > 20 || n < m{
		fmt.Println("arguments invalid n")
		return
	}
	if len(args) != n + 2 {
		fmt.Println("arguments invalid length", len(args), n)
		return
	}
	pubkeys := make([]*address.AddressPubKey, 0, n)
	for _, arg := range args[2:] {
		key, err := decodef(arg)
		if err != nil {
			fmt.Println("arguments invalid, failed decode pubkey")
			return
		}
		pubkeyAddr, err := address.NewAddressPubKey(key)
		if err != nil {
			fmt.Println("arguments invalid ")
			return
		}
		pubkeys = append(pubkeys, pubkeyAddr)
	}

	script, err := txscript.MultiSigScript(pubkeys, m)
	if err != nil {
		fmt.Println("Failed to build script")
		return
	}
	scripthash := common.Hash160(script)
	addr, _ := common.NewAddressWithId(common.PubKeyHashAddrID, scripthash)

	fmt.Println("Build multi sig address:")
	fmt.Println("    {\n        Address:", encodef(addr.ScriptAddress()), ",")
	fmt.Println("        Script:", encodef(script))
	fmt.Println("    }")
}