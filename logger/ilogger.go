// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package logger

// Logger is an interface which describes a level-based logger.  A default
// implementation of Logger is implemented by this package and can be created
// by calling (*Backend).Logger.
type Logger interface {
	// Tracef formats message according to format specifier and writes to
	// to logger with LevelTrace.
	Tracef(format string, params ...interface{})

	// Debugf formats message according to format specifier and writes to
	// logger with LevelDebug.
	Debugf(format string, params ...interface{})

	// Infof formats message according to format specifier and writes to
	// logger with LevelInfo.
	Infof(format string, params ...interface{})

	// Warnf formats message according to format specifier and writes to
	// to logger with LevelWarn.
	Warnf(format string, params ...interface{})

	// Errorf formats message according to format specifier and writes to
	// to logger with LevelError.
	Errorf(format string, params ...interface{})

	// Criticalf formats message according to format specifier and writes to
	// logger with LevelCritical.
	Criticalf(format string, params ...interface{})

	// Trace formats message using the default formats for its operands
	// and writes to logger with LevelTrace.
	Trace(v ...interface{})

	// Debug formats message using the default formats for its operands
	// and writes to logger with LevelDebug.
	Debug(v ...interface{})

	// Info formats message using the default formats for its operands
	// and writes to logger with LevelInfo.
	Info(v ...interface{})

	// Warn formats message using the default formats for its operands
	// and writes to logger with LevelWarn.
	Warn(v ...interface{})

	// Error formats message using the default formats for its operands
	// and writes to logger with LevelError.
	Error(v ...interface{})

	// Critical formats message using the default formats for its operands
	// and writes to logger with LevelCritical.
	Critical(v ...interface{})

	// Level returns the current logging level.
	Level() Level

	// SetLevel changes the logging level to the passed level.
	SetLevel(level Level)
}
