package main

import (
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
)

type Rpcrequest struct {
	Id        uint32    `json:"id"`
	Jsonrpc   string    `json:"jsonrpc"`
	Method    string    `json:"method"`
	Params    []int32   `json:"params"`
}


func main() {

	//test vvs-core request getblockchaininfo:
	param := []int32{}
	q := Rpcrequest{
		Id: 1,
		Jsonrpc:"2.0",
		Method:"getblockchaininfo",
		Params:param,
	}
	//post请求提交json数据
	ba, _ := json.Marshal(q)
	res, err := http.Post("http://127.0.0.1:8000/","application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("fakeBitcoin getblockchaininfo result: %s\n", string(body))


	//test vvs-core request listblockminerinfo:
	param = []int32{10,2}
	q = Rpcrequest{
		Id: 1,
		Jsonrpc:"2.0",
		Method:"listblockminerinfo",
		Params:param,
	}
	//post请求提交json数据
	ba, _ = json.Marshal(q)
	res, err = http.Post("http://127.0.0.1:8000/","application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	fmt.Printf("fakeBitcoin listblockminerinfo result: %s\n", string(body))
	return
}