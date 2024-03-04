package server

import "golang.org/x/crypto/acme/autocert"

// Option updates the setting in the server configuration.
type Option func(*Server)

// WithAddress updates the server address setting.
func WithAddress(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}

// WithTLS updates the server TLS settings.
func WithTLS(cfg *TLSConfig) Option {
	return func(s *Server) {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache(cfg.DirCert),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.HostWhitelist...),
		}

		s.TLSConfig = manager.TLSConfig()
	}
}
