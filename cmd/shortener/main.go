package main

import (
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/Galish/url-shortener/internal/app/server"
)

func main() {
	cfg := config.New()

	logger.Initialize(cfg.LogLevel)

	store, err := repository.New(cfg)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	router := handlers.NewRouter(cfg, store)
	httpServer := server.NewHTTPServer(cfg.Addr, router)

	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
