package server

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/config"
)

type server interface {
	Run() error
}

// New creates a server instance.
func New(cfg *config.Config, handler http.Handler) server {
	if cfg.IsHTTPSEnabled {
		return NewHTTPSServer(cfg.ServAddr, handler)
	}

	return NewHTTPServer(cfg.ServAddr, handler)
}
