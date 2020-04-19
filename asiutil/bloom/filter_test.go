// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package bloom_test

import (
	"testing"

	"github.com/AsimovNetwork/asimov/asiutil/bloom"
	"github.com/AsimovNetwork/asimov/protos"
)

// TestFilterLarge ensures a maximum sized filter can be created.
func TestFilterLarge(t *testing.T) {
	f := bloom.NewFilter(100000000, 0, 0.01, protos.BloomUpdateNone)
	if len(f.MsgFilterLoad().Filter) > protos.MaxFilterLoadFilterSize {
		t.Errorf("TestFilterLarge test failed: %d > %d",
			len(f.MsgFilterLoad().Filter), protos.MaxFilterLoadFilterSize)
	}
}

// TestFilterLoad ensures loading and unloading of a filter pass.
func TestFilterLoad(t *testing.T) {
	merkle := protos.MsgFilterLoad{}

	f := bloom.LoadFilter(&merkle)
	if !f.IsLoaded() {
		t.Errorf("TestFilterLoad IsLoaded test failed: want %v got %v",
			true, !f.IsLoaded())
		return
	}
	f.Unload()
	if f.IsLoaded() {
		t.Errorf("TestFilterLoad IsLoaded test failed: want %v got %v",
			f.IsLoaded(), false)
		return
	}
}

func TestFilterReload(t *testing.T) {
	f := bloom.NewFilter(10, 0, 0.000001, protos.BloomUpdateAll)

	bFilter := bloom.LoadFilter(f.MsgFilterLoad())
	if bFilter.MsgFilterLoad() == nil {
		t.Errorf("TestFilterReload LoadFilter test failed")
		return
	}
	bFilter.Reload(nil)

	if bFilter.MsgFilterLoad() != nil {
		t.Errorf("TestFilterReload Reload test failed")
	}
}
