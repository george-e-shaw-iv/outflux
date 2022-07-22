package main

import (
	"context"
	"fmt"

	"github.com/george-e-shaw-iv/outflux/internal/grpc"
)

// Outflux implements the grpc.OutfluxServer interface.
type Outflux struct {
	grpc.UnimplementedOutfluxServer
}

// Sync is the handler that gets invoked when the Sync RPC gets invoked from a client
// wanting to sync its influx data to the server.
func (*Outflux) Sync(ctx context.Context, req *grpc.SyncRequest) (*grpc.SyncResponse, error) {
	fmt.Println("received sync request", len(req.DataPoints))

	// THIS IS WHERE WE WILL INTERACT WITH THE INFLUXDB INSTANCE ON THE SERVER
	// WHICH WILL HAPPEN IN THE NEXT PR

	return &grpc.SyncResponse{
		Failed: []uint32{},
	}, nil
}
