// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"compress/bzip2"
	"fmt"
	"github.com/AsimovNetwork/asimov/asiutil"
	"github.com/AsimovNetwork/asimov/chaincfg"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/database"
	"os"
	"path/filepath"
	"runtime"
)

var (
	asimovdHomeDir     = asiutil.AppDataDir("asimovd", false)

	// blockDataNet is the expected network in the test block data.
	blockDataNet = common.TestNet

	// blockDataFile is the path to a file containing the first 256 blocks
	// of the block chain.
	blockDataFile = filepath.Join("..", "testdata", "blocks1-256.bz2")

	params = &chaincfg.TestNetParams

	// Default global config.
	DataDir =  filepath.Join(asimovdHomeDir, "data", params.Name())
	DbType =  database.FFLDB
)

func realMain() error {
	// Open the file that contains the blocks for reading.
	fi, err := os.Open(DataDir)
	if err != nil {
		fmt.Printf("failed to open file %v, err %v\n", blockDataFile, err)
		return err
	}
	defer func() {
		if err := fi.Close(); err != nil {
			fmt.Printf("failed to close file %v %v\n", blockDataFile, err)
		}
	}()
	bzip2.NewReader(fi)
	return nil
}

func main() {
	// Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	// Work around defer not working after os.Exit()
	if err := realMain(); err != nil {
		os.Exit(1)
	}
}
