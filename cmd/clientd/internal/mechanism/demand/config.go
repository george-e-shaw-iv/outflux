package demand

import "errors"

// Config is a type used to configure the on-demand mechanism.
type Config struct {
	Port     int    `json:"port" yaml:"port"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// DefaultConfig represents the default values defined for the struct fields
// on the Config type.
var DefaultConfig = Config{
	Port:     8000,
	Endpoint: "sync",
}

// Validate validates the receiver struct fields.
func (c *Config) Validate() error {
	if c.Port <= 1023 || c.Port > 65535 {
		return errors.New("invalid port")
	}

	if c.Endpoint == "" {
		return errors.New("endpoint cannot be empty string")
	}

	return nil
}
