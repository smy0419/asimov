// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package minersync

import (
	"bytes"
	"errors"
	"github.com/AsimovNetwork/asimov/chaincfg"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/database"
	"github.com/AsimovNetwork/asimov/database/dbdriver"
	"os"
	"path/filepath"
	"testing"
)

type FakeBtcClient struct {
	minerInfo []GetBitcoinBlockMinerInfoResult
	chainInfo *GetBlockChainInfoResult
}

func NewFakeBtcClient(height, count int32) *FakeBtcClient {
	minerInfo := make([]GetBitcoinBlockMinerInfoResult, 0, count)
	c := &FakeBtcClient{
		minerInfo: minerInfo,
	}
	c.push(height,count)
	return c
}

func (c *FakeBtcClient) Push(count int32) {
	height := c.minerInfo[len(c.minerInfo)-1].Height + 1
	c.push(height, count)
}

func (c *FakeBtcClient) push(height, count int32) {

	for i := int32(0); i < count; i++ {
		c.minerInfo = append(c.minerInfo, GetBitcoinBlockMinerInfoResult{
			Height:  height + i,
			Pool:    "",
			Hash:    "0x123456",
			Address: "11",
		})
	}

	c.chainInfo = &GetBlockChainInfoResult{
		Blocks: height + count - 1,
	}
}

func (c *FakeBtcClient) GetBitcoinMinerInfo(result interface{}, height, count int32) error {
	begin := c.minerInfo[0].Height
	last := c.minerInfo[len(c.minerInfo)-1].Height
	if height < begin || height+count-1 > last {
		return errors.New("miners not find")
	}
	miners := result.(*[]GetBitcoinBlockMinerInfoResult)
	*miners = c.minerInfo[height-begin : height-begin+count]
	return nil
}

func (c *FakeBtcClient) GetBitcoinBlockChainInfo(result interface{}) error {
	chaininfo := result.(*GetBlockChainInfoResult)
	*chaininfo = *c.chainInfo
	return nil
}

func TestBtcPowerInfo(t *testing.T) {
	winfo := MinerInfo{
		height:       1,
		minerName:    "22",
		minerAddress: "addr",
		blockHash:    common.HexToHash("11"),
		time:         1234,
		txs: []*ValidatorTx{
			&ValidatorTx{
				TxId: "111",
				Vin: []string{
					"123",
				},
				OutAddress: common.HexToAddress("333"),
			},
		},
	}
	w := bytes.NewBuffer(make([]byte, 0, winfo.SerializeSize()))
	c := w.Cap()
	err := winfo.Serialize(w)
	if err != nil {
		t.Errorf("Serialize error %v", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	var rinfo MinerInfo
	err = rinfo.Deserialize(r)
	if err != nil {
		t.Errorf("Deserialize error %v", err)
	}

	if w.Cap() != c {
		t.Errorf("SerializeSize error")
	}

	if rinfo.minerName != winfo.minerName ||
		rinfo.blockHash != winfo.blockHash ||
		rinfo.time != winfo.time ||
		rinfo.minerAddress != winfo.minerAddress ||
		len(rinfo.txs) != len(winfo.txs) ||
		rinfo.txs[0].TxId != winfo.txs[0].TxId ||
		len(rinfo.txs[0].Vin) != len(winfo.txs[0].Vin) ||
		rinfo.txs[0].Vin[0] != winfo.txs[0].Vin[0] ||
		rinfo.txs[0].OutAddress != winfo.txs[0].OutAddress {
		t.Error("Serialize or Deserialize error")
	}
}

func TestUpdateMinerInfo(t *testing.T) {
	dbPath := chaincfg.DefaultDataDir + "/devNetUnitTest"
	db, _ := loadBlockDB(dbPath)
	if db == nil {
		t.Errorf("loadBlockDB error ")
	}

	defer os.RemoveAll(dbPath)
	defer db.Close()

	tests := []struct {
		getCount  int32
		wantCount int32
		wantErr   bool
	}{
		{
			0,
			0,
			false,
		},
		{
			1,
			1,
			false,
		},
		{
			maxUpdateMiner - maxBtcRollback,
			maxUpdateMiner - maxBtcRollback,
			false,
		},
		{
			maxUpdateMiner,
			maxUpdateMiner - maxBtcRollback,
			false,
		},
		{
			maxUpdateMiner * 2,
			maxUpdateMiner - maxBtcRollback,
			false,
		},
	}

	for i, v := range tests {
		test_begin_height := int32(6100) + int32(i)*maxUpdateMiner*2
		client :=NewFakeBtcClient(test_begin_height, maxBtcRollback+ 1)
		bpc := &BtcPowerCollector{
			client:        client,
			db:            db,
			initBtcHeight: test_begin_height,
		}

		infos, err := bpc.Init(0, 1)
		if err != nil {
			t.Errorf("init error %v", err)
			return
		}
		client.Push(maxUpdateMiner - maxBtcRollback)

		if len(infos) != 1 || infos[0].height != test_begin_height {
			t.Errorf("init error %v", err)
		}

		_, count, err := bpc.UpdateMinerInfo(v.getCount)
		if count != v.wantCount || v.wantErr != (err != nil) {
			t.Errorf("test error in index :%v", i)
		}

	}

}

func TestGetBitcoinMiner(t *testing.T) {
	dbPath := chaincfg.DefaultDataDir + "/devNetUnitTest"
	db, _ := loadBlockDB(dbPath)
	if db == nil {
		t.Errorf("loadBlockDB error ")
	}

	defer os.RemoveAll(dbPath)
	defer db.Close()

	test_begin_height := int32(6100)
	bpc := &BtcPowerCollector{
		client:        NewFakeBtcClient(test_begin_height, maxUpdateMiner*2),
		db:            db,
		initBtcHeight: test_begin_height,
	}

	infos, err := bpc.Init(0, 144)
	if err != nil {
		t.Errorf("init error %v", err)
		return
	}

	if len(infos) != 144 || infos[0].height != test_begin_height {
		t.Errorf("init error %v", err)
		return
	}

	remove_count := int32(6)
	bpc.minerInfos = bpc.minerInfos[remove_count : len(bpc.minerInfos)-int(remove_count)]
	count := len(bpc.minerInfos)
	begin_height := bpc.minerInfos[0].height - test_begin_height
	last_height := bpc.minerInfos[len(bpc.minerInfos)-1].height - test_begin_height

	tests := []struct {
		beginHeight int32
		count       int32
		err         string
	}{
		//all in memory
		{
			begin_height,
			int32(count),
			"GetBitcoinMiner error when all in memory",
		},
		//low height in disk,high height in memory
		{
			0,
			int32(count) + remove_count,
			"GetBitcoinMiner error when low height in disk and high height in memory",
		},
		//all in disk and less than begin height
		{
			0,
			begin_height,
			"GetBitcoinMiner error when all in disk and less than begin height",
		},
		//low height in memory,high height in disk
		{
			begin_height,
			int32(count) + remove_count,
			"GetBitcoinMiner error when low height in memory and high height in disk",
		},
		//all in disk and greater than last height
		{
			last_height + 1,
			remove_count,
			"GetBitcoinMiner error when all in disk and greater than last height",
		},
		//all in disk and greater than last height
		{
			last_height + 2,
			remove_count - 1,
			"GetBitcoinMiner error when all in disk and greater than last height",
		},
		//get all
		{
			0,
			int32(count) + remove_count*2,
			"GetBitcoinMiner error when get all",
		},
	}

	for _, v := range tests {
		infos, err = bpc.GetBitcoinMiner(v.beginHeight, v.count)
		if err != nil {
			t.Errorf("%s %v", v.err, err)
			return
		}

		if len(infos) != int(v.count) ||
			infos[0].height != test_begin_height+v.beginHeight ||
			infos[len(infos)-1].height != test_begin_height+v.beginHeight+v.count-1 {
			t.Error(v.err)
			return
		}
	}

}

func blockDbPath(dataDir, dbType string) string {
	// The database name is based on the database type.
	dbName := "block" + "_" + dbType
	if dbType == "sqlite" {
		dbName = dbName + ".db"
	}
	dbPath := filepath.Join(dataDir, dbName)
	return dbPath
}

func loadBlockDB(dataDir string) (database.Database, error) {
	dbPath := blockDbPath(dataDir, "ffldb")

	db, err := dbdriver.Open("ffldb", dbPath, chaincfg.TestNetParams.Net)
	if err != nil {
		// Return the error if it's not because the database doesn't
		// exist.
		if dbErr, ok := err.(database.Error); !ok || dbErr.ErrorCode !=
			database.ErrDbDoesNotExist {
			return nil, err
		}

		// Create the db if it does not exist.
		err = os.MkdirAll(dataDir, 0700)
		if err != nil {
			return nil, err
		}

		db, err = dbdriver.Create("ffldb", dbPath, chaincfg.TestNetParams.Net)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
