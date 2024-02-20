package handlers

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/pkg/server"
)

var (
	hostWhitelist = []string{"urlshortener.io", "www.urlshortener.io"}
	dirCert       = "certs"
)

// NewServer configures and creates an HTTP server.
func NewServer(cfg *config.Config, handler http.Handler) *server.Server {
	var options = []server.Option{
		server.WithAddress(cfg.ServAddr),
	}

	if cfg.IsTLSEnabled {
		options = append(
			options,
			server.WithTLS(&server.TLSConfig{
				DirCert:       dirCert,
				HostWhitelist: hostWhitelist,
			}),
		)
	}

	return server.New(
		handler,
		options...,
	)
}
