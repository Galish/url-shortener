package grpc

import (
	"log"
	"net"

	pb "github.com/Galish/url-shortener/api/grpc/proto"
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/usecase"
	"github.com/Galish/url-shortener/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ShortenerServer represents GRPC server.
type ShortenerServer struct {
	pb.UnimplementedShortenerServer

	cfg     *config.Config
	server  *grpc.Server
	usecase usecase.Shortener
}

// NewServer configures and creates a GRPC server.
func NewServer(cfg *config.Config, uc usecase.Shortener) *ShortenerServer {
	s := grpc.NewServer()
	reflection.Register(s)

	server := &ShortenerServer{
		cfg:     cfg,
		server:  s,
		usecase: uc,
	}

	pb.RegisterShortenerServer(s, server)

	return server
}

// Run listens and serves GRPC requests.
func (s *ShortenerServer) Run() error {
	listen, err := net.Listen("tcp", s.cfg.GRPCAddr)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("running GRPC server on ", s.cfg.GRPCAddr)

	return s.server.Serve(listen)
}

// Close  is executed to release the memory
func (s *ShortenerServer) Close() error {
	logger.Info("shutting down the GRPC server")

	s.server.GracefulStop()

	return nil
}

// Get user urls

// Delete user urls
