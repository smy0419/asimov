package test

import (
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/hexutil"
	"github.com/AsimovNetwork/asimov/crypto"
	"testing"
)

// convert asimov address to ethereum address
func TestAddr(t *testing.T) {
	// 242 37 164 18 63 177 228 117 228 158 235 147 157 24 88 208 28 29 210 108
	// the address customer sees
	customerShowAddr := "ch35v8ZAufPnqAMqCRir6PvRuspr5hN6wf"
	// to bytes
	byteAddr, _ := hexutil.Decode(customerShowAddr)
	fmt.Println(byteAddr)
	// to ethereum address
	ethAddr := common.BytesToAddress(byteAddr).String()
	fmt.Println(ethAddr)
}

func TestTopic(t *testing.T) {
	// name + (argsType1,argsType2)
	topic := "Instructor(string,uint256)"
	topicHash := common.Bytes2Hex(crypto.Keccak256([]byte(topic)))
	fmt.Println("0x" + topicHash)
}
