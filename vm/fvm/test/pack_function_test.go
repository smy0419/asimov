package test

import (
	"fmt"
	"github.com/AsimovNetwork/asimov/vm/fvm"
	"github.com/AsimovNetwork/asimov/common"
	"testing"
)

func TestPackConstructor(t *testing.T) {
	abiStr := `[
	{
		"constant": true,
		"inputs": [
			{
				"name": "id",
				"type": "uint32"
			}
		],
		"name": "testUint32",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "pure",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "say",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"name": "_parent",
				"type": "string"
			},
			{
				"name": "_greeting",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	}
]`

	// address1 := common.HexToAddress("0x948ab52cc7b5107efd4b03a51f0d1688b4a49a54")
	// address2 := common.HexToAddress("0x948ab52cc7b5107efd4b03a51f0d1688b4a49a54")
	// definition, err := abi.JSON(strings.NewReader(abiStr))
	// if err != nil {
	// 	panic(err)
	// }
	// out, err := definition.Pack("", address)
	var args uint32 = 10
	out, err := fvm.PackFunctionArgs(abiStr, "testUint32", args)
	if err != nil {
		panic(err)
	}
	fmt.Println(common.Bytes2Hex(out))
}
