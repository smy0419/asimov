// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package ainterface

type Consensus interface {
	Start() error
	Halt() error
	GetRoundInterval() int64
}
