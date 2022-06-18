// Package config defines configuration which is consumed in the form of YAML or
// JSON files for the outflux client daemon.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the type that defines the structure of the outflux client
// configuration file.
type Config struct {
	File   string `json:"file" yaml:"file"`
	Server struct {
		Address string `json:"address" yaml:"address"`
	} `json:"server" yaml:"server"`

	Chunk *struct {
		MaxSizeBytes int `json:"maxSizeBytes" yaml:"maxSizeBytes"`
		MaxNumPoints int `json:"maxNumPoints" yaml:"maxNumPoints"`
	} `json:"chunk" yaml:"chunk"`

	Mechanism_ map[string]interface{} `json:"mechanism" yaml:"mechanism"`

	// mechanismIdentifier maps to one of the constant identifiers defined
	// in a constant block above this type.
	mechanismIdentifier int
	mechanismName       string
	mechanism           interface{}
}

// UnmarshalJSON implements the json.Unmarshaler interface for the Config type.
func (c *Config) UnmarshalJSON(b []byte) error {
	var err error

	// create an intermediate type to avoid infinite recursion for unmarshaling.
	type _config Config

	var _c _config
	if err = json.Unmarshal(b, &_c); err != nil {
		return err
	}

	if _c.mechanismIdentifier, _c.mechanism, err = parseMechanism(_c.Mechanism_); err != nil {
		return fmt.Errorf("parse mechanism: %w", err)
	}

	// We don't need this information anymore, its wasted memory.
	_c.Mechanism_ = nil

	// Set the receiver.
	*c = Config(_c)

	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for the Config type.
func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	var err error

	// create an intermediate type to avoid infinite recursion for unmarshaling.
	type _config Config

	var _c _config
	if err = value.Decode(&_c); err != nil {
		return err
	}

	if _c.mechanismIdentifier, _c.mechanism, err = parseMechanism(_c.Mechanism_); err != nil {
		return fmt.Errorf("parse mechanism: %w", err)
	}

	// We don't need this information anymore, its wasted memory.
	_c.Mechanism_ = nil

	// Set the receiver.
	*c = Config(_c)

	return nil
}

// Defaults sets defaults on unset fields on the Config type.
func (c *Config) Defaults() {
	if strings.TrimSpace(c.File) == "" {
		c.File = "/etc/outflux/metrics.out"
	}
}

// MechanismName returns c.mechanismName, parsed from Config.Mechanism_ dynamically.
func (c *Config) MechanismName() string {
	return c.mechanismName
}

// Parse takes the filepath of a configuration file and parses the configuration.
func Parse(fp string) (*Config, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	splitFp := strings.Split(fp, string(os.PathSeparator))
	splitBasename := strings.Split(splitFp[len(splitFp)-1], ".")

	var c Config

	switch ext := splitBasename[len(splitBasename)-1]; ext {
	case "json":
		if err := json.NewDecoder(f).Decode(&c); err != nil {
			return nil, fmt.Errorf("unmarshal json: %w", err)
		}
	case "yaml", "yml":
		if err := yaml.NewDecoder(f).Decode(&c); err != nil {
			return nil, fmt.Errorf("unmarshal yaml: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown config file extension %q, valid extensions are \"yaml\", \"yml\", and \"json\"", ext)
	}

	return &c, nil
}
