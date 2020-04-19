// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package peer

import (
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/logger"
	"github.com/AsimovNetwork/asimov/protos"
	"strings"
	"time"
)

const (
	// maxRejectReasonLen is the maximum length of a sanitized reject reason
	// that will be logged.
	maxRejectReasonLen = 250
)

// logger is a logger that is initialized with no output filters.  This
// means the package will not perform any logging by default until the caller
// requests it.
var log logger.Logger

// The default amount of logging is none.
func init() {
	log = logger.GetLogger("PEER")
}

// LogClosure is a closure that can be printed with %v to be used to
// generate expensive-to-create data for a detailed logger level and avoid doing
// the work if the data isn't printed.
type logClosure func() string

func (c logClosure) String() string {
	return c()
}

func newLogClosure(c func() string) logClosure {
	return logClosure(c)
}

// directionString is a helper function that returns a string that represents
// the direction of a connection (inbound or outbound).
func directionString(inbound bool) string {
	if inbound {
		return "inbound"
	}
	return "outbound"
}

// formatLockTime returns a transaction lock time as a human-readable string.
func formatLockTime(lockTime uint32) string {
	// The lock time field of a transaction is either a block height at
	// which the transaction is finalized or a timestamp depending on if the
	// value is before the lockTimeThreshold.  When it is under the
	// threshold it is a block height.
	if lockTime < common.LockTimeThreshold {
		return fmt.Sprintf("height %d", lockTime)
	}

	return time.Unix(int64(lockTime), 0).String()
}

// invSummary returns an inventory message as a human-readable string.
func invSummary(invList []*protos.InvVect) string {
	// No inventory.
	invLen := len(invList)
	if invLen == 0 {
		return "empty"
	}

	// One inventory item.
	if invLen == 1 {
		iv := invList[0]
		switch iv.Type {
		case protos.InvTypeError:
			return fmt.Sprintf("error %s", iv.Hash)
		case protos.InvTypeBlock:
			return fmt.Sprintf("block %s", iv.Hash)
		case protos.InvTypeTx:
			return fmt.Sprintf("tx %s", iv.Hash)
		case protos.InvTypeSignature:
			return fmt.Sprintf("signature %s", iv.Hash)
		}

		return fmt.Sprintf("unknown (%d) %s", uint32(iv.Type), iv.Hash)
	}

	// More than one inv item.
	return fmt.Sprintf("size %d", invLen)
}

// locatorSummary returns a block locator as a human-readable string.
func locatorSummary(locator []*common.Hash, stopHash *common.Hash) string {
	if len(locator) > 0 {
		return fmt.Sprintf("locator %s, stop %s", locator[0], stopHash)
	}

	return fmt.Sprintf("no locator, stop %s", stopHash)

}

// sanitizeString strips any characters which are even remotely dangerous, such
// as html control characters, from the passed string.  It also limits it to
// the passed maximum size, which can be 0 for unlimited.  When the string is
// limited, it will also add "..." to the string to indicate it was truncated.
func sanitizeString(str string, maxLength uint) string {
	const safeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY" +
		"Z01234567890 .,;_/:?@"

	// Strip any characters not in the safeChars string removed.
	str = strings.Map(func(r rune) rune {
		if strings.ContainsRune(safeChars, r) {
			return r
		}
		return -1
	}, str)

	// Limit the string to the max allowed length.
	if maxLength > 0 && uint(len(str)) > maxLength {
		str = str[:maxLength]
		str = str + "..."
	}
	return str
}

// messageSummary returns a human-readable string which summarizes a message.
// Not all messages have or need a summary.  This is used for debug logging.
func messageSummary(msg protos.Message) string {
	switch msg := msg.(type) {
	case *protos.MsgVersion:
		return fmt.Sprintf("agent %s, pver %d, block %d",
			msg.UserAgent, msg.ProtocolVersion, msg.LastBlock)

	case *protos.MsgVerAck:
		// No summary.

	case *protos.MsgGetAddr:
		// No summary.

	case *protos.MsgAddr:
		return fmt.Sprintf("%d addr", len(msg.AddrList))

	case *protos.MsgPing:
		// No summary - perhaps add nonce.

	case *protos.MsgPong:
		// No summary - perhaps add nonce.

	case *protos.MsgMemPool:
		// No summary.

	case *protos.MsgTx:
		return fmt.Sprintf("hash %s, %d inputs, %d outputs, lock %s",
			msg.TxHash(), len(msg.TxIn), len(msg.TxOut),
			formatLockTime(msg.LockTime))

	case *protos.MsgBlock:
		header := &msg.Header
		return fmt.Sprintf("height %d, hash %s, ver %d, %d tx, %v",
			msg.Header.Height, msg.BlockHash(),
			header.Version, len(msg.Transactions), header.Timestamp)

	case *protos.MsgInv:
		return invSummary(msg.InvList)

	case *protos.MsgNotFound:
		return invSummary(msg.InvList)

	case *protos.MsgGetData:
		return invSummary(msg.InvList)

	case *protos.MsgGetBlocks:
		return locatorSummary(msg.BlockLocatorHashes, &msg.HashStop)

	case *protos.MsgGetHeaders:
		return locatorSummary(msg.BlockLocatorHashes, &msg.HashStop)

	case *protos.MsgHeaders:
		return fmt.Sprintf("num %d", len(msg.Headers))

	case *protos.MsgReject:
		// Ensure the variable length strings don't contain any
		// characters which are even remotely dangerous such as HTML
		// control characters, etc.  Also limit them to sane length for
		// logging.
		rejCommand := sanitizeString(msg.Cmd, common.CommandSize)
		rejReason := sanitizeString(msg.Reason, maxRejectReasonLen)
		summary := fmt.Sprintf("cmd %v, code %v, reason %v", rejCommand,
			msg.Code, rejReason)
		if rejCommand == protos.CmdBlock || rejCommand == protos.CmdTx {
			summary += fmt.Sprintf(", hash %v", msg.Hash)
		}
		return summary
	}

	// No summary for other messages.
	return ""
}
