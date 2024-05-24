// Package log provides a global logger for logging.
package log

import (
	"context"
	"fmt"
	"io"
	"os"

	"omni/observer/logging"
)

// Logger is the global logger.
var Logger = logging.New(os.Stderr).With().Timestamp().Logger()

// Output duplicates the global logger and sets w as its output.
func Output(w io.Writer) logging.Logger {
	return Logger.Output(w)
}

// With creates a child logger with the field added to its context.
func With() logging.Context {
	return Logger.With()
}

// Level creates a child logger with the minimum accepted level set to level.
func Level(level logging.Level) logging.Logger {
	return Logger.Level(level)
}

// Sample returns a logger with the s sampler.
func Sample(s logging.Sampler) logging.Logger {
	return Logger.Sample(s)
}

// Hook returns a logger with the h Hook.
func Hook(h logging.Hook) logging.Logger {
	return Logger.Hook(h)
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func Err(err error) *logging.Event {
	return Logger.Err(err)
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *logging.Event {
	return Logger.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *logging.Event {
	return Logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *logging.Event {
	return Logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *logging.Event {
	return Logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *logging.Event {
	return Logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *logging.Event {
	return Logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *logging.Event {
	return Logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func WithLevel(level logging.Level) *logging.Event {
	return Logger.WithLevel(level)
}

// Log starts a new message with no level. Setting logging.GlobalLevel to
// logging.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *logging.Event {
	return Logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	Logger.Debug().CallerSkipFrame(1).Msg(fmt.Sprint(v...))
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Logger.Debug().CallerSkipFrame(1).Msgf(format, v...)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *logging.Logger {
	return logging.Ctx(ctx)
}
