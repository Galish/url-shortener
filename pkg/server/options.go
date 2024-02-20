package server

type Option func(*Server)

func WithAddress(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithSecureTransport(cnf *TLSConfig) Option {
	return func(s *Server) {
		s.tlsConfig = cnf
	}
}
