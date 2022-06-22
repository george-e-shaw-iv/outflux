// Package sync does the actual syncing between the client and the server instances
// of outflux and their InfluxDB integrations.
package sync

import (
	"context"

	"github.com/george-e-shaw-iv/outflux/cmd/clientd/internal/config"
)

// Runner is the type that has receiver functions that perform the syncing
// operation between the client and the server instances of outflux and their
// InfluxDB integrations.
type Runner struct {
	cfg *config.Config
}

// NewRunner returns a new initialized instance of Runner.
func NewRunner(cfg *config.Config) *Runner {
	return &Runner{
		cfg: cfg,
	}
}

// Do performs the sync between the client and the server instances of outflux
// and their InfluxDB integrations.
func (r *Runner) Do(ctx context.Context) error {
	return nil
}
