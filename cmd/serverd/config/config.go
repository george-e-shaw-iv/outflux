// Package config defines configuration which is consumed in the form of YAML or
// JSON files for the outflux server daemon.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the type that defines the structure of the outflux server
// configuration file.
type Config struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

// Default represents the default values defined for the struct fields on the
// Config type.
var Default = Config{
	Host: "127.0.0.1",
	Port: 8000,
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
