// Package main is the entry point to the Shortener application.
package main

import (
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/api/grpc"
	restapi "github.com/Galish/url-shortener/api/rest"
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/Galish/url-shortener/internal/app/usecase"
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

	logger.Init()

	cfg := config.New()

	logger.SetLevel(cfg.LogLevel)

	store, err := repository.New(cfg)
	if err != nil {
		panic(err)
	}

	shortener := usecase.New(cfg, store)

	handler := restapi.NewHandler(shortener, store)
	router := restapi.NewRouter(cfg, handler)
	restServer := restapi.NewServer(cfg, router)

	grpcServer := grpc.NewServer(cfg, shortener)

	sd := shutdowner.New(restServer, grpcServer, shortener, store)

	go func() {
		if err := restServer.Run(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	go func() {
		if err := grpcServer.Run(); err != nil {
			panic(err)
		}
	}()

	sd.Wait()
}
