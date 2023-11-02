package server

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
)

type httpServer struct {
	addr    string
	handler http.Handler
}

func NewHTTPServer(addr string, handler http.Handler) *httpServer {
	return &httpServer{
		addr,
		handler,
	}
}

func (s *httpServer) Run() error {
	logger.Info("Running server on", s.addr)
	return http.ListenAndServe(s.addr, s.handler)
}
