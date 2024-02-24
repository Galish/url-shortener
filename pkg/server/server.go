package server

import (
	"context"
	"net/http"

	"github.com/Galish/url-shortener/pkg/logger"
)

// Server represents HTTP server.
type Server struct {
	*http.Server
}

// TLSConfig represents TLS server configuration.
type TLSConfig struct {
	DirCert       string
	HostWhitelist []string
}

// New creates an instance of HTTP server.
func New(handler http.Handler, options ...Option) *Server {
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	server := &Server{httpServer}

	for _, opt := range options {
		opt(server)
	}

	return server
}

// Run listens and serves requests sent to HTTP handlers.
func (s Server) Run() error {
	if s.TLSConfig == nil {
		logger.Info("running HTTP server on ", s.Addr)

		return s.ListenAndServe()
	}

	logger.Info("running HTTPS server on ", s.Addr)

	return s.ListenAndServeTLS("", "")
}

func (s *Server) Close() error {
	logger.Info("shutting down the server")

	return s.Shutdown(context.Background())
}
