package server

import (
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/Galish/url-shortener/internal/app/logger"
)

// Server represents HTTP server.
type Server struct {
	addr      string
	handler   http.Handler
	tlsConfig *TLSConfig
}

// TLSConfig represents TLS server configuration.
type TLSConfig struct {
	DirCert       string
	HostWhitelist []string
}

// New creates an instance of HTTP server.
func New(handler http.Handler, options ...Option) *Server {
	server := &Server{
		addr:    ":8080",
		handler: handler,
	}

	for _, opt := range options {
		opt(server)
	}

	return server
}

// Run listens and serves requests sent to HTTP handlers.
func (s Server) Run() error {
	logger.Info("running HTTP server on ", s.addr)

	if s.tlsConfig == nil {
		return http.ListenAndServe(s.addr, s.handler)
	}

	manager := &autocert.Manager{
		Cache:      autocert.DirCache(s.tlsConfig.DirCert),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(s.tlsConfig.HostWhitelist...),
	}

	server := &http.Server{
		Addr:      s.addr,
		Handler:   s.handler,
		TLSConfig: manager.TLSConfig(),
	}

	logger.Info("running server HTTPS on ", s.addr)

	return server.ListenAndServeTLS("", "")
}
