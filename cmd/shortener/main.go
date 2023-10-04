package main

import (
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/server"
	"github.com/Galish/url-shortener/internal/app/storage"
)

func main() {
	cfg := config.Get()
	store := storage.NewKeyValueStorage()
	router := handlers.NewRouter(cfg, store)
	httpServer := server.NewHTTPServer(cfg.Addr, router)

	err := httpServer.Run()
	if err != nil {
		panic(err)
	}
}
