package server

// Option updates the setting in the server configuration.
type Option func(*Server)

// WithAddress updates the server address setting.
func WithAddress(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

// WithTLS updates the server TLS settings.
func WithTLS(cnf *TLSConfig) Option {
	return func(s *Server) {
		s.tlsConfig = cnf
	}
}
