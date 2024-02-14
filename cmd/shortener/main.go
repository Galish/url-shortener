// Package main is the entry point to the Shortener application.
package main

import (
	"fmt"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/Galish/url-shortener/pkg/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf(
		"Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildVersion,
		buildDate,
		buildCommit,
	)

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
