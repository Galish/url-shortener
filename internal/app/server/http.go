package server

import (
	"fmt"
	"net/http"
)

type httpServer struct {
	addr    string
	handler http.Handler
}

func NewHTTPServer(addr string, handler http.Handler) httpServer {
	return httpServer{
		addr,
		handler,
	}
}

func (srv *httpServer) Run() error {
	fmt.Println("Running server on", srv.addr)
	return http.ListenAndServe(srv.addr, srv.handler)
}
