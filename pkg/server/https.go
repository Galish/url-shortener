package server

import (
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/Galish/url-shortener/internal/app/logger"
)

var (
	hostWhitelist  = []string{"urlshortener.io", "www.urlshortener.io"}
	dirCertificate = "certs"
)

type httpsServer struct {
	server *http.Server
	addr   string
}

// NewHTTPSServer creates an instance of HTTPS server.
func NewHTTPSServer(addr string, handler http.Handler) *httpsServer {
	manager := &autocert.Manager{
		Cache:      autocert.DirCache(dirCertificate),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hostWhitelist...),
	}

	return &httpsServer{
		server: &http.Server{
			Addr:      addr,
			Handler:   handler,
			TLSConfig: manager.TLSConfig(),
		},
		addr: addr,
	}
}

// Run listens and serves requests sent to HTTP handlers.
func (s *httpsServer) Run() error {
	logger.Info("running server HTTPS on", s.addr)

	return s.server.ListenAndServeTLS("", "")
}
