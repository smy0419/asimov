// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package minersync

import (
	"github.com/AsimovNetwork/asimov/logger"
)

// logger is a logger that is initialized with no output filters.  This
// means the package will not perform any logging by default until the caller
// requests it.
var log logger.Logger

// The default amount of logging is none.
func init() {
	log = logger.GetLogger("BTCN")
}
