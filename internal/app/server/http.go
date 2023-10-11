package server

import (
	"fmt"
	"net/http"
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
	fmt.Println("Running server on", s.addr)
	return http.ListenAndServe(s.addr, s.handler)
}
