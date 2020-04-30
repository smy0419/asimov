// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mining

import (
	"container/heap"
	"github.com/AsimovNetwork/asimov/ainterface"
	"github.com/AsimovNetwork/asimov/asiutil"
	"github.com/AsimovNetwork/asimov/blockchain"
	"github.com/AsimovNetwork/asimov/chaincfg"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/address"
	"github.com/AsimovNetwork/asimov/crypto"
	"github.com/AsimovNetwork/asimov/protos"
	"math"
	"math/rand"
	"testing"
)

// TestTxPriceHeap ensures the priority queue for transaction fees and
// priorities works as expected.
func TestTxPriceHeap(t *testing.T) {
	// Create some fake priority items that exercise the expected sort
	// edge conditions.
	testItems := []*TxPrioItem{
		{gasPrice: 5678,},
		{gasPrice: 1234,},
		{gasPrice: 10001,},
		{gasPrice: 0,},
	}

	// Add random data in addition to the edge conditions already manually
	// specified.
	randSeed := rand.Int63()
	defer func() {
		if t.Failed() {
			t.Logf("Random numbers using seed: %v", randSeed)
		}
	}()
	prng := rand.New(rand.NewSource(randSeed))
	for i := 0; i < 1000; i++ {
		testItems = append(testItems, &TxPrioItem{
			gasPrice: prng.Float64() * 10000,
		})
	}

	// Test sorting by fee per KB then priority.
	var highest *TxPrioItem
	priorityQueue := NewTxPriorityQueue(len(testItems))
	for i := 0; i < len(testItems); i++ {
		prioItem := testItems[i]
		if highest == nil {
			highest = prioItem
		}
		if prioItem.gasPrice >= highest.gasPrice {

			highest = prioItem
		}
		heap.Push(priorityQueue, prioItem)
	}

	for i := 0; i < len(testItems); i++ {
		prioItem := heap.Pop(priorityQueue).(*TxPrioItem)
		if prioItem.gasPrice > highest.gasPrice {
			t.Fatalf("fee sort: item (fee per KB: %v) higher than than prev "+
				"(fee per KB: %v)", prioItem.gasPrice, highest.gasPrice, )
		}
		highest = prioItem
	}
}

func TestCreateCoinbaseTx(t *testing.T) {
	privKey, _ := crypto.NewPrivateKey(crypto.S256())
	pkaddr, _ := address.NewAddressPubKey(privKey.PubKey().SerializeCompressed())
	addr := pkaddr.AddressPubKeyHash()
	tests := []struct {
		validater common.IAddress
		height    int32
		wantErr   bool
	}{
		{
			pkaddr,
			1,
			false,
		}, {
			addr,
			1,
			false,
		}, {
			&common.Address{},
			1,
			true,
		}, {
			nil,
			1,
			false,
		}, {
			nil,
			0,
			false,
		}, {
			nil,
			math.MaxInt32,
			false,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		coinbaseScript, err := StandardCoinbaseScript(test.height, 0)
		_, _, err = CreateCoinbaseTx(&chaincfg.DevelopNetParams, coinbaseScript, test.height, test.validater, nil)
		if test.wantErr != (err != nil) {
			t.Errorf("tests #%d error %v", i, err)
		}
	}
}

func TestNewBlockTemplate(t *testing.T) {
	policy := Policy{
		BlockProductedTimeOut: DefaultBlockProductedTimeOut,
		TxConnectTimeOut: DefaultTxConnectTimeOut,
		UtxoValidateTimeOut: UtxoValidateTimeOut,
	}
	chain, teardownFunc, err := newFakeChain(&chaincfg.DevelopNetParams)
	if err != nil {
		t.Error("newFakeChain error: ", err)
		return
	}
	fakeTxSource := &fakeTxSource{make(map[common.Hash]*TxDesc)}
	fakeSigSource := &fakeSigSource{make([]*asiutil.BlockSign, 0)}

	g := NewBlkTmplGenerator(
		&policy,
		fakeTxSource,
		fakeSigSource,
		chain,
	)

	defer teardownFunc()

	global_view := blockchain.NewUtxoViewpoint()

	g.FetchUtxoView = func(tx *asiutil.Tx, dolock bool) (viewpoint ainterface.IUtxoViewpoint, e error) {
		neededSet := make(map[protos.OutPoint]struct{})
		prevOut := protos.OutPoint{Hash: *tx.Hash()}
		for txOutIdx := range tx.MsgTx().TxOut {
			prevOut.Index = uint32(txOutIdx)
			neededSet[prevOut] = struct{}{}
		}
		if !blockchain.IsCoinBase(tx) {
			for _, txIn := range tx.MsgTx().TxIn {
				neededSet[txIn.PreviousOutPoint] = struct{}{}
			}
		}

		// Request the utxos from the point of view of the end of the main
		// chain.
		view := blockchain.NewUtxoViewpoint()
		for k, _ := range neededSet {
			view.AddEntry(k,global_view.LookupEntry(k))
		}
		return view, nil
	}

	invaildAsset := protos.NewAssets(0, 0, 1)

	keys := []*crypto.Account{}
	for i := 0; i < 16; i++ {
		privKey, _ := crypto.NewPrivateKey(crypto.S256())
		pkaddr, _ := address.NewAddressPubKey(privKey.PubKey().SerializeCompressed())
		addr := pkaddr.AddressPubKeyHash()
		keys = append(keys, &crypto.Account {*privKey, *privKey.PubKey(),addr})
	}

	fakeTxs := TxDescList{
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1e8, &asiutil.FlowCoinAsset, 0, false, 0, common.HexToHash("1"),
			},
		}, []*fakeOut{
			{
				keys[1].Address, 1e8 - 1e4, &asiutil.FlowCoinAsset,
			},
		}, global_view), GasPrice: 1},
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1e8, &asiutil.FlowCoinAsset, 1, false, 0, common.HexToHash("1"),
			},
		}, []*fakeOut{
			{
				keys[1].Address, 1e8, &asiutil.FlowCoinAsset,
			},
		}, global_view), GasPrice: 2},
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1e18, &asiutil.FlowCoinAsset, 0, false, 0, common.HexToHash("2"),
			}, {
				keys[0], 1e4, &asiutil.FlowCoinAsset, 1, false, 0, common.HexToHash("3"),
			},
		}, []*fakeOut{
			{
				keys[0].Address, 1e18 - 1e12, &asiutil.FlowCoinAsset,
			},
		}, global_view), GasPrice: 3},
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1234567890, &asiutil.FlowCoinAsset, 3, false, 0, common.HexToHash("4"),
			}, {
				keys[1], 1e6, &asiutil.FlowCoinAsset, 5, false, 0, common.HexToHash("4"),
			}, {
				keys[3], 1e4, &asiutil.FlowCoinAsset, 8, false, 0, common.HexToHash("5"),
			},
		}, []*fakeOut{
			{
				keys[2].Address, 1234567890 + 1e6, &asiutil.FlowCoinAsset,
			},
		}, global_view), GasPrice: 4},
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1e4, &asiutil.FlowCoinAsset, 0, false, 0, common.HexToHash("5"),
			},
		}, []*fakeOut{
			{
				keys[0].Address, 1e3, &asiutil.FlowCoinAsset,
			}, {
				keys[2].Address, 1e3, &asiutil.FlowCoinAsset,
			}, {
				keys[1].Address, 1e3, &asiutil.FlowCoinAsset,
			}, {
				keys[1].Address, 1e3, &asiutil.FlowCoinAsset,
			}, {
				keys[0].Address, 6e3 - 1, &asiutil.FlowCoinAsset,
			},
		}, global_view), GasPrice: 5},
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1, &asiutil.FlowCoinAsset, 0, false, 0, common.HexToHash("6"),
			}, {
				keys[1], 1e6, &asiutil.FlowCoinAsset, 1, false, 0, common.HexToHash("6"),
			}, {
				keys[2], 1e4, &asiutil.FlowCoinAsset, 2, false, 0, common.HexToHash("6"),
			}, {
				keys[2], 1e4, &asiutil.FlowCoinAsset, 2, false, 0, common.HexToHash("7"),
			}, {
				keys[3], 1e4, &asiutil.FlowCoinAsset, 4, false, 0, common.HexToHash("7"),
			},
		}, []*fakeOut{
			{
				keys[2].Address, 1e6, &asiutil.FlowCoinAsset,
			}, {
				keys[2].Address, 1e4 - 1, &asiutil.FlowCoinAsset,
			}, {
				keys[4].Address, 1e4, &asiutil.FlowCoinAsset,
			}, {
				keys[5].Address, 1e3, &asiutil.FlowCoinAsset,
			}, {
				keys[5].Address, 1e3, &asiutil.FlowCoinAsset,
			}, {
				keys[6].Address, 8e3 - 1, &asiutil.FlowCoinAsset,
			},
		}, global_view), GasPrice: 6},
	}
	//create tx depend last tx
	fakeTxs = append(fakeTxs, &TxDesc{Tx: createFakeTx([]*fakeIn{
		{
			keys[5], 1e3, &asiutil.FlowCoinAsset, 4, false, 0x7FFFFFFF, *fakeTxs[len(fakeTxs)-1].Tx.Hash(),
		},
	}, []*fakeOut{
		{
			keys[0].Address, 1e3 - 2, &asiutil.FlowCoinAsset,
		},
	}, nil), GasPrice: 7})

	invalidFakeTxs := TxDescList{
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1e8, &asiutil.FlowCoinAsset, 0, false, 0, common.HexToHash("1"),
			},
		}, []*fakeOut{
			{
				keys[1].Address, 1e8 - 1e4, &asiutil.FlowCoinAsset,
			},
		}, global_view), GasPrice: 1},
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1e8, &asiutil.FlowCoinAsset, math.MaxUint32, true, 0, common.HexToHash("0"),
			},
		}, []*fakeOut{
			{
				keys[1].Address, 1e8 - 1, &asiutil.FlowCoinAsset,
			},
		}, global_view), GasPrice: 1},
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1e8, &asiutil.FlowCoinAsset, 0, false, 0, common.HexToHash("8"),
			},
		}, []*fakeOut{
			{
				keys[1].Address, 1e8 + 1, &asiutil.FlowCoinAsset,
			},
		}, global_view), GasPrice: 1},
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1e8, &asiutil.FlowCoinAsset, 3, false, 0, common.HexToHash("8"),
			},
		}, []*fakeOut{
			{
				keys[1].Address, 1e8 - 1, invaildAsset,
			},
		}, global_view), GasPrice: 1},
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1e8, invaildAsset, 4, false, 0, common.HexToHash("8"),
			},
		}, []*fakeOut{
			{
				keys[1].Address, 1e8 - 1, invaildAsset,
			},
		}, global_view), GasPrice: 1},
		{Tx: createFakeTx([]*fakeIn{
			{
				keys[0], 1e8, &asiutil.FlowCoinAsset, 5, false, 0, common.HexToHash("8"),
			},
		}, []*fakeOut{
			{
				keys[1].Address, 1e8 - 1, &asiutil.FlowCoinAsset,
			},
		}, nil), GasPrice: 1},
	}

	getFees := func(amounts int64) map[protos.Assets]int64 {
		res := make(map[protos.Assets]int64)
		res[asiutil.FlowCoinAsset] = amounts
		return res
	}

	privateKey := "0xd0f0461b7b4d26cf370e6c73b58ef7fa26e8e30853a8cee901ed42cf0879cb6e"
	account,_ := crypto.NewAccount(privateKey)

	tests := []struct {
		validator   *crypto.Account
		gasFloor    uint64
		gasCeil     uint64
		round       uint32
		slot        uint16
		txs         TxDescList
		wantTx      []*common.Hash
		wantFees    map[protos.Assets]int64
		wantOpCosts []int64
		wantWeight  uint16
		wantErr     bool
	}{
		{
			account, 160000000, 160000000, 1, 0, TxDescList{},
			[]*common.Hash{},
			make(map[protos.Assets]int64),
			[]int64{1}, 120, false,
		}, {
			account, 160000000, 160000000, 1, 0, fakeTxs[0:1],
			[]*common.Hash{fakeTxs[0].Tx.Hash()},
			getFees(1e4),
			[]int64{1, 1}, 120, false,
		}, {
			account, 160000000, 160000000, 1, 0, fakeTxs[1:7],
			[]*common.Hash{fakeTxs[5].Tx.Hash(), fakeTxs[6].Tx.Hash(), fakeTxs[4].Tx.Hash(), fakeTxs[3].Tx.Hash(), fakeTxs[2].Tx.Hash(), fakeTxs[1].Tx.Hash()},
			getFees(1 + 1 + 1e12 + 1e4 + 1 + 1e4 + 3),
			[]int64{1, 6, 1, 5, 1, 1, 1}, 120, false,
		}, {
			account, 160000000, 160000000, 1, 0, invalidFakeTxs,
			[]*common.Hash{},
			make(map[protos.Assets]int64),
			[]int64{1}, 120, false,
		}, {
			keys[0], 160000000, 160000000, 1, 0, TxDescList{},
			[]*common.Hash{},
			make(map[protos.Assets]int64),
			[]int64{1}, 0, true,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {

		fakeTxSource.clear()
		for _, v := range test.txs {
			fakeTxSource.push(v)
		}

		block, err := g.ProduceNewBlock(test.validator, test.gasFloor, test.gasCeil, test.round, test.slot, 5*100000)
		if err != nil {
			if test.wantErr != true {
				t.Errorf("tests #%d error %v", i, err)
			}
			continue
		}

		txs := block.MsgBlock().Transactions

		if block.MsgBlock().Header.CoinBase != *test.validator.Address ||
			block.MsgBlock().Header.Round != test.round ||
			block.MsgBlock().Header.SlotIndex != test.slot ||
			block.MsgBlock().Header.Weight != test.wantWeight {
			t.Errorf("tests #%d Coinbase: %v ,Round: %v ,Slot: %v Weight: %v",
				i, block.MsgBlock().Header.CoinBase, block.MsgBlock().Header.Round, block.MsgBlock().Header.SlotIndex, block.MsgBlock().Header.Weight)
		}

		outTxEqual := func(ltxs []*protos.MsgTx, rtxs []*common.Hash) bool {
			if len(ltxs) != len(rtxs) {
				return false
			}
			for k, v := range ltxs {
				if v.TxHash() != *rtxs[k] {
					return false
				}
			}
			return true
		}
		t.Log(i)
		for _, v := range txs {
			t.Log(v.TxHash())
		}
		t.Log(test.wantTx)

		if !outTxEqual(txs[:len(txs)-1], test.wantTx) {
			t.Errorf("tests #%d out tx error, txlen %d, want tx: %v", i, len(txs), test.wantTx)
		}

		feesEqual := func(outs []*protos.TxOut, r map[protos.Assets]int64) bool {

			for _, out := range outs {
				if out.Assets != asiutil.FlowCoinAsset{
					if out.Value != r[out.Assets]{
						return false
					}
				}
			}
			return true
		}
		coinbase := txs[len(txs)-1]
		if !feesEqual(coinbase.TxOut, test.wantFees) {
			t.Errorf("tests #%d fees error,coinbase out: %v ,want fees: %v", i, coinbase.TxOut, test.wantFees)
		}
	}
}
