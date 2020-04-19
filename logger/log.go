// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package logger

import (
	"fmt"
	"github.com/jrick/logrotate/rotator"
	"os"
	"path/filepath"
	"sort"
)

// logWriter implements an io.Writer that outputs to both standard output and
// the write-end pipe of an initialized logger rotator.
type logWriter struct{}

func (logWriter) Write(p []byte) (n int, err error) {
	os.Stdout.Write(p)
	if logRotator != nil {
		logRotator.Write(p)
	}
	return len(p), nil
}

// Loggers per subsystem.  A single backend logger is created and all subsytem
// loggers created from it will write to the backend.  When adding new
// subsystems, add the subsystem logger variable here and to the
// subsystemLoggers map.
//
// Loggers can not be used before the logger rotator has been initialized with a
// logger file.  This must be performed early during application startup by calling
// initLogRotator.
var (
	// backendLog is the logging backend used to create all subsystem loggers.
	// The backend must not be used before the logger rotator has been initialized,
	// or data races and/or nil pointer dereferences will occur.
	backendLog = NewBackend(logWriter{})

	// logRotator is one of the logging outputs.  It should be closed on
	// application shutdown.
	logRotator *rotator.Rotator

	adxrLog  = backendLog.Logger("ADXR")
	amgrLog  = backendLog.Logger("AMGR")
	mainLog  = backendLog.Logger("ASIM")
	btcnLog  = backendLog.Logger("BTCN")
	chanLog  = backendLog.Logger("CHAN")
	cmgrLog  = backendLog.Logger("CMGR")
	discLog  = backendLog.Logger("DISC")
	ffdbLog  = backendLog.Logger("FFDB")
	indxLog  = backendLog.Logger("INDX")
	minrLog  = backendLog.Logger("MINR")
	peerLog  = backendLog.Logger("PEER")
	rpcsLog  = backendLog.Logger("RPCS")
	scrpLog  = backendLog.Logger("SCRP")
	srvrLog  = backendLog.Logger("SRVR")
	syncLog  = backendLog.Logger("SYNC")
	txmpLog  = backendLog.Logger("TXMP")
	consLog  = backendLog.Logger("CONS")
	solclog  = backendLog.Logger("SOLC")
	contLog  = backendLog.Logger("CONT")
	nodeLog  = backendLog.Logger("NODE")
	contractLog = backendLog.Logger("CNTR")
)

// Initialize package-global logger variables.
func init() {
}

// subsystemLoggers maps each subsystem identifier to its associated logger.
var subsystemLoggers = map[string]Logger{
	"ADXR":     adxrLog,
	"AMGR":     amgrLog,
	"ASIM":     mainLog,
	"BTCN":     btcnLog,
	"CMGR":     cmgrLog,
	"FFDB":     ffdbLog,
	"CHAN":     chanLog,
	"DISC":     discLog,
	"INDX":     indxLog,
	"MINR":     minrLog,
	"PEER":     peerLog,
	"RPCS":     rpcsLog,
	"SCRP":     scrpLog,
	"SRVR":     srvrLog,
	"SYNC":     syncLog,
	"TXMP":     txmpLog,
	"CONS":     consLog,
	"SOLC":     solclog,
	"CONT":     contLog,
	"NODE":     nodeLog,
	"CNTR":     contractLog,
}

func GetLog() Logger {
	return mainLog
}

// initLogRotator initializes the logging rotater to write logs to logFile and
// create roll files in the same directory.  It must be called before the
// package-global log rotater variables are used.
func InitLogRotator(logFile string) {
	logDir, _ := filepath.Split(logFile)
	err := os.MkdirAll(logDir, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger directory: %v\n", err)
		os.Exit(1)
	}
	r, err := rotator.New(logFile, 1024*1024, false, 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create file rotator: %v\n", err)
		os.Exit(1)
	}

	logRotator = r
}

// setLogLevel sets the logging level for provided subsystem.  Invalid
// subsystems are ignored.  Uninitialized subsystems are dynamically created as
// needed.
func SetLogLevel(subsystemID string, logLevel string) bool {
	// Ignore invalid subsystems.
	logger, ok := subsystemLoggers[subsystemID]
	if !ok {
		return false
	}

	// Defaults to info if the logger level is invalid.
	level, _ := LevelFromString(logLevel)
	logger.SetLevel(level)
	return true
}

// setLogLevels sets the logger level for all subsystem loggers to the passed
// level.  It also dynamically creates the subsystem loggers as needed, so it
// can be used to initialize the logging system.
func SetLogLevels(logLevel string) {
	// Configure all sub-systems with the new logging level.  Dynamically
	// create loggers as needed.
	for subsystemID := range subsystemLoggers {
		SetLogLevel(subsystemID, logLevel)
	}
}

// CloseLogRotator closes the log rotator if it has been created
func CloseLogRotator() {
	if logRotator != nil {
		logRotator.Close()
	}
}

func GetLogger(name string) Logger {
	log, _ := subsystemLoggers[name]
	return log
}

// validLogLevel returns whether or not logLevel is a valid debug logger level.
func ValidLogLevel(logLevel string) bool {
	switch logLevel {
	case "trace":
		fallthrough
	case "debug":
		fallthrough
	case "info":
		fallthrough
	case "warn":
		fallthrough
	case "error":
		fallthrough
	case "critical":
		return true
	}
	return false
}

// SupportedSubsystems returns a sorted slice of the supported subsystems for
// logging purposes.
func SupportedSubsystems() []string {
	// Convert the subsystemLoggers map keys to a slice.
	subsystems := make([]string, 0, len(subsystemLoggers))
	for subsysID := range subsystemLoggers {
		subsystems = append(subsystems, subsysID)
	}

	// Sort the subsystems for stable display.
	sort.Strings(subsystems)
	return subsystems
}
