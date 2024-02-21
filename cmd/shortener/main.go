// Package main is the entry point to the Shortener application.
package main

import (
	"fmt"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository"
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

	fmt.Printf("Config: %+v\n", cfg)

	logger.Initialize(cfg.LogLevel)

	store, err := repository.New(cfg)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	handler := handlers.NewHandler(cfg, store)
	defer handler.Close()

	router := handlers.NewRouter(handler)
	server := handlers.NewServer(cfg, router)

	if err := server.Run(); err != nil {
		panic(err)
	}
}
