package log

// Fields is a type used to embed fields in log lines.
type Fields map[string]interface{}

// WithError returns Fields that include an error.
func WithError(err error) Fields {
	return Fields{
		"error": err.Error(),
	}
}
