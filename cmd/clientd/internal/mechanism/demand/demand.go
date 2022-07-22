// Package demand defines the on-demand syncing mechanism, syncing on-demand
// via an HTTP request.
package demand

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/george-e-shaw-iv/outflux/internal/log"
)

// Syncer implements the mechanism.Syncer type for the on-demand mechanism.
type Syncer struct {
	cfg *Config

	sync   chan struct{}
	server *http.Server

	once sync.Once
}

// NewSyncer validates the configuration passed in and returns a new syncer
// for interval on-demand syncing.
func NewSyncer(cfg *Config) (*Syncer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate on-demand config: %w", err)
	}

	return &Syncer{
		cfg: cfg,
	}, nil
}

// handler is the HTTP handler function used to intake the on-demand request
// to sync.
func (s *Syncer) handler(w http.ResponseWriter, r *http.Request) {
	if len(s.sync) == 0 {
		// Only send a sync signal if there isn't already one waiting.
		s.sync <- struct{}{}
	}

	w.WriteHeader(http.StatusNoContent)
}

// initialize initializes an HTTP server and channel to perform on-demand
// syncing via HTTP.
func (s *Syncer) initialize() {
	endpoint := strings.TrimSuffix(strings.TrimPrefix(s.cfg.Endpoint, "/"), "/")

	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, s.handler)

	s.sync = make(chan struct{}, 1)
	s.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", s.cfg.Port),
		Handler:        mux,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}

	log.Info("on-demand mechanism runner started http server", log.Fields{
		"address": fmt.Sprintf(":%d/%s", s.cfg.Port, endpoint),
	})

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Fatal("fatal error encountered running on-demand http server", log.Fields{
				"error": err.Error(),
			})
		}
	}()
}

// Name helps implement the mechanism.Syncer interface.
func (*Syncer) Name() string {
	return "on-demand"
}

// Sync helps implements the mechanism.Syncer interface. This function blocks
// until it is time to sync with the server.
func (s *Syncer) Sync(ctx context.Context) error {
	s.once.Do(s.initialize)

	// Block until the sync channel is sent on from receiving an HTTP or the
	// request or the context is canceled.
	select {
	case <-s.sync:
	case <-ctx.Done():
	}

	return nil
}

// Close helps implement the mechanism.Syncer interface and cleans up the opened
// time.Ticker.
func (s *Syncer) Close(ctx context.Context) error {
	if s.server != nil {
		if err := s.server.Shutdown(ctx); err != nil {
			log.Error("error gracefully shutting down http server for on-demand syncing, closing it forcefully now", log.Fields{
				"error": err.Error(),
			})

			if err := s.server.Close(); err != nil {
				log.Error("error forcefully shutting down http server for on-demand syncing", log.Fields{
					"error": err.Error(),
				})
			}
		}
	}

	if s.sync != nil {
		close(s.sync)
	}

	return nil
}
