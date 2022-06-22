package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/omekov/superapp-backend/internal/salecar/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Server ...
type Server struct {
	httpServer *http.Server
	grpcServer *grpc.Server
}

// NewServer ...
func NewServer(cfg *config.Config, port string, handler http.Handler) *Server {
	opts := []grpc.ServerOption{grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
		Timeout:           15 * time.Second,
		MaxConnectionAge:  5 * time.Minute,
		Time:              120 * time.Minute,
	})}
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + port,
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		grpcServer: grpc.NewServer(opts...),
	}
}

// RunGRPC ...
func (s *Server) RunGRPC(lis net.Listener) error {
	return s.grpcServer.Serve(lis)
}

// StopGRPC ...
func (s *Server) StopGRPC() {
	s.grpcServer.GracefulStop()
}

// RunHTTP ...
func (s *Server) RunHTTP() error {
	return s.httpServer.ListenAndServe()
}

// StopHTTP ...
func (s *Server) StopHTTP(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
