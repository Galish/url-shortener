// Package main is the entry point to the Shortener application.
package main

import (
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/Galish/url-shortener/pkg/logger"
	"github.com/Galish/url-shortener/pkg/shutdowner"
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

	logger.SetLevel(cfg.LogLevel)

	store, err := repository.New(cfg)
	if err != nil {
		panic(err)
	}

	handler := handlers.NewHandler(cfg, store)
	router := handlers.NewRouter(handler)
	server := handlers.NewServer(cfg, router)

	sd := shutdowner.New(server, handler, store)

	if err := server.Run(); err != http.ErrServerClosed {
		panic(err)
	}

	sd.Wait()
}
