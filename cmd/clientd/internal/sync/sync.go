// Package sync does the actual syncing between the client and the server instances
// of outflux and their InfluxDB integrations.
package sync

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/george-e-shaw-iv/outflux/internal/async"
	"github.com/george-e-shaw-iv/outflux/internal/grpc"
	"github.com/george-e-shaw-iv/outflux/internal/log"
)

//go:embed embed/copy-and-truncate.sh
var copyAndTruncateScript []byte

// Runner is the type that has receiver functions that perform the syncing
// operation between the client and the server instances of outflux and their
// InfluxDB integrations.
type Runner struct {
	mu     async.Mutex
	client *grpc.Client

	copyAndTruncateScriptPath string
}

// NewRunner returns a new initialized instance of Runner.
func NewRunner(ctx context.Context, outfluxClient *grpc.Client) (*Runner, error) {
	// We need to write the copy and truncate script that locks the file to
	// a temporary directory that will exist for the lifetime of the program.
	f, err := os.CreateTemp("", "copy-and-truncate-*.sh")
	if err != nil {
		return nil, fmt.Errorf("create temporary file for copy and truncate script: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(copyAndTruncateScript); err != nil {
		return nil, fmt.Errorf("write to temporary file for copy and truncate script: %w", err)
	}

	return &Runner{
		client:                    outfluxClient,
		copyAndTruncateScriptPath: f.Name(),
	}, nil
}

// Do performs the sync between the client and the server instances of outflux
// and their InfluxDB integrations.
func (r *Runner) Do(ctx context.Context) error {
	r.mu.Lock(ctx)
	defer r.mu.Unlock()

	// THIS WILL ALL HAPPEN IN THE NEXT PR:
	//
	// execute the script against the metrics file
	// read the outputted file (capture output filename from invocation of script)
	// marshal the influx line protocol format
	// send it over the wire to the server in chunks

	return nil
}

func (r *Runner) Close(ctx context.Context) error {
	if r.copyAndTruncateScriptPath != "" {
		if err := os.Remove(r.copyAndTruncateScriptPath); err != nil {
			log.Error("remove temporary copy and truncate script", log.Fields{
				"error": err.Error(),
			})
		}
	}

	return nil
}
