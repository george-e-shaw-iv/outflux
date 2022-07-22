package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ServerOption is the type that implements functional options for the Server
// type.
type ServerOption func(*serverOptions)

// WithPort overrides the default port for the gRPC server, which is 8000.
func WithPort(port int) ServerOption {
	return func(s *serverOptions) {
		s.port = port
	}
}

// WithHost overrides the default host for the gRPC server, which is 127.0.0.1.
func WithHost(host string) ServerOption {
	return func(s *serverOptions) {
		s.host = host
	}
}

// serverOptions is a struct that contains all of the configurable options that
// the Server type can utilize.
type serverOptions struct {
	port int
	host string
}

// defaultServerOptions is just serverOptions with each applicable struct field
// given an explicit default value.
var defaultServerOptions = serverOptions{
	port: 8000,
	host: "127.0.0.1",
}

// Server is a type that wraps the gRPC server for Outflux.
type Server struct {
	opts     serverOptions
	listener net.Listener
	gs       *grpc.Server
}

// NewServer returns a new instance of Server configured to serve gRPC requests using
// the interface passed to it and any options given in the parameters.
//
// *Server.Close should be deferred as soon as this function returns with a nil error.
func NewServer(ctx context.Context, outflux OutfluxServer, opts ...ServerOption) (*Server, error) {
	var err error

	serverOpts := defaultServerOptions
	for i := range opts {
		opts[i](&serverOpts)
	}

	s := Server{
		opts: serverOpts,
	}

	lc := &net.ListenConfig{}
	s.listener, err = lc.Listen(ctx, "tcp", s.Address())
	if err != nil {
		return nil, fmt.Errorf("attempt to listen on address %q using tcp protocol: %w", s.Address(), err)
	}

	s.gs = grpc.NewServer()
	RegisterOutfluxServer(s.gs, outflux)
	reflection.Register(s.gs)

	return &s, nil
}

// Address returns the server address that it is listening on.
func (s *Server) Address() string {
	return fmt.Sprintf("%s:%d", s.opts.host, s.opts.port)
}

// Listen is a blocking function that effectively runs the gRPC server. This function
// will return if the passed in context is canceled or the server encounters a fatal
// error while serving.
func (s *Server) Listen(ctx context.Context) error {
	serverErr := make(chan error, 1)
	go func() {
		if err := s.gs.Serve(s.listener); err != nil {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case e := <-serverErr:
		return e
	}
}

// Close cleans up underlying in-flight RPCs, connections to the server, and the underlying
// TCP listener.
func (s *Server) Close(ctx context.Context) error {
	if s.gs != nil {
		closed := make(chan struct{})
		go func() {
			s.gs.GracefulStop()
			close(closed)
		}()

		select {
		case <-ctx.Done():
			// If the passed in context is canceled before the graceful shutdown
			// process finishes, force stop the server.
			s.gs.Stop()
		case <-closed:
		}
	}

	// This should realistically be closed almost always by this point in
	// time, but this is insurance just incase it isn't.
	if s.listener != nil {
		// We eat the error here because we don't care if it's already closed,
		// which it will report via error.
		_ = s.listener.Close()
	}
	return nil
}
