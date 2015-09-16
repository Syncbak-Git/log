// Package log implements a simple log file. It is similar to the standard library log package,
// but it introduces log levels to control which log entries are actually written.
// Note that the various SetXXX() functions are not thread-safe and should be called before
// writing log entries (or at least while there are no parallel routines writing log entries).
// Package log is the successor to github.com/Syncbak-Git/logging.
package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

var output io.Writer = os.Stderr
var logLevel Level = LevelAll
var timestamp func() string = func() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

// Level is a logging level.
type Level uint64

const (
	LevelDebug Level = 1 << iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelPanic
	LevelCustom
	LevelAll  = 0xFFFF
	LevelNone = 0
)

// SetLogLevel controls which log entries are actually written.
// Multiple logging levels can be combined by ORing individual
// Level values, eg. LevelDebug|LevelError will log both DEBUG and ERROR entries. Alternatively,
// specific log levels can be suppresed via XORing with LevelAll, eg. LevelAll ^ LevelDebug
// will log everything except DEBUG entries.
// The default log level is LevelAll.
func SetLogLevel(l Level) {
	logLevel = l
}

// SetOutput directs log output to w. The default output is written to os.Stderr.
func SetOutput(w io.Writer) {
	output = w
}

// SetOuputFile is a convenience function to wrap SetOutput() for writing to a file.
func SetOutputFile(f string) error {
	w, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	SetOutput(w)
	return nil
}

// SetTimestamp allows the user to replace the default RFC3339Nano timestamp string. It
// is intended for creating deterministic test cases, but may be generally useful.
func SetTimestamp(f func() string) {
	timestamp = f
}

// Debug writes a DEBUG entry to the log file.
func Debug(format string, args ...interface{}) error {
	if logLevel&LevelDebug == 0 {
		return nil
	}
	return writeEntry("DEBUG", format, args...)
}

// Info writes a INFO entry to the log file.
func Info(format string, args ...interface{}) error {
	if logLevel&LevelInfo == 0 {
		return nil
	}
	return writeEntry("INFO", format, args...)
}

// Warning writes a WARNING entry to the log file.
func Warning(format string, args ...interface{}) error {
	if logLevel&LevelWarning == 0 {
		return nil
	}
	return writeEntry("WARNING", format, args...)
}

// Error writes an ERROR entry to the log file.
func Error(format string, args ...interface{}) error {
	if logLevel&LevelError == 0 {
		return nil
	}
	return writeEntry("ERROR", format, args...)
}

// Fatal writes a FATAL entry to the log file and then exits
// via os.Exit(1).
func Fatal(format string, args ...interface{}) error {
	if logLevel&LevelFatal == 0 {
		return nil
	}
	err := writeEntry("FATAL", format, args...)
	os.Exit(1)
	return err // won't actually execute
}

// Panic writes a PANIC entry to the log file and then
// calls panic() with the log entry.
func Panic(format string, args ...interface{}) error {
	if logLevel&LevelPanic == 0 {
		return nil
	}
	err := writeEntry("PANIC", format, args...)
	panic(fmt.Sprintf(format, args...))
	return err // won't actually execute
}

// Custom writes an entry with a caller-supplied log level string.
func Custom(level string, format string, args ...interface{}) error {
	if logLevel&LevelCustom == 0 {
		return nil
	}
	return writeEntry(level, format, args...)
}

func writeEntry(level string, format string, args ...interface{}) error {
	_, err := fmt.Fprintf(output, "%s\t%s\t%s\n", timestamp(), level, fmt.Sprintf(format, args...))
	return err
}
