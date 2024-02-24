// Package main is the entry point to the Shortener application.
package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/Galish/url-shortener/pkg/logger"
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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	store, err := repository.New(cfg)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	handler := handlers.NewHandler(cfg, store)

	g.Go(func() error {
		<-gCtx.Done()
		handler.Close()
		return nil
	})

	router := handlers.NewRouter(handler)
	server := handlers.NewServer(cfg, router)

	g.Go(func() error {
		return server.Run()
	})

	g.Go(func() error {
		<-gCtx.Done()
		return server.Close()
	})

	if err := g.Wait(); err != nil {
		fmt.Println("---shutting down---", err)
	}
}
