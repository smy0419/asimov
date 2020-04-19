// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ffldb

// BenchmarkBlockHeader benchmarks how long it takes to load the mainnet genesis
// block header.
//func BenchmarkBlockHeader(b *testing.B) {
//	// Start by creating a new database and populating it with the mainnet
//	// genesis block.
//	dbPath := filepath.Join(os.TempDir(), "ffldb-benchblkhdr")
//	_ = os.RemoveAll(dbPath)
//	db, err := Create("ffldb", dbPath, blockDataNet)
//	if err != nil {
//		b.Fatal(err)
//	}
//	defer os.RemoveAll(dbPath)
//	defer db.Close()
//	err = db.Update(func(tx database.Tx) error {
//		block := vvsutil.NewBlock(chaincfg.MainNetParams.GenesisBlock)
//		blockBytes, _ := block.Bytes()
//		return tx.StoreBlock(block.Hash(), blockBytes )
//	})
//	if err != nil {
//		b.Fatal(err)
//	}
//
//	b.ReportAllocs()
//	b.ResetTimer()
//	err = db.View(func(tx database.Tx) error {
//		blockHash := chaincfg.MainNetParams.GenesisHash
//		for i := 0; i < b.N; i++ {
//			_, err := tx.FetchBlockHeader(blockHash)
//			if err != nil {
//				return err
//			}
//		}
//		return nil
//	})
//	if err != nil {
//		b.Fatal(err)
//	}
//
//	// Don't benchmark teardown.
//	b.StopTimer()
//}
//
//// BenchmarkBlockHeader benchmarks how long it takes to load the mainnet genesis
//// block.
//func BenchmarkBlock(b *testing.B) {
//	// Start by creating a new database and populating it with the mainnet
//	// genesis block.
//	dbPath := filepath.Join(os.TempDir(), "ffldb-benchblk")
//	_ = os.RemoveAll(dbPath)
//	db, err := Create("ffldb", dbPath, blockDataNet)
//	if err != nil {
//		b.Fatal(err)
//	}
//	defer os.RemoveAll(dbPath)
//	defer db.Close()
//	err = db.Update(func(tx database.Tx) error {
//		block := vvsutil.NewBlock(chaincfg.MainNetParams.GenesisBlock)
//		blockBytes, _ := block.Bytes()
//		return tx.StoreBlock(block.Hash(), blockBytes)
//	})
//	if err != nil {
//		b.Fatal(err)
//	}
//
//	b.ReportAllocs()
//	b.ResetTimer()
//	err = db.View(func(tx database.Tx) error {
//		blockHash := chaincfg.MainNetParams.GenesisHash
//		for i := 0; i < b.N; i++ {
//			_, err := tx.FetchBlock(blockHash)
//			if err != nil {
//				return err
//			}
//		}
//		return nil
//	})
//	if err != nil {
//		b.Fatal(err)
//	}
//
//	// Don't benchmark teardown.
//	b.StopTimer()
//}
