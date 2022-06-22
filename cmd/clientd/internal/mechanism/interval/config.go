package interval

import (
	"errors"
	"time"
)

// Config is a type used to configure the interval mechanism.
type Config struct {
	Duration time.Duration
}

// DefaultConfig represents the default values defined for the struct fields
// on the Config type.
var DefaultConfig = Config{
	Duration: time.Minute,
}

// Validate validates the receiver struct fields.
func (c *Config) Validate() error {
	if c.Duration <= time.Second*5 {
		return errors.New("duration must be greater than 5s")
	}

	return nil
}
