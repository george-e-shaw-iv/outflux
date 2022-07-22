// Package mechanism defines the running of and orchestration of mechanisms
// for syncing.
package mechanism

import (
	"context"
	"errors"
	"fmt"
	"strings"
	_sync "sync"

	"github.com/george-e-shaw-iv/outflux/cmd/clientd/internal/config"
	"github.com/george-e-shaw-iv/outflux/cmd/clientd/internal/mechanism/demand"
	"github.com/george-e-shaw-iv/outflux/cmd/clientd/internal/mechanism/interval"
	"github.com/george-e-shaw-iv/outflux/cmd/clientd/internal/sync"
	"github.com/george-e-shaw-iv/outflux/internal/log"
)

// Syncer is the interface that a mechanism should implement in order to be
// used during the runtime of the syncing process for outflux.
type Syncer interface {
	// Name just returns the name of the Syncer.
	Name() string

	// Sync is a blocking function that returns a non-nil error when a sync
	// should occur.
	Sync(ctx context.Context) error

	// Close gets called when outflux wants to gracefully exit.
	Close(ctx context.Context) error
}

// Runner encapsulates mechanism running logic on its receivers to do orchestration
// for running multiple mechanisms and interacting with the primary syncing
// package.
type Runner struct {
	cfg     *config.Config
	sync    *sync.Runner
	syncers []Syncer

	wg _sync.WaitGroup
}

// NewRunner returns an initialized instance of Runner.
func NewRunner(cfg *config.Config, syncRunner *sync.Runner) *Runner {
	return &Runner{
		cfg:  cfg,
		sync: syncRunner,
	}
}

// RunAll runs all of the configured mechanisms and effectively starts the syncing
// loop. This is a blocking function.
func (r *Runner) RunAll(ctx context.Context) error {
	var syncers []Syncer

	if r.cfg.Mechanism.Interval != nil {
		intervalSyncer, err := interval.NewSyncer(r.cfg.Mechanism.Interval)
		if err != nil {
			return fmt.Errorf("initialize interval mechanism: %w", err)
		}
		syncers = append(syncers, intervalSyncer)
	}

	if r.cfg.Mechanism.OnDemand != nil {
		onDemandSyncer, err := demand.NewSyncer(r.cfg.Mechanism.OnDemand)
		if err != nil {
			return fmt.Errorf("initialize on-demand mechanism: %w", err)
		}
		syncers = append(syncers, onDemandSyncer)
	}

	if len(syncers) == 0 {
		return errors.New("no mechanisms specified for syncing")
	}

	for i := range syncers {
		syncer := syncers[i]

		r.wg.Add(1)
		go func(ctx context.Context, syncer Syncer) {
			defer r.wg.Done()

			for ctx.Err() == nil {
				if err := syncer.Sync(ctx); err != nil {
					log.Error("syncing mechanism failed to report sync", log.Fields{
						"error":     err.Error(),
						"mechanism": syncer.Name(),
					})
					return
				}

				log.Info("sync event fired", log.Fields{
					"mechanism": syncer.Name(),
				})

				if err := r.sync.Do(ctx); err != nil {
					log.Error("sync to server", log.Fields{
						"error": err.Error(),
					})
				}
			}
		}(ctx, syncer)
	}

	// Block until all go-routines running the syncing mechanisms have returned,
	// which will only happen on a context cancellation event.
	r.wg.Wait()

	return nil
}

// Close closes all of the syncers.
func (r *Runner) Close(ctx context.Context) error {
	var errs []string
	for i := range r.syncers {
		if err := r.syncers[i].Close(ctx); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}
