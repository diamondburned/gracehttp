package gracehttp

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/http2"
)

// WaitForInterrupt blocks until an interrupt signal is received.
func WaitForInterrupt() os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	return <-sig
}

// Server is the wrapper around a regular net/http server. It is not
// thread-safe, and must never be shared across goroutines.
type Server struct {
	*http.Server
}

// ListenAndServeAsync listens and serves the given address. The returned
// server can be used to gracefully shut down the server.
func ListenAndServeAsync(addr string, h http.Handler) *Server {
	s := MustServer(addr, h)
	if err := s.ListenAndServeAsync(context.Background()); err != nil {
		log.Fatalln("failed to listen and serve:", err)
	}

	return s
}

// MustServer creates a new server and panics if it fails.
func MustServer(addr string, h http.Handler) *Server {
	s, err := NewServer(addr, h)
	if err != nil {
		log.Fatalln("failed to create gracehttp server:", err)
	}
	return s
}

// NewServer creates a new graceful server instance with defaults.
func NewServer(addr string, h http.Handler) (*Server, error) {
	h1 := http.Server{
		Addr:    addr,
		Handler: h,
	}
	h2 := http2.Server{
		MaxHandlers:          10240,
		MaxConcurrentStreams: 4096,
	}

	return NewCustomServer(&h1, &h2)
}

// NewCustomServer creates a new graceful server instance from the given
// configs. If h2 is nil, then HTTP2 is not forced.
func NewCustomServer(h1 *http.Server, h2 *http2.Server) (*Server, error) {
	if h2 != nil {
		if err := http2.ConfigureServer(h1, h2); err != nil {
			return nil, errors.Wrap(err, "failed to configure http2")
		}
	}

	return &Server{Server: h1}, nil
}

// ListenAndServe listens to and serves the server's address. The context is
// used for timing out the initial listen.
func (s *Server) ListenAndServe(ctx context.Context) error {
	l, err := ListenAddrCfg(ctx, s.Addr, net.ListenConfig{})
	if err != nil {
		return err
	}

	return s.Server.Serve(l)
}

// ListenAndServeAsync listens to the server's address and serves in a
// background goroutine. The context is used for timing out the initial listen.
func (s *Server) ListenAndServeAsync(ctx context.Context) error {
	l, err := ListenAddrCfg(ctx, s.Addr, net.ListenConfig{})
	if err != nil {
		return err
	}

	go s.Server.Serve(l)
	return nil
}

// ShutdownTimeout is a convenient function to allow graceful shutdown for the
// given duration.
func (s *Server) ShutdownTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.Server.Shutdown(ctx)
}
