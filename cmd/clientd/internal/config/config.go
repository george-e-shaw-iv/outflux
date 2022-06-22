// Package config defines configuration which is consumed in the form of YAML or
// JSON files for the outflux client daemon.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/george-e-shaw-iv/outflux/cmd/clientd/internal/mechanism/demand"
	"github.com/george-e-shaw-iv/outflux/cmd/clientd/internal/mechanism/interval"
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

	Mechanism struct {
		Interval *interval.Config `json:"interval" yaml:"interval"`
		OnDemand *demand.Config   `json:"onDemand" yaml:"onDemand"`
	} `json:"mechanism" yaml:"mechanism"`
}

// Default represents the default values defined for the struct fields on the
// Config type.
var Default = Config{
	File: "/etc/outflux/metrics.out",
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
