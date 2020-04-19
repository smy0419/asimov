package rpc

import (
	"github.com/AsimovNetwork/asimov/logger"
)

// logger is a logger that is initialized with no output filters.  This
// means the package will not perform any logging by default until the caller
// requests it.
var rpcLog logger.Logger

// The default amount of logging is none.
func init() {
	rpcLog = logger.GetLogger("RPCS")
}
