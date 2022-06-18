package log

// singleton is the singleton version of the Logger type.
var singleton Logger

// Debug prints a debug level log line.
func Debug(message string, fields Fields) {
	singleton.Debug(message, fields)
}

// Info prints an info level log line.
func Info(message string, fields Fields) {
	singleton.Info(message, fields)
}

// Warn prints a warn level log line.
func Warn(message string, fields Fields) {
	singleton.Warn(message, fields)
}

// Error prints a error level log line.
func Error(message string, fields Fields) {
	singleton.Error(message, fields)
}

// Fatal prints a fatal level log line.
func Fatal(message string, fields Fields) {
	singleton.Fatal(message, fields)
}
