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
	router := handlers.NewRouter(cfg, repository.New())
	httpServer := server.NewHTTPServer(cfg.Addr, router)
	logger.Initialize(cfg.LogLevel)

	err := httpServer.Run()
	if err != nil {
		panic(err)
	}
}
