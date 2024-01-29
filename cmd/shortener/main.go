// Package main is the entry point to the Shortener application.
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

	handler := handlers.NewHandler(cfg, store)
	defer handler.Close()

	router := handlers.NewRouter(handler)
	httpServer := server.NewHTTPServer(cfg.ServAddr, router)

	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
