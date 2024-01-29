// Package server provides a HTTP server.
package server

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
)

type httpServer struct {
	addr    string
	handler http.Handler
}

// NewHTTPServer creates an instance of HTTP server.
func NewHTTPServer(addr string, handler http.Handler) *httpServer {
	return &httpServer{
		addr,
		handler,
	}
}

// Run listens and serves requests sent to HTTP handlers.
func (s *httpServer) Run() error {
	logger.Info("running server on", s.addr)
	return http.ListenAndServe(s.addr, s.handler)
}
