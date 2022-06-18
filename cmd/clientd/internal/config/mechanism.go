package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Mechanism Identifier Constants
//
// These are set on Config.mechanismIdentifier type to help quickly assert
// Config.mechanism which contains the contents of the Config.Mechanism
// generic map. The Config.Mechanism generic map is marshaled into the
// appropriate type upon marshal.
const (
	// interval refers to the mechanism defined by the Interval type.
	interval = iota
)

// Interval is a type that represents the configuration for the interval
// mechanism. This mechanism attempts to sync the client with the server on
// a defined interval.
type Interval struct {
	Duration time.Duration `json:"duration" yaml:"duration"`
}

// parseMechanism takes an unparsed mechanism from Config.Mechanism and parses
// it into its identifier (int) and type (interface{}). This function also returns
// an error if applicable.
func parseMechanism(unparsed map[string]interface{}) (int, interface{}, error) {
	mechanismTypeInterface, exists := unparsed["type"]
	if !exists {
		return 0, nil, errors.New("mechanism.type needs to be set")
	}

	mechanismType, ok := mechanismTypeInterface.(string)
	if !ok {
		return 0, nil, errors.New("mechanism.type should be a string")
	}

	delete(unparsed, "type")
	b, err := json.Marshal(unparsed)
	if err != nil {
		return 0, nil, fmt.Errorf("marshal mechanism back into bytes: %w", err)
	}

	switch m := strings.ToLower(mechanismType); m {
	case "interval":
		var parsed Interval
		if err := json.Unmarshal(b, &parsed); err != nil {
			return 0, nil, fmt.Errorf("unmarshal generic mechanism into interval: %w", err)
		}
		return interval, parsed, nil
	default:
		return 0, nil, fmt.Errorf("unknown mechanism.type %q", m)
	}
}
