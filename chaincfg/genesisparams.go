// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package chaincfg

import (
	"encoding/json"
	"github.com/AsimovNetwork/asimov/common"
)

// Params for genesis contracts
var (
	OfficialAddress = common.HexToAddress("0x660000000000000000000000000000000000000000")
	NetConstructorArgsMap = map[string]map[string][]common.Address{
		common.DevelopNet.String(): {
			// GenesisOrganization
			"genesisCitizens": {
				common.HexToAddress("0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a"),
				common.HexToAddress("0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad"),
				common.HexToAddress("0x66a8c1def9e07589415ebf59b55f76b9d5115064c9"),
				common.HexToAddress("0x662306d43258293846ab198810879e31feac759fac"),
				common.HexToAddress("0x6697a7f6e03dcff5cd9cbf99b5a77e7b31dd07153c"),
				common.HexToAddress("0x660832b804c7ef67bcc224de8a48c9e7a3a55b3c1f"),
				common.HexToAddress("0x6661d84abf1832cae49adc891d0331cbbba0d92963"),
			},
			// ValidatorCommittee
			"_validators": {
				common.HexToAddress("0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a"),
				common.HexToAddress("0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad"),
				common.HexToAddress("0x66a8c1def9e07589415ebf59b55f76b9d5115064c9"),
				common.HexToAddress("0x662306d43258293846ab198810879e31feac759fac"),
				common.HexToAddress("0x6697a7f6e03dcff5cd9cbf99b5a77e7b31dd07153c"),
			},
			// ConsensusPOA, used in subchain.
			"_admins": {
				common.HexToAddress("0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a"),
			},
			"_miners": {
				common.HexToAddress("0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad"),
				common.HexToAddress("0x66a8c1def9e07589415ebf59b55f76b9d5115064c9"),
			},
		},
		common.TestNet.String(): {
			// GenesisOrganization
			"genesisCitizens": {
				common.HexToAddress("0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a"),
				common.HexToAddress("0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad"),
				common.HexToAddress("0x66a8c1def9e07589415ebf59b55f76b9d5115064c9"),
				common.HexToAddress("0x662306d43258293846ab198810879e31feac759fac"),
				common.HexToAddress("0x6697a7f6e03dcff5cd9cbf99b5a77e7b31dd07153c"),
				common.HexToAddress("0x660832b804c7ef67bcc224de8a48c9e7a3a55b3c1f"),
				common.HexToAddress("0x6661d84abf1832cae49adc891d0331cbbba0d92963"),
			},
			// ValidatorCommittee
			"_validators": {
				common.HexToAddress("0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a"),
				common.HexToAddress("0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad"),
				common.HexToAddress("0x66a8c1def9e07589415ebf59b55f76b9d5115064c9"),
				common.HexToAddress("0x662306d43258293846ab198810879e31feac759fac"),
				common.HexToAddress("0x6697a7f6e03dcff5cd9cbf99b5a77e7b31dd07153c"),
			},
			// ConsensusPOA, used in subchain.
			"_admins": {
				common.HexToAddress("0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a"),
			},
			"_miners": {
				common.HexToAddress("0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad"),
				common.HexToAddress("0x66a8c1def9e07589415ebf59b55f76b9d5115064c9"),
			},
		},
	}
)

// ContractInfo defines all fields for a contract.
type ContractInfo struct {
	// name of contract
	Name    string
	// an address for some contract instance
	Address []byte
	// compiled code for contract, contains ctor and params.
	Code    string
	// abis for contract
	AbiInfo string
	// init params
	InitCode    string
	// block height where the contract active.
	BlockHeight int32
}

// Decode a bytes array into a map, it is only used for genesis block's data.
func TransferGenesisData(data []byte) map[common.ContractCode][]ContractInfo {
	var cMap map[common.ContractCode][]ContractInfo
	err := json.Unmarshal(data, &cMap)
	if err != nil {
		panic(err)
	}
	return cMap
}
