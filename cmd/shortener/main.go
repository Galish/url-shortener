package main

import (
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/server"
	"github.com/Galish/url-shortener/internal/app/storage/kvstorage"
)

func main() {
	cfg := config.New()
	router := handlers.NewRouter(cfg, kvstorage.New())
	httpServer := server.NewHTTPServer(cfg.Addr, router)

	err := httpServer.Run()
	if err != nil {
		panic(err)
	}
}
