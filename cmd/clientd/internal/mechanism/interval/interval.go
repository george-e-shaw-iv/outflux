// Package interval defines the interval syncing mechanism, syncing on an
// interval.
package interval

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Syncer implements the mechanism.Syncer type for the interval mechanism.
type Syncer struct {
	cfg *Config

	once   sync.Once
	ticker *time.Ticker
}

// NewSyncer validates the configuration passed in and returns a new syncer
// for interval-based syncing.
func NewSyncer(cfg *Config) (*Syncer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate interval config: %w", err)
	}

	return &Syncer{
		cfg: cfg,
	}, nil
}

// initialize initializes a time.Ticker for interval-based syncing that the
// receiver carries out.
func (s *Syncer) initialize() {
	s.ticker = time.NewTicker(s.cfg.Duration)
}

// Name helps implement the mechanism.Syncer interface.
func (*Syncer) Name() string {
	return "interval"
}

// Sync helps implements the mechanism.Syncer interface. This function blocks
// until it is time to sync with the server.
func (s *Syncer) Sync(ctx context.Context) error {
	s.once.Do(s.initialize)

	// Block until the ticker on the receiver "ticks" or the context is
	// canceled.
	select {
	case <-s.ticker.C:
	case <-ctx.Done():
	}

	return nil
}

// Close helps implement the mechanism.Syncer interface and cleans up the opened
// time.Ticker.
func (s *Syncer) Close(ctx context.Context) error {
	if s.ticker != nil {
		s.ticker.Stop()
	}

	return nil
}
