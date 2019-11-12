package log

import (
	"fmt"
	"io"
	"os"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// Log prints tagged log messages. It is a stub for future reporting of
// such messages to the nebula service
func Log(writer io.Writer, level LogLevel, log string) {
	fmt.Fprintln(
		writer,
		log,
	)
}

// Info reports an informational log
func Info(log string) {
	Log(os.Stdout, LogLevelInfo, log)
}

// Warn reports a warning
func Warn(log string) {
	Log(os.Stderr, LogLevelWarn, log)
}

// Error reports an error
func Error(log string) {
	Log(os.Stderr, LogLevelError, log)
}

// Fatal reports a fatal error then exits the process
func Fatal(log string) {
	Log(os.Stderr, LogLevelFatal, log)
	os.Exit(1)
}
