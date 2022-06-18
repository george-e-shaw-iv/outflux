// Package log implements a simple logger using only stdlib imports.
//
// This was created in order to not take a dependency on a cumbersome logging
// package.
package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Fields is a type used to embed fields in log lines.
type Fields map[string]interface{}

// Logger is the type that has the actual logger implementation.
type Logger struct {
	level   int
	format  int
	writers *writers

	// timestampFormat is passed to *time.Time.Format when logging.
	timestampFormat string
}

// NewLogger returns a Logger configured with the passed in functional options.
func NewLogger(opts ...Option) (*Logger, error) {
	l := Logger{
		level:           LevelInfo,
		format:          FormatJSON,
		writers:         defaultWriters(),
		timestampFormat: time.RFC3339,
	}

	for i := range opts {
		opts[i](&l)
	}

	return &l, l.validate()
}

// validate ensures that the options passed to the Logger are valid.
func (l *Logger) validate() error {
	if l.level < LevelDebug || l.level > LevelFatal {
		return fmt.Errorf("log level must be in range [%d,%d]", LevelDebug, LevelFatal)
	}

	if l.format < FormatJSON || l.format > FormatHumanReadable {
		return errors.New("invalid format")
	}

	if l.timestampFormat == "" {
		return errors.New("timestamp format cannot be empty")
	}

	return nil
}

// print is the underlying print function that all of the receiver functions on
// the Logger type that print use.
func (l *Logger) print(level int, message string, fields Fields) {
	if level < l.level {
		// If we're trying to log below the level we're configured to log for,
		// then we should shallow return here.
		return
	}

	timestamp := time.Now().UTC().Format(l.timestampFormat)

	switch l.format {
	case FormatJSON:
		line := formatJSON{
			Timestamp: timestamp,
			Level:     levels[level],
			Message:   message,
			Fields:    fields,
		}

		b, err := json.Marshal(line)
		if err == nil {
			fmt.Fprintln(l.writers.ByLevel(level), string(b))
		}

		fmt.Fprintf(l.writers.ByLevel(level), "error marshaling log line, falling back to human readable format: %v\n", err)
		fallthrough
	case FormatHumanReadable:
		fmt.Fprintf(l.writers.ByLevel(level), formatHumanReadable, timestamp, levels[level], message, fields)
	}
}

// Debug prints a debug level log line.
func (l *Logger) Debug(message string, fields Fields) {
	l.print(LevelDebug, message, fields)
}

// Info prints an info level log line.
func (l *Logger) Info(message string, fields Fields) {
	l.print(LevelInfo, message, fields)
}

// Warn prints a warn level log line.
func (l *Logger) Warn(message string, fields Fields) {
	l.print(LevelWarn, message, fields)
}

// Error prints a error level log line.
func (l *Logger) Error(message string, fields Fields) {
	l.print(LevelError, message, fields)
}

// Fatal prints a fatal level log line.
func (l *Logger) Fatal(message string, fields Fields) {
	l.print(LevelFatal, message, fields)
	os.Exit(1)
}
