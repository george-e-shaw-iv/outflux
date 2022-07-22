package log

import (
	"encoding/json"
	"io"
)

// Constant block for log level definition identifiers.
const (
	LevelDebug = iota - 1
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// levels maps level identifiers defined in the constant block above to string
// representations.
var levels map[int]string = map[int]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
}

// Constant block for formatting identifiers.
const (
	FormatJSON = iota
	FormatHumanReadable
)

// formatJSON is the JSON representation of a log line.
type formatJSON struct {
	Timestamp    string   `json:"timestamp"`
	Level        string   `json:"level"`
	Message      string   `json:"message"`
	Fields       []Fields `json:"-"`
	MergedFields Fields   `json:"fields,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface for formatJSON. Right now
// the only purpose this serves is to generate the MergedFields struct field.
func (f formatJSON) MarshalJSON() ([]byte, error) {
	if len(f.Fields) > 0 {
		f.MergedFields = make(Fields)
	}

	for i := range f.Fields {
		for k, v := range f.Fields[i] {
			f.MergedFields[k] = v
		}
	}

	// This data is now worthless.
	f.Fields = nil

	// Create an intermediate type to avoid infinite recursion on the json.Marshaler
	// interface implementation.
	type _formatJSON formatJSON
	_f := _formatJSON(f)

	return json.Marshal(_f)
}

// formatHumanReadable is the human readable version of a log line, in the form
// of a string with verbs that should be passed to fmt.Printf or an alike method.
var formatHumanReadable = `TIMESTAMP=%s LEVEL=%s\tMESSAGE=%s\tFIELDS=%v\n`

// Option is the functional option type to configure the Logger type.
type Option func(*Logger)

// WithLevel sets the log level on the Logger type. By default, the log level is
// set to LevelInfo.
func WithLevel(level int) Option {
	return func(l *Logger) {
		l.level = level
	}
}

// WithFormat sets the log format on the Logger type. By default the format used
// is FormatJSON.
func WithFormat(format int) Option {
	return func(l *Logger) {
		l.format = format
	}
}

// WithTimestampFormat sets the timestamp format on each log line. This format string
// gets passed to time.Time.Parse. By default the format used is time.RFC3339.
func WithTimestampFormat(format string) Option {
	return func(l *Logger) {
		l.timestampFormat = format
	}
}

// WithWriter sets the writer for all log levels. By default the writers for each
// log level are the following:
//	debug: stdout
//	info: stdout
//	warn: stdout
//	error: stderr
//	fatal: stderr
func WithWriter(writer io.Writer) Option {
	return func(l *Logger) {
		l.writers.SetWriterOnAll(writer)
	}
}

// WithWriterForLevels sets the writer for groups of individual log levels. By default
// the writers for each log level are the following:
//	debug: stdout
//	info: stdout
//	warn: stdout
//	error: stderr
//	fatal: stderr
func WithWriterForLevels(writer io.Writer, levels ...int) Option {
	return func(l *Logger) {
		l.writers.SetWriterOnLevels(writer, levels...)
	}
}
