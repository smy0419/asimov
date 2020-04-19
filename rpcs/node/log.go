package node

import (
	"github.com/AsimovNetwork/asimov/logger"
)

// logger is a logger that is initialized with no output filters.  This
// means the package will not perform any logging by default until the caller
// requests it.
var nodeLog logger.Logger

// The default amount of logging is none.
func init() {
	nodeLog = logger.GetLogger("NODE")
}
