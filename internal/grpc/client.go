package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is a wrapper type around a client connection to the gRPC server
// used to power outflux.
type Client struct {
	cc *grpc.ClientConn
	OutfluxClient
}

// NewClient returns a new OutfluxClient connected to the server at the address
// given as a parameter.
//
// TODO(george-e-shaw-iv): Expand this to allow options to be passed and be secure
// by default.
func NewClient(ctx context.Context, address string) (*Client, error) {
	cc, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to outflux gRPC server: %w", err)
	}

	return &Client{
		cc:            cc,
		OutfluxClient: NewOutfluxClient(cc),
	}, nil
}

// Close closes the underlying client connection to the Outflux gRPC server.
func (c *Client) Close() error {
	if err := c.cc.Close(); err != nil {
		return fmt.Errorf("close grpc client connection: %w", err)
	}
	return nil
}
