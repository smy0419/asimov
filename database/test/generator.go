// Copyright (c) 2018-2020 The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package test

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/database"
	"github.com/AsimovNetwork/asimov/database/dbdriver"
	"github.com/AsimovNetwork/asimov/database/dbimpl/ethdb"
	"os"
	"path/filepath"
	"strings"
)

const ()

type Generator interface {
	Generate() error
}

// db
type DbGenerator struct {
	dbType  string
	netType common.AsimovNet
	isMemdb bool

	db database.Database
}

func NewDbGenerator(
	dbType string,
	netType common.AsimovNet,
	isMemdb bool) *DbGenerator {

	return &DbGenerator{
		dbType:  dbType,
		netType: netType,
		isMemdb: isMemdb,
	}

}

func (dbg *DbGenerator) Generate() error {
	var db database.Database
	var err error

	if dbg.dbType == database.ETHDB {

		db, err = dbg.generateEthdb()

	} else if dbg.dbType == database.FFLDB {

		db, err = dbg.generateFfldb()

	}

	if db == nil {
		return errors.New(fmt.Sprintf("DbGenerator fail:%s", err.Error()))
	}

	dbg.db = db
	return nil
}

func (dbg *DbGenerator) getDbPath() string {
	baseDataDir := "testdata"
	dbName := "blocks_" + dbg.dbType
	dbPath := filepath.Join(baseDataDir, strings.ToLower(dbg.netType.String()), dbName)
	return dbPath
}

func (dbg *DbGenerator) generateEthdb() (database.Database, error) {
	if dbg.isMemdb {
		return ethdb.NewMemDatabase(), nil
	}

	dbPath := dbg.getDbPath()
	db, err := ethdb.NewLDBDatabase(dbPath, 768, 1024)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (dbg *DbGenerator) generateFfldb() (database.Database, error) {
	if dbg.isMemdb {
		db, err := dbdriver.Create("memdb")
		if err != nil {
			return nil, errors.New("ffldb do not support memdb")
		}

		return db.(database.Database), err
	}

	dbPath := dbg.getDbPath()
	db, err := dbdriver.Open(dbg.dbType, dbPath, dbg.netType)

	if err != nil {
		dbErr, ok := err.(database.Error)
		if !ok || dbErr.ErrorCode != database.ErrDbDoesNotExist {
			return nil, err
		}

		err = os.MkdirAll(dbPath, 0700)
		if err != nil {
			return nil, err
		}
		db, err = dbdriver.Create(dbg.dbType, dbPath, dbg.netType)

		if err != nil {
			return nil, err
		}
	}

	return db.(database.Database), nil
}

// bytes
type BytesGenerator struct {
}

func NewBytesGenerator() *BytesGenerator {
	return &BytesGenerator{}
}

func (bg *BytesGenerator) Generate() error {
	return nil
}

func (bg *BytesGenerator) CustomGenerate(sz uint32) []byte {
	token := make([]byte, sz)
	rand.Read(token)
	return token
}

func (bg *BytesGenerator) GenerateHash() common.Hash {
	data := bg.CustomGenerate(common.HashLength)

	h := common.Hash{}
	h.SetBytes(data)
	return h
}

// block
type BlockGenerator struct {
}

func (bg *BlockGenerator) Generate() error {
	return nil
}

// block tx loc
type BlockTxLocGenerator struct {
}

func (btlg *BlockTxLocGenerator) Generate() error {
	return nil
}

// ffldb path
type FfldbPathGenerator struct {
}

func (fpg *FfldbPathGenerator) Generate() error {
	return nil
}
